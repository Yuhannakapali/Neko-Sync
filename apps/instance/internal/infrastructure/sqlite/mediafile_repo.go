package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"nekosync-instance/internal/domain/mediafile"
	"nekosync-instance/internal/domain/shared"
)

// MediaFileRepository implements mediafile.Repository.
type MediaFileRepository struct{ db *sql.DB }

// NewMediaFileRepository builds a MediaFileRepository.
func NewMediaFileRepository(db *sql.DB) *MediaFileRepository { return &MediaFileRepository{db: db} }

const mediaFileCols = `id, library_id, path, size, mod_time, kind, container, duration, streams,
	parsed_title, parsed_year, season, episode, provider_ids, match_status, created_at, updated_at`

func (r *MediaFileRepository) Get(ctx context.Context, id shared.UUID) (*mediafile.MediaFile, error) {
	return r.one(ctx, `SELECT `+mediaFileCols+` FROM media_files WHERE id = ?`, string(id))
}

func (r *MediaFileRepository) GetByPath(ctx context.Context, path string) (*mediafile.MediaFile, error) {
	return r.one(ctx, `SELECT `+mediaFileCols+` FROM media_files WHERE path = ?`, path)
}

func (r *MediaFileRepository) ListByLibrary(ctx context.Context, libraryID shared.UUID) ([]*mediafile.MediaFile, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+mediaFileCols+` FROM media_files
		WHERE library_id = ? ORDER BY parsed_title, season, episode, path`, string(libraryID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*mediafile.MediaFile
	for rows.Next() {
		f, err := scanMediaFile(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *MediaFileRepository) Save(ctx context.Context, f *mediafile.MediaFile) error {
	streams, err := toJSON(nonNilStreams(f.Streams))
	if err != nil {
		return err
	}
	ids := f.ProviderIDs
	if ids == nil {
		ids = map[string]string{}
	}
	providerIDs, err := toJSON(ids)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO media_files (`+mediaFileCols+`)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
			library_id = excluded.library_id, path = excluded.path, size = excluded.size,
			mod_time = excluded.mod_time, kind = excluded.kind, container = excluded.container,
			duration = excluded.duration, streams = excluded.streams,
			parsed_title = excluded.parsed_title, parsed_year = excluded.parsed_year,
			season = excluded.season, episode = excluded.episode,
			provider_ids = excluded.provider_ids, match_status = excluded.match_status,
			updated_at = excluded.updated_at`,
		string(f.ID), string(f.LibraryID), f.Path, f.Size, toNanos(f.ModTime), string(f.Kind),
		f.Container, f.Duration, streams,
		f.ParsedTitle, nullInt(f.ParsedYear), nullInt(f.Season), nullFloat(f.Episode),
		providerIDs, string(f.MatchStatus), toNanos(f.CreatedAt), toNanos(f.UpdatedAt))
	return err
}

func (r *MediaFileRepository) Delete(ctx context.Context, id shared.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM media_files WHERE id = ?`, string(id))
	return err
}

func (r *MediaFileRepository) one(ctx context.Context, q string, arg any) (*mediafile.MediaFile, error) {
	f, err := scanMediaFile(r.db.QueryRowContext(ctx, q, arg))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, mediafile.ErrNotFound
	}
	return f, err
}

func scanMediaFile(s scanner) (*mediafile.MediaFile, error) {
	var (
		f                         mediafile.MediaFile
		id, libID, kind, status   string
		streams, providerIDs      string
		modTime, created, updated int64
		year, season              sql.NullInt64
		episode                   sql.NullFloat64
	)
	err := s.Scan(&id, &libID, &f.Path, &f.Size, &modTime, &kind, &f.Container, &f.Duration,
		&streams, &f.ParsedTitle, &year, &season, &episode, &providerIDs, &status, &created, &updated)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(streams), &f.Streams); err != nil {
		return nil, fmt.Errorf("media file %s: decode streams: %w", id, err)
	}
	if err := json.Unmarshal([]byte(providerIDs), &f.ProviderIDs); err != nil {
		return nil, fmt.Errorf("media file %s: decode provider_ids: %w", id, err)
	}
	f.ID, f.LibraryID = shared.UUID(id), shared.UUID(libID)
	f.Kind, f.MatchStatus = mediafile.Kind(kind), mediafile.MatchStatus(status)
	f.ModTime = fromNanos(modTime)
	f.ParsedYear, f.Season, f.Episode = intFromNull(year), intFromNull(season), floatFromNull(episode)
	f.CreatedAt, f.UpdatedAt = fromNanos(created), fromNanos(updated)
	return &f, nil
}

func nonNilStreams(s []mediafile.Stream) []mediafile.Stream {
	if s == nil {
		return []mediafile.Stream{}
	}
	return s
}
