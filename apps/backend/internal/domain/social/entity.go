package social

import (
	"nekosync/internal/domain/shared"
	"time"
)

type ReportStatus string

const (
	ReportPending  ReportStatus = "pending"
	ReportReviewed ReportStatus = "reviewed"
	ReportResolved ReportStatus = "resolved"
)

type Discussion struct {
	shared.BaseEntity
	UserID    shared.UUID  `json:"user_id" db:"user_id"`
	Title     string       `json:"title" db:"title"`
	Content   string       `json:"content" db:"content"`
	ContentID *shared.UUID `json:"content_id" db:"content_id"`
	EpisodeID *shared.UUID `json:"episode_id" db:"episode_id"`
	ChapterID *shared.UUID `json:"chapter_id" db:"chapter_id"`
	IsLocked  bool         `json:"is_locked" db:"is_locked"`
	IsPinned  bool         `json:"is_pinned" db:"is_pinned"`
	Views     int          `json:"views" db:"views"`
}

type DiscussionPost struct {
	shared.BaseEntity
	DiscussionID shared.UUID  `json:"discussion_id" db:"discussion_id"`
	UserID       shared.UUID  `json:"user_id" db:"user_id"`
	ParentID     *shared.UUID `json:"parent_id" db:"parent_id"`
	Content      string       `json:"content" db:"content"`
}

type DiscussionReaction struct {
	shared.BaseEntity
	PostID    shared.UUID `json:"post_id" db:"post_id"`
	UserID    shared.UUID `json:"user_id" db:"user_id"`
	Emoji     string      `json:"emoji" db:"emoji"`
	ReactedAt time.Time   `json:"reacted_at" db:"reacted_at"`
}

type Comment struct {
	shared.BaseEntity
	UserID    shared.UUID  `json:"user_id" db:"user_id"`
	ContentID shared.UUID  `json:"content_id" db:"content_id"`
	Comment   string       `json:"comment" db:"comment"`
	ParentID  *shared.UUID `json:"parent_id" db:"parent_id"`
}

type Review struct {
	shared.BaseEntity
	UserID    shared.UUID `json:"user_id" db:"user_id"`
	ContentID shared.UUID `json:"content_id" db:"content_id"`
	Rating    int         `json:"rating" db:"rating"`
	Review    *string     `json:"review" db:"review"`
}

type Report struct {
	shared.BaseEntity
	ReporterID     shared.UUID  `json:"reporter_id" db:"reporter_id"`
	ReportedUserID *shared.UUID `json:"reported_user_id" db:"reported_user_id"`
	ContentType    string       `json:"content_type" db:"content_type"`
	ContentID      shared.UUID  `json:"content_id" db:"content_id"`
	Reason         string       `json:"reason" db:"reason"`
	Status         ReportStatus `json:"status" db:"status"`
}
