package history

import (
	"nekosync/internal/domain/shared"
	"time"
)

type WatchHistory struct {
	UserID          shared.UUID `json:"user_id" db:"user_id"`
	ContentID       shared.UUID `json:"content_id" db:"content_id"`
	EpisodeID       shared.UUID `json:"episode_id" db:"episode_id"`
	ProgressSeconds int         `json:"progress_seconds" db:"progress_seconds"`
	Status          string      `json:"status" db:"status"`
	WatchedAt       time.Time   `json:"watched_at" db:"watched_at"`
	Finished        bool        `json:"finished" db:"finished"`
}

type ReadHistory struct {
	UserID     shared.UUID `json:"user_id" db:"user_id"`
	ContentID  shared.UUID `json:"content_id" db:"content_id"`
	ChapterID  shared.UUID `json:"chapter_id" db:"chapter_id"`
	PageNumber int         `json:"page_number" db:"page_number"`
	Status     string      `json:"status" db:"status"`
	ReadAt     time.Time   `json:"read_at" db:"read_at"`
	Finished   bool        `json:"finished" db:"finished"`
}

type ListenHistory struct {
	UserID          shared.UUID `json:"user_id" db:"user_id"`
	MusicID         shared.UUID `json:"music_id" db:"music_id"`
	ProgressSeconds int         `json:"progress_seconds" db:"progress_seconds"`
	Status          string      `json:"status" db:"status"`
	ListenedAt      time.Time   `json:"listened_at" db:"listened_at"`
	Finished        bool        `json:"finished" db:"finished"`
}

type Favorite struct {
	UserID      shared.UUID `json:"user_id" db:"user_id"`
	ContentID   shared.UUID `json:"content_id" db:"content_id"`
	FavoritedAt time.Time   `json:"favorited_at" db:"favorited_at"`
}

type FavoriteMusic struct {
	UserID      shared.UUID `json:"user_id" db:"user_id"`
	MusicID     shared.UUID `json:"music_id" db:"music_id"`
	FavoritedAt time.Time   `json:"favorited_at" db:"favorited_at"`
}

type Playlist struct {
	shared.BaseEntity
	UserID        shared.UUID `json:"user_id" db:"user_id"`
	Name          string      `json:"name" db:"name"`
	Description   *string     `json:"description" db:"description"`
	CoverImageURL *string     `json:"cover_image_url" db:"cover_image_url"`
	IsPublic      bool        `json:"is_public" db:"is_public"`
}

type PlaylistMusic struct {
	PlaylistID shared.UUID `json:"playlist_id" db:"playlist_id"`
	MusicID    shared.UUID `json:"music_id" db:"music_id"`
	AddedAt    time.Time   `json:"added_at" db:"added_at"`
	OrderIndex int         `json:"order_index" db:"order_index"`
}
