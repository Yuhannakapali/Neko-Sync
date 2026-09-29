package work

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"nekosync/internal/platform/entity"
)

type pgRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

const workColumns = `id, kind, title, titles, year, synopsis, artwork,
	array_to_json(genres), array_to_json(tags), provider_ids, created_at, updated_at`

// jsonArg encodes composite fields for JSONB / text[] columns. Nil slices and
// maps are stored as empty values so reads never return SQL NULL.
func jsonArg(v any, empty string) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	if string(b) == "null" {
		return empty, nil
	}
	return string(b), nil
}

type workArgs struct {
	titles, artwork, genres, tags, providerIDs string
}

func encodeWork(w *Work) (workArgs, error) {
	var a workArgs
	var err error
	if a.titles, err = jsonArg(w.Titles, "[]"); err != nil {
		return a, err
	}
	if a.artwork, err = jsonArg(w.Artwork, "[]"); err != nil {
		return a, err
	}
	if a.genres, err = jsonArg(w.Genres, "[]"); err != nil {
		return a, err
	}
	if a.tags, err = jsonArg(w.Tags, "[]"); err != nil {
		return a, err
	}
	a.providerIDs, err = jsonArg(w.ProviderIDs, "{}")
	return a, err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanWork(row scanner) (*Work, error) {
	w := &Work{}
	var titles, artwork, genres, tags, providerIDs []byte
	if err := row.Scan(&w.ID, &w.Kind, &w.Title, &titles, &w.Year, &w.Synopsis, &artwork,
		&genres, &tags, &providerIDs, &w.CreatedAt, &w.UpdatedAt); err != nil {
		return nil, err
	}
	for _, f := range []struct {
		raw []byte
		dst any
	}{{titles, &w.Titles}, {artwork, &w.Artwork}, {genres, &w.Genres}, {tags, &w.Tags}, {providerIDs, &w.ProviderIDs}} {
		if err := json.Unmarshal(f.raw, f.dst); err != nil {
			return nil, fmt.Errorf("decode work %s: %w", w.ID, err)
		}
	}
	return w, nil
}

func (r *pgRepository) Create(ctx context.Context, w *Work) error {
	a, err := encodeWork(w)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO works (id, kind, title, titles, year, synopsis, artwork, genres, tags, provider_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7,
			ARRAY(SELECT jsonb_array_elements_text($8::jsonb)),
			ARRAY(SELECT jsonb_array_elements_text($9::jsonb)),
			$10, $11, $12)`,
		w.ID, w.Kind, w.Title, a.titles, w.Year, w.Synopsis, a.artwork, a.genres, a.tags, a.providerIDs,
		w.CreatedAt, w.UpdatedAt)
	return err
}

func (r *pgRepository) GetByID(ctx context.Context, id entity.UUID) (*Work, error) {
	w, err := scanWork(r.db.QueryRowContext(ctx, `SELECT `+workColumns+` FROM works WHERE id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return w, err
}

func (r *pgRepository) GetByProviderID(ctx context.Context, provider, id string) (*Work, error) {
	w, err := scanWork(r.db.QueryRowContext(ctx,
		`SELECT `+workColumns+` FROM works WHERE provider_ids @> jsonb_build_object($1::text, $2::text) LIMIT 1`,
		provider, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return w, err
}

func (r *pgRepository) Update(ctx context.Context, w *Work) error {
	a, err := encodeWork(w)
	if err != nil {
		return err
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE works SET kind = $2, title = $3, titles = $4, year = $5, synopsis = $6, artwork = $7,
			genres = ARRAY(SELECT jsonb_array_elements_text($8::jsonb)),
			tags = ARRAY(SELECT jsonb_array_elements_text($9::jsonb)),
			provider_ids = $10
		WHERE id = $1`,
		w.ID, w.Kind, w.Title, a.titles, w.Year, w.Synopsis, a.artwork, a.genres, a.tags, a.providerIDs)
	return affected(res, err, ErrNotFound)
}

func (r *pgRepository) Delete(ctx context.Context, id entity.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM works WHERE id = $1`, id)
	return affected(res, err, ErrNotFound)
}

// List returns works of the given kind, or of every kind when kind is "".
func (r *pgRepository) List(ctx context.Context, kind Kind, limit, offset int) ([]*Work, error) {
	return r.queryWorks(ctx, `SELECT `+workColumns+` FROM works
		WHERE ($1 = '' OR kind = $1) ORDER BY title, id LIMIT $2 OFFSET $3`, kind, limit, offset)
}

// Search matches the primary or any localized title, case-insensitively.
func (r *pgRepository) Search(ctx context.Context, query string, kind Kind, limit, offset int) ([]*Work, error) {
	pattern := "%" + escapeLike(query) + "%"
	return r.queryWorks(ctx, `SELECT `+workColumns+` FROM works
		WHERE ($2 = '' OR kind = $2)
		  AND (title ILIKE $1 OR EXISTS (
		        SELECT 1 FROM jsonb_array_elements(titles) t WHERE t->>'title' ILIKE $1))
		ORDER BY title, id LIMIT $3 OFFSET $4`, pattern, kind, limit, offset)
}

func (r *pgRepository) queryWorks(ctx context.Context, query string, args ...any) ([]*Work, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var works []*Work
	for rows.Next() {
		w, err := scanWork(rows)
		if err != nil {
			return nil, err
		}
		works = append(works, w)
	}
	return works, rows.Err()
}

const childColumns = `id, work_id, parent_id, kind, ordinal, title, synopsis, duration, release_date,
	provider_ids, created_at, updated_at`

func scanChild(row scanner) (*WorkChild, error) {
	c := &WorkChild{}
	var providerIDs []byte
	if err := row.Scan(&c.ID, &c.WorkID, &c.ParentID, &c.Kind, &c.Ordinal, &c.Title, &c.Synopsis,
		&c.Duration, &c.ReleaseDate, &providerIDs, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(providerIDs, &c.ProviderIDs); err != nil {
		return nil, fmt.Errorf("decode work child %s: %w", c.ID, err)
	}
	return c, nil
}

func (r *pgRepository) CreateChild(ctx context.Context, c *WorkChild) error {
	providerIDs, err := jsonArg(c.ProviderIDs, "{}")
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO work_children (id, work_id, parent_id, kind, ordinal, title, synopsis, duration,
			release_date, provider_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		c.ID, c.WorkID, c.ParentID, c.Kind, c.Ordinal, c.Title, c.Synopsis, c.Duration,
		c.ReleaseDate, providerIDs, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *pgRepository) GetChild(ctx context.Context, id entity.UUID) (*WorkChild, error) {
	c, err := scanChild(r.db.QueryRowContext(ctx, `SELECT `+childColumns+` FROM work_children WHERE id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrChildNotFound
	}
	return c, err
}

func (r *pgRepository) ListChildren(ctx context.Context, workID entity.UUID) ([]*WorkChild, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+childColumns+` FROM work_children
		WHERE work_id = $1 ORDER BY parent_id NULLS FIRST, kind, ordinal`, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []*WorkChild
	for rows.Next() {
		c, err := scanChild(rows)
		if err != nil {
			return nil, err
		}
		children = append(children, c)
	}
	return children, rows.Err()
}

func (r *pgRepository) GetChildByOrdinal(ctx context.Context, workID entity.UUID, parentID *entity.UUID, kind ChildKind, ordinal float64) (*WorkChild, error) {
	c, err := scanChild(r.db.QueryRowContext(ctx, `SELECT `+childColumns+` FROM work_children
		WHERE work_id = $1 AND parent_id IS NOT DISTINCT FROM $2 AND kind = $3 AND ordinal = $4`,
		workID, parentID, kind, ordinal))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrChildNotFound
	}
	return c, err
}

func (r *pgRepository) UpdateChild(ctx context.Context, c *WorkChild) error {
	providerIDs, err := jsonArg(c.ProviderIDs, "{}")
	if err != nil {
		return err
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE work_children SET parent_id = $2, kind = $3, ordinal = $4, title = $5, synopsis = $6,
			duration = $7, release_date = $8, provider_ids = $9
		WHERE id = $1`,
		c.ID, c.ParentID, c.Kind, c.Ordinal, c.Title, c.Synopsis, c.Duration, c.ReleaseDate, providerIDs)
	return affected(res, err, ErrChildNotFound)
}

func (r *pgRepository) DeleteChild(ctx context.Context, id entity.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM work_children WHERE id = $1`, id)
	return affected(res, err, ErrChildNotFound)
}

func affected(res sql.Result, err error, notFound error) error {
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return notFound
	}
	return nil
}

func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
