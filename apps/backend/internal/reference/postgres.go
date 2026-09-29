package reference

import (
	"context"
	"database/sql"
	"errors"

	"nekosync/internal/platform/entity"
)

var ErrNotFound = errors.New("content reference not found")

type pgRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

const columns = `id, user_id, work_id, child_id, source, locator, quality, created_at, updated_at`

func scan(row interface{ Scan(...any) error }) (*ContentReference, error) {
	r := &ContentReference{}
	err := row.Scan(&r.ID, &r.UserID, &r.WorkID, &r.ChildID, &r.Source, &r.Locator, &r.Quality, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func (p *pgRepository) Create(ctx context.Context, r *ContentReference) error {
	_, err := p.db.ExecContext(ctx, `
		INSERT INTO content_references (`+columns+`)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		r.ID, r.UserID, r.WorkID, r.ChildID, r.Source, r.Locator, r.Quality, r.CreatedAt, r.UpdatedAt)
	return err
}

// Upsert inserts r, or — when the user already has this exact source+locator
// for the same work/child — refreshes its quality. r.ID is set to the stored
// row's ID either way, so re-scans keep stable reference IDs.
func (p *pgRepository) Upsert(ctx context.Context, r *ContentReference) error {
	return p.db.QueryRowContext(ctx, `
		INSERT INTO content_references (`+columns+`)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id, work_id, child_id, source, locator)
		DO UPDATE SET quality = EXCLUDED.quality
		RETURNING id, created_at`,
		r.ID, r.UserID, r.WorkID, r.ChildID, r.Source, r.Locator, r.Quality, r.CreatedAt, r.UpdatedAt,
	).Scan(&r.ID, &r.CreatedAt)
}

func (p *pgRepository) Update(ctx context.Context, r *ContentReference) error {
	res, err := p.db.ExecContext(ctx, `
		UPDATE content_references SET child_id = $2, source = $3, locator = $4, quality = $5
		WHERE id = $1`,
		r.ID, r.ChildID, r.Source, r.Locator, r.Quality)
	return affected(res, err)
}

func (p *pgRepository) Delete(ctx context.Context, id entity.UUID) error {
	res, err := p.db.ExecContext(ctx, `DELETE FROM content_references WHERE id = $1`, id)
	return affected(res, err)
}

// Resolve returns the user's references for one work (childID nil) or one
// child. Order is insertion order; callers apply Rank for playback preference.
func (p *pgRepository) Resolve(ctx context.Context, userID, workID entity.UUID, childID *entity.UUID) ([]*ContentReference, error) {
	return p.query(ctx, `SELECT `+columns+` FROM content_references
		WHERE user_id = $1 AND work_id = $2 AND child_id IS NOT DISTINCT FROM $3
		ORDER BY created_at, id`, userID, workID, childID)
}

func (p *pgRepository) ListByUser(ctx context.Context, userID entity.UUID, limit, offset int) ([]*ContentReference, error) {
	return p.query(ctx, `SELECT `+columns+` FROM content_references
		WHERE user_id = $1 ORDER BY created_at, id LIMIT $2 OFFSET $3`, userID, limit, offset)
}

func (p *pgRepository) query(ctx context.Context, q string, args ...any) ([]*ContentReference, error) {
	rows, err := p.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []*ContentReference
	for rows.Next() {
		r, err := scan(rows)
		if err != nil {
			return nil, err
		}
		refs = append(refs, r)
	}
	return refs, rows.Err()
}

func affected(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
