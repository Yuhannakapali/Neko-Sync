package social

import (
	"context"
	"nekosync/internal/domain/shared"
)

type Repository interface {
	CreateDiscussion(ctx context.Context, d *Discussion) error
	GetDiscussionByID(ctx context.Context, id shared.UUID) (*Discussion, error)
	GetDiscussionsByContentID(ctx context.Context, contentID shared.UUID, limit, offset int) ([]*Discussion, error)
	UpdateDiscussion(ctx context.Context, d *Discussion) error
	DeleteDiscussion(ctx context.Context, id shared.UUID) error

	CreatePost(ctx context.Context, p *DiscussionPost) error
	GetPostsByDiscussionID(ctx context.Context, discussionID shared.UUID) ([]*DiscussionPost, error)
	DeletePost(ctx context.Context, id shared.UUID) error

	AddReaction(ctx context.Context, r *DiscussionReaction) error
	RemoveReaction(ctx context.Context, postID, userID shared.UUID) error

	CreateComment(ctx context.Context, c *Comment) error
	GetCommentsByContentID(ctx context.Context, contentID shared.UUID, limit, offset int) ([]*Comment, error)
	DeleteComment(ctx context.Context, id shared.UUID) error

	CreateReview(ctx context.Context, r *Review) error
	GetReviewsByContentID(ctx context.Context, contentID shared.UUID, limit, offset int) ([]*Review, error)
	UpdateReview(ctx context.Context, r *Review) error
	DeleteReview(ctx context.Context, id shared.UUID) error

	CreateReport(ctx context.Context, r *Report) error
	GetReportsByStatus(ctx context.Context, status ReportStatus, limit, offset int) ([]*Report, error)
	UpdateReportStatus(ctx context.Context, id shared.UUID, status ReportStatus) error
}
