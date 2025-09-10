package contentrecommendation

// ContentRecommendation represents a content recommendation for a user.
type ContentRecommendation struct {
	ID                   string `json:"id" db:"id"`
	UserID               string `json:"user_id" db:"user_id"`
	ContentID            string `json:"content_id" db:"content_id"`
	RecommendationReason string `json:"recommendation_reason" db:"recommendation_reason"`
	RecommendationType   string `json:"recommendation_type" db:"recommendation_type"`
}
