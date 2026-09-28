package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"nekosync-instance/internal/domain/library"
	"nekosync-instance/internal/domain/mediafile"
	"nekosync-instance/internal/domain/shared"
)

func openTest(t *testing.T) *sql.DB {
	t.Helper()
	db, err := Open(context.Background(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestMigrateIsIdempotent(t *testing.T) {
	db := openTest(t)
	if err := Migrate(context.Background(), db); err != nil {
		t.Fatalf("second migrate: %v", err)
	}
	var v int
	_ = db.QueryRow("PRAGMA user_version").Scan(&v)
	if v != 1 {
		t.Fatalf("user_version = %d, want 1", v)
	}
}

func TestRepositoriesRoundTrip(t *testing.T) {
	ctx := context.Background()
	db := openTest(t)
	libs, files := NewLibraryRepository(db), NewMediaFileRepository(db)

	now := time.Now().UTC()
	lib := &library.Library{Name: "Anime", Kind: mediafile.KindAnime, Root: "/media/anime"}
	lib.ID, lib.CreatedAt, lib.UpdatedAt = shared.NewUUID(), now, now
	if err := libs.Save(ctx, lib); err != nil {
		t.Fatal(err)
	}
	if got, err := libs.GetByRoot(ctx, "/media/anime"); err != nil || got.ID != lib.ID {
		t.Fatalf("GetByRoot = %+v, %v", got, err)
	}

	ep := 7.5
	f := &mediafile.MediaFile{
		LibraryID: lib.ID, Path: "/media/anime/Gintama/Gintama - 7.5.mkv",
		Size: 123, ModTime: time.Unix(0, 1_700_000_000_123_456_789), Kind: mediafile.KindAnime,
		Technical: mediafile.Technical{Container: "matroska,webm", Duration: 1420.5, Streams: []mediafile.Stream{
			{Index: 1, Type: mediafile.StreamAudio, Codec: "aac", Language: "jpn", Default: true},
		}},
		ParsedTitle: "Gintama", Episode: &ep,
		ProviderIDs: map[string]string{"tmdb": "57911"}, MatchStatus: mediafile.MatchMatched,
	}
	f.ID, f.CreatedAt, f.UpdatedAt = shared.NewUUID(), now, now
	if err := files.Save(ctx, f); err != nil {
		t.Fatal(err)
	}

	got, err := files.GetByPath(ctx, f.Path)
	if err != nil {
		t.Fatal(err)
	}
	if !got.ModTime.Equal(f.ModTime) {
		t.Errorf("ModTime lost precision: %v vs %v", got.ModTime, f.ModTime)
	}
	if got.Episode == nil || *got.Episode != 7.5 || got.Season != nil || got.ParsedYear != nil {
		t.Errorf("nullable fields: ep=%v season=%v year=%v", got.Episode, got.Season, got.ParsedYear)
	}
	if got.ProviderIDs["tmdb"] != "57911" || len(got.Streams) != 1 || got.Streams[0].Language != "jpn" {
		t.Errorf("json fields: %+v %+v", got.ProviderIDs, got.Streams)
	}

	// Upsert keeps the row and updates fields.
	f.Size = 456
	if err := files.Save(ctx, f); err != nil {
		t.Fatal(err)
	}
	list, _ := files.ListByLibrary(ctx, lib.ID)
	if len(list) != 1 || list[0].Size != 456 {
		t.Fatalf("after upsert: %d rows, size %d", len(list), list[0].Size)
	}

	if err := files.Delete(ctx, f.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := files.Get(ctx, f.ID); !errors.Is(err, mediafile.ErrNotFound) {
		t.Fatalf("after delete: %v", err)
	}
}
