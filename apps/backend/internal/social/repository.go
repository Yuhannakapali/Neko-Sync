package social

import (
	"context"
	"nekosync/internal/platform/entity"
)

type Repository interface {
	CreateDiscussion(ctx context.Context, d *Discussion) error
	GetDiscussionByID(ctx context.Context, id entity.UUID) (*Discussion, error)
	GetDiscussionsByContentID(ctx context.Context, contentID entity.UUID, limit, offset int) ([]*Discussion, error)
	UpdateDiscussion(ctx context.Context, d *Discussion) error
	DeleteDiscussion(ctx context.Context, id entity.UUID) error

	CreatePost(ctx context.Context, p *DiscussionPost) error
	GetPostsByDiscussionID(ctx context.Context, discussionID entity.UUID) ([]*DiscussionPost, error)
	DeletePost(ctx context.Context, id entity.UUID) error

	AddReaction(ctx context.Context, r *DiscussionReaction) error
	RemoveReaction(ctx context.Context, postID, userID entity.UUID) error

	CreateComment(ctx context.Context, c *Comment) error
	GetCommentsByContentID(ctx context.Context, contentID entity.UUID, limit, offset int) ([]*Comment, error)
	DeleteComment(ctx context.Context, id entity.UUID) error

	CreateReview(ctx context.Context, r *Review) error
	GetReviewsByContentID(ctx context.Context, contentID entity.UUID, limit, offset int) ([]*Review, error)
	UpdateReview(ctx context.Context, r *Review) error
	DeleteReview(ctx context.Context, id entity.UUID) error

	CreateReport(ctx context.Context, r *Report) error
	GetReportsByStatus(ctx context.Context, status ReportStatus, limit, offset int) ([]*Report, error)
	UpdateReportStatus(ctx context.Context, id entity.UUID, status ReportStatus) error
}
