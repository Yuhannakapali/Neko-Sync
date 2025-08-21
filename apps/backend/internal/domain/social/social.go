package social

import (
	"nekosync/internal/domain"
	"time"
)

// ========== DISCUSSIONS / FORUM ==========

type Discussion struct {
	domain.BaseEntity
	UserID    domain.UUID  `json:"user_id" db:"user_id"`
	Title     string       `json:"title" db:"title"`
	Content   string       `json:"content" db:"content"`
	ContentID *domain.UUID `json:"content_id" db:"content_id"`
	EpisodeID *domain.UUID `json:"episode_id" db:"episode_id"`
	ChapterID *domain.UUID `json:"chapter_id" db:"chapter_id"`
	IsLocked  bool         `json:"is_locked" db:"is_locked"`
	IsPinned  bool         `json:"is_pinned" db:"is_pinned"`
	Views     int          `json:"views" db:"views"`
}

type DiscussionPost struct {
	domain.BaseEntity
	DiscussionID domain.UUID  `json:"discussion_id" db:"discussion_id"`
	UserID       domain.UUID  `json:"user_id" db:"user_id"`
	ParentID     *domain.UUID `json:"parent_id" db:"parent_id"`
	Content      string       `json:"content" db:"content"`
}

type DiscussionReaction struct {
	domain.BaseEntity
	PostID    domain.UUID `json:"post_id" db:"post_id"`
	UserID    domain.UUID `json:"user_id" db:"user_id"`
	Emoji     string      `json:"emoji" db:"emoji"`
	ReactedAt time.Time   `json:"reacted_at" db:"reacted_at"`
}

// ========== COMMENTS & REVIEWS ==========

type ContentComment struct {
	domain.BaseEntity
	UserID    domain.UUID  `json:"user_id" db:"user_id"`
	ContentID domain.UUID  `json:"content_id" db:"content_id"`
	Comment   string       `json:"comment" db:"comment"`
	ParentID  *domain.UUID `json:"parent_id" db:"parent_id"`
}

type ContentReview struct {
	domain.BaseEntity
	UserID    domain.UUID `json:"user_id" db:"user_id"`
	ContentID domain.UUID `json:"content_id" db:"content_id"`
	Rating    int         `json:"rating" db:"rating"`
	Review    *string     `json:"review" db:"review"`
}

// ========== REPORTS & MODERATION ==========

type Report struct {
	domain.BaseEntity
	ReporterID     domain.UUID         `json:"reporter_id" db:"reporter_id"`
	ReportedUserID *domain.UUID        `json:"reported_user_id" db:"reported_user_id"`
	ContentType    string              `json:"content_type" db:"content_type"`
	ContentID      domain.UUID         `json:"content_id" db:"content_id"`
	Reason         string              `json:"reason" db:"reason"`
	Status         domain.ReportStatus `json:"status" db:"status"`
}
