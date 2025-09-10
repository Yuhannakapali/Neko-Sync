package collaborativefiltering

// CollaborativeFilteringRecommendation represents a recommendation based on collaborative filtering.
type CollaborativeFilteringRecommendation struct {
	ID              string `json:"id" db:"id"`
	UserID          string `json:"user_id" db:"user_id"`
	ContentID       string `json:"content_id" db:"content_id"`
	SimilarityScore string `json:"similarity_score" db:"similarity_score"`
	CreatedAt       string `json:"created_at" db:"created_at"`
	UpdatedAt       string `json:"updated_at" db:"updated_at"`
}
