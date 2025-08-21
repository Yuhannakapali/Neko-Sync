package entities

import "time"

// ========== SERIES / FRANCHISE ==========

// Series represents a series or franchise that can contain multiple content items
type Series struct {
	BaseEntity
	Title       string      `json:"title" db:"title"`
	Description *string     `json:"description" db:"description"`
	CoverURL    *string     `json:"cover_url" db:"cover_url"`
	Type        ContentType `json:"type" db:"type"`
}

// ========== BASE CONTENT ==========

// Content represents the base content entity
type Content struct {
	BaseEntity
	Title        string      `json:"title" db:"title"`
	Description  *string     `json:"description" db:"description"`
	ThumbnailURL *string     `json:"thumbnail_url" db:"thumbnail_url"`
	CoverPicture *string     `json:"cover_picture" db:"cover_picture"`
	Type         ContentType `json:"type" db:"type"`
	SeriesID     *UUID       `json:"series_id" db:"series_id"`
}

// Genre represents a content genre
type Genre struct {
	BaseEntity
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
}

// ContentGenre represents the many-to-many relationship between content and genres
type ContentGenre struct {
	ContentID UUID `json:"content_id" db:"content_id"`
	GenreID   UUID `json:"genre_id" db:"genre_id"`
}

// Tag represents a content tag
type Tag struct {
	BaseEntity
	Name string `json:"name" db:"name"`
}

// ContentTag represents the many-to-many relationship between content and tags
type ContentTag struct {
	ContentID UUID `json:"content_id" db:"content_id"`
	TagID     UUID `json:"tag_id" db:"tag_id"`
}

// ========== ANIME ==========

// Anime represents anime-specific content
type Anime struct {
	BaseEntity
	ContentID        UUID          `json:"content_id" db:"content_id"`
	NumberOfEpisodes *int          `json:"number_of_episodes" db:"number_of_episodes"`
	Season           *Season       `json:"season" db:"season"`
	Year             *int          `json:"year" db:"year"`
	Status           ContentStatus `json:"status" db:"status"`
	Rating           *float64      `json:"rating" db:"rating"`
}

// Episode represents an anime episode
type Episode struct {
	BaseEntity
	ContentID          UUID       `json:"content_id" db:"content_id"`
	EpisodeNumber      int        `json:"episode_number" db:"episode_number"`
	Title              string     `json:"title" db:"title"`
	Description        *string    `json:"description" db:"description"`
	VideoURL           *string    `json:"video_url" db:"video_url"`
	Duration           *int       `json:"duration" db:"duration"`
	ReleaseDate        *time.Time `json:"release_date" db:"release_date"`
	Subtitles          []string   `json:"subtitles" db:"subtitles"`                     // JSON
	LanguagesSupported []string   `json:"languages_supported" db:"languages_supported"` // JSON
	Views              int        `json:"views" db:"views"`
}

// ========== MANGA ==========

// Manga represents manga-specific content
type Manga struct {
	BaseEntity
	ContentID        UUID          `json:"content_id" db:"content_id"`
	NumberOfChapters *int          `json:"number_of_chapters" db:"number_of_chapters"`
	Status           ContentStatus `json:"status" db:"status"`
	Rating           *float64      `json:"rating" db:"rating"`
}

// Chapter represents a manga chapter
type Chapter struct {
	BaseEntity
	ContentID     UUID       `json:"content_id" db:"content_id"`
	ChapterNumber int        `json:"chapter_number" db:"chapter_number"`
	Title         string     `json:"title" db:"title"`
	Description   *string    `json:"description" db:"description"`
	Pages         []string   `json:"pages" db:"pages"` // JSON array of page URLs
	ReleaseDate   *time.Time `json:"release_date" db:"release_date"`
	Views         int        `json:"views" db:"views"`
}

// ========== MOVIES ==========

// Movie represents movie-specific content
type Movie struct {
	BaseEntity
	ContentID   UUID       `json:"content_id" db:"content_id"`
	Duration    *int       `json:"duration" db:"duration"`
	VideoURL    *string    `json:"video_url" db:"video_url"`
	ReleaseDate *time.Time `json:"release_date" db:"release_date"`
	Rating      *float64   `json:"rating" db:"rating"`
	Views       int        `json:"views" db:"views"`
}

// ========== MUSIC ==========

// Music represents music-specific content
type Music struct {
	BaseEntity
	ContentID   UUID       `json:"content_id" db:"content_id"`
	MusicType   MusicType  `json:"music_type" db:"music_type"`
	AudioURL    *string    `json:"audio_url" db:"audio_url"`
	Duration    *int       `json:"duration" db:"duration"`
	ReleaseDate *time.Time `json:"release_date" db:"release_date"`
	Plays       int        `json:"plays" db:"plays"`
}

// Artist represents a music artist
type Artist struct {
	BaseEntity
	Name      string  `json:"name" db:"name"`
	Bio       *string `json:"bio" db:"bio"`
	AvatarURL *string `json:"avatar_url" db:"avatar_url"`
}

// MusicArtist represents the relationship between music and artists
type MusicArtist struct {
	MusicID  UUID    `json:"music_id" db:"music_id"`
	ArtistID UUID    `json:"artist_id" db:"artist_id"`
	Role     *string `json:"role" db:"role"` // e.g., "composer", "performer", "vocalist"
}

// Business Rules for Content domain
func (a *Anime) IsCompleted() bool {
	return a.Status == ContentStatusCompleted
}

func (m *Manga) IsCompleted() bool {
	return m.Status == ContentStatusCompleted
}

func (e *Episode) IncrementViews() {
	e.Views++
	e.UpdatedAt = time.Now()
}

func (c *Chapter) IncrementViews() {
	c.Views++
	c.UpdatedAt = time.Now()
}

func (mv *Movie) IncrementViews() {
	mv.Views++
	mv.UpdatedAt = time.Now()
}

func (m *Music) IncrementPlays() {
	m.Plays++
	m.UpdatedAt = time.Now()
}
