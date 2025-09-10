package contentgenre

// ContentGenre represents the association between content and genre.
type ContentGenre struct {
	ContentID string `json:"content_id" db:"content_id"`
	GenreID   string `json:"genre_id" db:"genre_id"`
}
