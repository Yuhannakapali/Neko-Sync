package social

import "nekosync/internal/domain"

type Repository interface {
	// Discussion operations
	CreateDiscussion(discussion *Discussion) error
	GetDiscussion(id domain.UUID) (*Discussion, error)
	GetDiscussionsByContent(contentID domain.UUID) ([]*Discussion, error)
	GetDiscussionsByUser(userID domain.UUID) ([]*Discussion, error)
	UpdateDiscussion(discussion *Discussion) error
	DeleteDiscussion(id domain.UUID) error
	IncrementDiscussionViews(id domain.UUID) error

	// Discussion Post operations
	CreateDiscussionPost(post *DiscussionPost) error
	GetDiscussionPost(id domain.UUID) (*DiscussionPost, error)
	GetDiscussionPosts(discussionID domain.UUID) ([]*DiscussionPost, error)
	GetPostReplies(parentID domain.UUID) ([]*DiscussionPost, error)
	UpdateDiscussionPost(post *DiscussionPost) error
	DeleteDiscussionPost(id domain.UUID) error

	// Discussion Reaction operations
	CreateDiscussionReaction(reaction *DiscussionReaction) error
	DeleteDiscussionReaction(postID, userID domain.UUID, emoji string) error
	GetPostReactions(postID domain.UUID) ([]*DiscussionReaction, error)
	GetUserReaction(postID, userID domain.UUID) (*DiscussionReaction, error)

	// Content Comment operations
	CreateContentComment(comment *ContentComment) error
	GetContentComment(id domain.UUID) (*ContentComment, error)
	GetContentComments(contentID domain.UUID) ([]*ContentComment, error)
	GetCommentReplies(parentID domain.UUID) ([]*ContentComment, error)
	UpdateContentComment(comment *ContentComment) error
	DeleteContentComment(id domain.UUID) error

	// Content Review operations
	CreateContentReview(review *ContentReview) error
	GetContentReview(id domain.UUID) (*ContentReview, error)
	GetContentReviews(contentID domain.UUID) ([]*ContentReview, error)
	GetUserReview(userID, contentID domain.UUID) (*ContentReview, error)
	UpdateContentReview(review *ContentReview) error
	DeleteContentReview(id domain.UUID) error
	GetAverageRating(contentID domain.UUID) (float64, error)

	// Report operations
	CreateReport(report *Report) error
	GetReport(id domain.UUID) (*Report, error)
	GetReportsByStatus(status domain.ReportStatus) ([]*Report, error)
	GetReportsByReporter(reporterID domain.UUID) ([]*Report, error)
	UpdateReportStatus(id domain.UUID, status domain.ReportStatus) error
	DeleteReport(id domain.UUID) error
}
