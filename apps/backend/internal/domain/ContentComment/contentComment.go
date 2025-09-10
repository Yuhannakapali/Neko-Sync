package contentcomment

// ContentComment represents a comment on content.
type ContentComment struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"user_id" db:"user_id"`
	ContentID string `json:"content_id" db:"content_id"`
	Comment   string `json:"comment" db:"comment"`
	ParentID  string `json:"parent_id" db:"parent_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}
