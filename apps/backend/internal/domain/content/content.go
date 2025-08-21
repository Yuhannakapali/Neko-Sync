package content

import (
	"nekosync/internal/domain"
	"time"
)

// ========== SERIES / FRANCHISE ==========

type Series struct {
	domain.BaseEntity
	Title       string             `json:"title" db:"title"`
	Description *string            `json:"description" db:"description"`
	CoverURL    *string            `json:"cover_url" db:"cover_url"`
	Type        domain.ContentType `json:"type" db:"type"`
}

// ========== BASE CONTENT ==========

type Content struct {
	domain.BaseEntity
	Title        string             `json:"title" db:"title"`
	Description  *string            `json:"description" db:"description"`
	ThumbnailURL *string            `json:"thumbnail_url" db:"thumbnail_url"`
	CoverPicture *string            `json:"cover_picture" db:"cover_picture"`
	Type         domain.ContentType `json:"type" db:"type"`
	SeriesID     *domain.UUID       `json:"series_id" db:"series_id"`
}

type Genre struct {
	domain.BaseEntity
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
}

type ContentGenre struct {
	ContentID domain.UUID `json:"content_id" db:"content_id"`
	GenreID   domain.UUID `json:"genre_id" db:"genre_id"`
}

type Tag struct {
	domain.BaseEntity
	Name string `json:"name" db:"name"`
}

type ContentTag struct {
	ContentID domain.UUID `json:"content_id" db:"content_id"`
	TagID     domain.UUID `json:"tag_id" db:"tag_id"`
}

// ========== ANIME ==========

type Anime struct {
	domain.BaseEntity
	ContentID        domain.UUID          `json:"content_id" db:"id"`
	NumberOfEpisodes *int                 `json:"number_of_episodes" db:"number_of_episodes"`
	Season           *domain.Season       `json:"season" db:"season"`
	Year             *int                 `json:"year" db:"year"`
	Status           domain.ContentStatus `json:"status" db:"status"`
	Rating           *float64             `json:"rating" db:"rating"`
}

type Episode struct {
	domain.BaseEntity
	ContentID          domain.UUID `json:"content_id" db:"content_id"`
	EpisodeNumber      int         `json:"episode_number" db:"episode_number"`
	Title              string      `json:"title" db:"title"`
	Description        *string     `json:"description" db:"description"`
	VideoURL           *string     `json:"video_url" db:"video_url"`
	Duration           *int        `json:"duration" db:"duration"`
	ReleaseDate        *time.Time  `json:"release_date" db:"release_date"`
	Subtitles          *string     `json:"subtitles" db:"subtitles"`                     // JSON
	LanguagesSupported *string     `json:"languages_supported" db:"languages_supported"` // JSON
	Views              int         `json:"views" db:"views"`
}

// ========== MANGA ==========

type Manga struct {
	domain.BaseEntity
	ContentID        domain.UUID          `json:"content_id" db:"id"`
	NumberOfChapters *int                 `json:"number_of_chapters" db:"number_of_chapters"`
	Status           domain.ContentStatus `json:"status" db:"status"`
	Rating           *float64             `json:"rating" db:"rating"`
}

type Chapter struct {
	domain.BaseEntity
	ContentID     domain.UUID `json:"content_id" db:"content_id"`
	ChapterNumber int         `json:"chapter_number" db:"chapter_number"`
	Title         string      `json:"title" db:"title"`
	Description   *string     `json:"description" db:"description"`
	Pages         *string     `json:"pages" db:"pages"` // JSON
	ReleaseDate   *time.Time  `json:"release_date" db:"release_date"`
	Views         int         `json:"views" db:"views"`
}

// ========== MOVIES ==========

type Movie struct {
	domain.BaseEntity
	ContentID   domain.UUID `json:"content_id" db:"id"`
	Duration    *int        `json:"duration" db:"duration"`
	VideoURL    *string     `json:"video_url" db:"video_url"`
	ReleaseDate *time.Time  `json:"release_date" db:"release_date"`
	Rating      *float64    `json:"rating" db:"rating"`
	Views       int         `json:"views" db:"views"`
}

// ========== MUSIC ==========

type Music struct {
	domain.BaseEntity
	ContentID   domain.UUID      `json:"content_id" db:"id"`
	MusicType   domain.MusicType `json:"music_type" db:"music_type"`
	AudioURL    *string          `json:"audio_url" db:"audio_url"`
	Duration    *int             `json:"duration" db:"duration"`
	ReleaseDate *time.Time       `json:"release_date" db:"release_date"`
	Plays       int              `json:"plays" db:"plays"`
}

type Artist struct {
	domain.BaseEntity
	Name      string  `json:"name" db:"name"`
	Bio       *string `json:"bio" db:"bio"`
	AvatarURL *string `json:"avatar_url" db:"avatar_url"`
}

type MusicArtist struct {
	MusicID  domain.UUID `json:"music_id" db:"music_id"`
	ArtistID domain.UUID `json:"artist_id" db:"artist_id"`
	Role     *string     `json:"role" db:"role"`
}
