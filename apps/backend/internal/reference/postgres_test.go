package reference_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"nekosync/internal/platform/entity"
	"nekosync/internal/platform/postgres/pgtest"
	"nekosync/internal/reference"
)

func seed(t *testing.T, db *sql.DB) (userID, workID, childID entity.UUID) {
	t.Helper()
	userID, workID, childID = entity.NewUUID(), entity.NewUUID(), entity.NewUUID()
	for _, q := range []struct {
		sql  string
		args []any
	}{
		{`INSERT INTO users (id, username, email, password_hash) VALUES ($1, 'u', 'u@example.com', 'x')`, []any{userID}},
		{`INSERT INTO works (id, kind, title) VALUES ($1, 'series', 'Show')`, []any{workID}},
		{`INSERT INTO work_children (id, work_id, kind, ordinal) VALUES ($1, $2, 'episode', 1)`, []any{childID, workID}},
	} {
		if _, err := db.Exec(q.sql, q.args...); err != nil {
			t.Fatal(err)
		}
	}
	return userID, workID, childID
}

func ref(userID, workID entity.UUID, childID *entity.UUID, source reference.SourceKind, locator string) *reference.ContentReference {
	now := time.Now()
	return &reference.ContentReference{
		BaseEntity: entity.BaseEntity{ID: entity.NewUUID(), CreatedAt: now, UpdatedAt: now},
		UserID:     userID, WorkID: workID, ChildID: childID, Source: source, Locator: locator,
	}
}

func TestUpsertReconcilesRescans(t *testing.T) {
	ctx := context.Background()
	db := pgtest.Open(t)
	repo := reference.NewPostgresRepository(db)
	userID, workID, childID := seed(t, db)

	first := ref(userID, workID, &childID, reference.SourceInstance, "inst-1:node-9")
	if err := repo.Upsert(ctx, first); err != nil {
		t.Fatal(err)
	}

	hd := "1080p"
	rescan := ref(userID, workID, &childID, reference.SourceInstance, "inst-1:node-9")
	rescan.Quality = &hd
	if err := repo.Upsert(ctx, rescan); err != nil {
		t.Fatal(err)
	}
	if rescan.ID != first.ID {
		t.Fatalf("rescan got new id %s, want existing %s", rescan.ID, first.ID)
	}

	refs, err := repo.Resolve(ctx, userID, workID, &childID)
	if err != nil {
		t.Fatal(err)
	}
	if len(refs) != 1 || refs[0].Quality == nil || *refs[0].Quality != "1080p" {
		t.Fatalf("expected one refreshed reference, got %+v", refs)
	}

	// Whole-work references (movies) are distinct from child references and
	// also reconcile, even though child_id is NULL.
	for range 2 {
		if err := repo.Upsert(ctx, ref(userID, workID, nil, reference.SourceInstance, "inst-1:node-9")); err != nil {
			t.Fatal(err)
		}
	}
	whole, _ := repo.Resolve(ctx, userID, workID, nil)
	if len(whole) != 1 {
		t.Fatalf("whole-work references: %d, want 1", len(whole))
	}
}

func TestResolveAndRank(t *testing.T) {
	ctx := context.Background()
	db := pgtest.Open(t)
	repo := reference.NewPostgresRepository(db)
	userID, workID, childID := seed(t, db)

	for _, r := range []*reference.ContentReference{
		ref(userID, workID, &childID, reference.SourceOfficialLink, "https://licensed.example/ep1"),
		ref(userID, workID, &childID, reference.SourceInstance, "inst-1:node-9"),
	} {
		if err := repo.Create(ctx, r); err != nil {
			t.Fatal(err)
		}
	}

	refs, _ := repo.Resolve(ctx, userID, workID, &childID)
	ranked := reference.Rank(refs)
	if len(ranked) != 2 || ranked[0].Source != reference.SourceInstance {
		t.Fatalf("ranked: %+v", ranked)
	}

	if err := repo.Delete(ctx, ranked[0].ID); err != nil {
		t.Fatal(err)
	}
	if err := repo.Delete(ctx, ranked[0].ID); !errors.Is(err, reference.ErrNotFound) {
		t.Fatalf("second delete: %v", err)
	}
	all, _ := repo.ListByUser(ctx, userID, 10, 0)
	if len(all) != 1 {
		t.Fatalf("ListByUser: %d", len(all))
	}
}
