package history

import (
	"nekosync/internal/platform/entity"
	"nekosync/internal/progress"
	"time"
)

type Status string

const (
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusDropped    Status = "dropped"
	StatusPlanned    Status = "planned"
)

// Entry is one user's progress through one work, or one child of it
// (episode, chapter, track). It replaces the old per-medium watch/read/listen
// histories: the medium lives on the Work, the position in Progress.
type Entry struct {
	UserID    entity.UUID       `json:"user_id" db:"user_id"`
	WorkID    entity.UUID       `json:"work_id" db:"work_id"`
	ChildID   *entity.UUID      `json:"child_id" db:"child_id"`
	Progress  progress.Progress `json:"progress"`
	Status    Status            `json:"status" db:"status"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
}

type Favorite struct {
	UserID      entity.UUID `json:"user_id" db:"user_id"`
	WorkID      entity.UUID `json:"work_id" db:"work_id"`
	FavoritedAt time.Time   `json:"favorited_at" db:"favorited_at"`
}

// Playlist is an ordered, user-curated list of works or children of any medium.
type Playlist struct {
	entity.BaseEntity
	UserID        entity.UUID `json:"user_id" db:"user_id"`
	Name          string      `json:"name" db:"name"`
	Description   *string     `json:"description" db:"description"`
	CoverImageURL *string     `json:"cover_image_url" db:"cover_image_url"`
	IsPublic      bool        `json:"is_public" db:"is_public"`
}

type PlaylistItem struct {
	PlaylistID entity.UUID  `json:"playlist_id" db:"playlist_id"`
	WorkID     entity.UUID  `json:"work_id" db:"work_id"`
	ChildID    *entity.UUID `json:"child_id" db:"child_id"`
	AddedAt    time.Time    `json:"added_at" db:"added_at"`
	OrderIndex int          `json:"order_index" db:"order_index"`
}
