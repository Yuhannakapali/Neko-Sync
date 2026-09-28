package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"nekosync-instance/internal/domain/library"
	"nekosync-instance/internal/domain/mediafile"
	"nekosync-instance/internal/domain/shared"
)

// LibraryRepository implements library.Repository.
type LibraryRepository struct{ db *sql.DB }

// NewLibraryRepository builds a LibraryRepository.
func NewLibraryRepository(db *sql.DB) *LibraryRepository { return &LibraryRepository{db: db} }

const libraryCols = `id, name, kind, root, created_at, updated_at`

func (r *LibraryRepository) Get(ctx context.Context, id shared.UUID) (*library.Library, error) {
	return r.one(ctx, `SELECT `+libraryCols+` FROM libraries WHERE id = ?`, string(id))
}

func (r *LibraryRepository) GetByRoot(ctx context.Context, root string) (*library.Library, error) {
	return r.one(ctx, `SELECT `+libraryCols+` FROM libraries WHERE root = ?`, root)
}

func (r *LibraryRepository) List(ctx context.Context) ([]*library.Library, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+libraryCols+` FROM libraries ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*library.Library
	for rows.Next() {
		l, err := scanLibrary(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *LibraryRepository) Save(ctx context.Context, l *library.Library) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO libraries (`+libraryCols+`) VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
			name = excluded.name, kind = excluded.kind, root = excluded.root,
			updated_at = excluded.updated_at`,
		string(l.ID), l.Name, string(l.Kind), l.Root, toNanos(l.CreatedAt), toNanos(l.UpdatedAt))
	return err
}

func (r *LibraryRepository) one(ctx context.Context, q string, arg any) (*library.Library, error) {
	l, err := scanLibrary(r.db.QueryRowContext(ctx, q, arg))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, library.ErrNotFound
	}
	return l, err
}

type scanner interface{ Scan(dest ...any) error }

func scanLibrary(s scanner) (*library.Library, error) {
	var (
		l                library.Library
		id, kind         string
		created, updated int64
	)
	if err := s.Scan(&id, &l.Name, &kind, &l.Root, &created, &updated); err != nil {
		return nil, err
	}
	l.ID = shared.UUID(id)
	l.Kind = mediafile.Kind(kind)
	l.CreatedAt, l.UpdatedAt = fromNanos(created), fromNanos(updated)
	return &l, nil
}
