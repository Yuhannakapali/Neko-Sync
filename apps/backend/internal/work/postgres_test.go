package work_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"nekosync/internal/platform/entity"
	"nekosync/internal/platform/postgres/pgtest"
	"nekosync/internal/work"
)

func newWork(title string, kind work.Kind, providers map[string]string) *work.Work {
	now := time.Now().UTC().Truncate(time.Microsecond)
	return &work.Work{
		BaseEntity:  entity.BaseEntity{ID: entity.NewUUID(), CreatedAt: now, UpdatedAt: now},
		Kind:        kind,
		Title:       title,
		Titles:      []work.LocalizedTitle{{Lang: "ja", Title: title + " (ja)"}},
		Genres:      []string{"action", "drama"},
		ProviderIDs: providers,
	}
}

func TestWorkRoundTrip(t *testing.T) {
	ctx := context.Background()
	repo := work.NewPostgresRepository(pgtest.Open(t))

	year := 2013
	w := newWork("Attack on Titan", work.KindAnime, map[string]string{"anilist": "16498"})
	w.Year = &year
	w.Artwork = []work.Artwork{{Kind: work.ArtworkPoster, URL: "https://img.example/aot.jpg"}}
	if err := repo.Create(ctx, w); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetByID(ctx, w.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != w.Title || *got.Year != 2013 || len(got.Genres) != 2 || got.Titles[0].Lang != "ja" ||
		got.Artwork[0].URL != w.Artwork[0].URL || got.ProviderIDs["anilist"] != "16498" {
		t.Fatalf("round trip mismatch: %+v", got)
	}

	byProvider, err := repo.GetByProviderID(ctx, "anilist", "16498")
	if err != nil || byProvider.ID != w.ID {
		t.Fatalf("GetByProviderID: %v, %v", byProvider, err)
	}
	if _, err := repo.GetByProviderID(ctx, "anilist", "nope"); !errors.Is(err, work.ErrNotFound) {
		t.Fatalf("missing provider id: %v", err)
	}

	w.Title = "Shingeki no Kyojin"
	w.Tags = nil
	if err := repo.Update(ctx, w); err != nil {
		t.Fatal(err)
	}
	got, _ = repo.GetByID(ctx, w.ID)
	if got.Title != "Shingeki no Kyojin" || got.Tags == nil {
		t.Fatalf("update: %+v", got)
	}

	if err := repo.Delete(ctx, w.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetByID(ctx, w.ID); !errors.Is(err, work.ErrNotFound) {
		t.Fatalf("after delete: %v", err)
	}
	if err := repo.Delete(ctx, w.ID); !errors.Is(err, work.ErrNotFound) {
		t.Fatalf("second delete: %v", err)
	}
}

func TestWorkListAndSearch(t *testing.T) {
	ctx := context.Background()
	repo := work.NewPostgresRepository(pgtest.Open(t))

	for _, w := range []*work.Work{
		newWork("Berserk", work.KindManga, nil),
		newWork("Bleach", work.KindAnime, nil),
		newWork("100%_Real", work.KindMovie, nil),
	} {
		if err := repo.Create(ctx, w); err != nil {
			t.Fatal(err)
		}
	}

	all, err := repo.List(ctx, "", 10, 0)
	if err != nil || len(all) != 3 {
		t.Fatalf("List all: %d, %v", len(all), err)
	}
	anime, _ := repo.List(ctx, work.KindAnime, 10, 0)
	if len(anime) != 1 || anime[0].Title != "Bleach" {
		t.Fatalf("List anime: %+v", anime)
	}

	found, _ := repo.Search(ctx, "BER", "", 10, 0)
	if len(found) != 1 || found[0].Title != "Berserk" {
		t.Fatalf("Search title: %+v", found)
	}
	found, _ = repo.Search(ctx, "bleach (ja)", "", 10, 0)
	if len(found) != 1 {
		t.Fatalf("Search localized title: %+v", found)
	}
	found, _ = repo.Search(ctx, "%_", "", 10, 0)
	if len(found) != 1 || found[0].Title != "100%_Real" {
		t.Fatalf("LIKE wildcards must be literal: %+v", found)
	}
}

func TestWorkChildren(t *testing.T) {
	ctx := context.Background()
	repo := work.NewPostgresRepository(pgtest.Open(t))

	w := newWork("Show", work.KindSeries, nil)
	if err := repo.Create(ctx, w); err != nil {
		t.Fatal(err)
	}
	child := func(parent *entity.UUID, kind work.ChildKind, ordinal float64) *work.WorkChild {
		return &work.WorkChild{
			BaseEntity: entity.BaseEntity{ID: entity.NewUUID(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
			WorkID:     w.ID, ParentID: parent, Kind: kind, Ordinal: ordinal,
		}
	}

	s1, s2 := child(nil, work.ChildSeason, 1), child(nil, work.ChildSeason, 2)
	for _, c := range []*work.WorkChild{s1, s2, child(&s1.ID, work.ChildEpisode, 1), child(&s2.ID, work.ChildEpisode, 1), child(&s1.ID, work.ChildEpisode, 1.5)} {
		if err := repo.CreateChild(ctx, c); err != nil {
			t.Fatalf("create %v %v: %v", c.Kind, c.Ordinal, err)
		}
	}
	if err := repo.CreateChild(ctx, child(&s1.ID, work.ChildEpisode, 1)); err == nil {
		t.Fatal("duplicate episode within a season must be rejected")
	}
	if err := repo.CreateChild(ctx, child(nil, work.ChildSeason, 1)); err == nil {
		t.Fatal("duplicate top-level season must be rejected")
	}

	ep, err := repo.GetChildByOrdinal(ctx, w.ID, &s2.ID, work.ChildEpisode, 1)
	if err != nil || *ep.ParentID != s2.ID {
		t.Fatalf("season 2 episode 1: %+v, %v", ep, err)
	}
	top, err := repo.GetChildByOrdinal(ctx, w.ID, nil, work.ChildSeason, 2)
	if err != nil || top.ID != s2.ID {
		t.Fatalf("top-level season 2: %+v, %v", top, err)
	}
	if _, err := repo.GetChildByOrdinal(ctx, w.ID, &s1.ID, work.ChildEpisode, 9); !errors.Is(err, work.ErrChildNotFound) {
		t.Fatalf("missing child: %v", err)
	}

	children, _ := repo.ListChildren(ctx, w.ID)
	if len(children) != 5 {
		t.Fatalf("ListChildren: %d", len(children))
	}

	if err := repo.Delete(ctx, w.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetChild(ctx, s1.ID); !errors.Is(err, work.ErrChildNotFound) {
		t.Fatalf("children should cascade with their work: %v", err)
	}
}
