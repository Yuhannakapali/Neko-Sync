package entities

import "time"

type ReportStatus string

const (
	ReportStatusPending  ReportStatus = "pending"
	ReportStatusReviewed ReportStatus = "reviewed"
	ReportStatusResolved ReportStatus = "resolved"
)

type Discussion struct {
	BaseEntity
	UserID    UUID    `json:"user_id" db:"user_id"`
	Title     string  `json:"title" db:"title"`
	Content   string  `json:"content" db:"content"`
	ContentID *UUID   `json:"content_id" db:"content_id"`
	EpisodeID *UUID   `json:"episode_id" db:"episode_id"`
	ChapterID *UUID   `json:"chapter_id" db:"chapter_id"`
	IsLocked  bool    `json:"is_locked" db:"is_locked"`
	IsPinned  bool    `json:"is_pinned" db:"is_pinned"`
	Views     int     `json:"views" db:"views"`
}

type DiscussionPost struct {
	BaseEntity
	DiscussionID UUID   `json:"discussion_id" db:"discussion_id"`
	UserID       UUID   `json:"user_id" db:"user_id"`
	ParentID     *UUID  `json:"parent_id" db:"parent_id"`
	Content      string `json:"content" db:"content"`
}

type DiscussionReaction struct {
	BaseEntity
	PostID    UUID      `json:"post_id" db:"post_id"`
	UserID    UUID      `json:"user_id" db:"user_id"`
	Emoji     string    `json:"emoji" db:"emoji"`
	ReactedAt time.Time `json:"reacted_at" db:"reacted_at"`
}

type ContentComment struct {
	BaseEntity
	UserID    UUID    `json:"user_id" db:"user_id"`
	ContentID UUID    `json:"content_id" db:"content_id"`
	Comment   string  `json:"comment" db:"comment"`
	ParentID  *UUID   `json:"parent_id" db:"parent_id"`
}

type ContentReview struct {
	BaseEntity
	UserID    UUID    `json:"user_id" db:"user_id"`
	ContentID UUID    `json:"content_id" db:"content_id"`
	Rating    int     `json:"rating" db:"rating"`
	Review    *string `json:"review" db:"review"`
}

type Report struct {
	BaseEntity
	ReporterID     UUID         `json:"reporter_id" db:"reporter_id"`
	ReportedUserID *UUID        `json:"reported_user_id" db:"reported_user_id"`
	ContentType    string       `json:"content_type" db:"content_type"`
	ContentID      UUID         `json:"content_id" db:"content_id"`
	Reason         string       `json:"reason" db:"reason"`
	Status         ReportStatus `json:"status" db:"status"`
}
