package chapter

// Chapter represents a chapter of content (e.g., manga chapter).
type Chapter struct {
	ID            string `json:"id" db:"id"`
	ContentID     string `json:"content_id" db:"content_id"`
	ChapterNumber int    `json:"chapter_number" db:"chapter_number"`
	Title         string `json:"title" db:"title"`
	Description   string `json:"description" db:"description"`
	Pages         string `json:"pages" db:"pages"`
	ReleaseDate   string `json:"release_date" db:"release_date"`
	View          int    `json:"view" db:"view"`
	CreatedAt     string `json:"created_at" db:"created_at"`
	UpdatedAt     string `json:"updated_at" db:"updated_at"`
}
