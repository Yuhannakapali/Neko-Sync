package content

import (
	"nekosync/internal/domain/shared"
	"time"
)

type Type string

const (
	TypeAnime Type = "anime"
	TypeManga Type = "manga"
	TypeMovie Type = "movie"
	TypeMusic Type = "music"
)

type Season string

const (
	SeasonWinter Season = "winter"
	SeasonSpring Season = "spring"
	SeasonSummer Season = "summer"
	SeasonFall   Season = "fall"
)

type Status string

const (
	StatusOngoing   Status = "ongoing"
	StatusCompleted Status = "completed"
	StatusUpcoming  Status = "upcoming"
)

type MusicType string

const (
	MusicTypeSong       MusicType = "song"
	MusicTypeAlbum      MusicType = "album"
	MusicTypeSoundtrack MusicType = "soundtrack"
)

// ========== BASE CONTENT ==========

type Content struct {
	shared.BaseEntity
	Title        string  `json:"title" db:"title"`
	Description  *string `json:"description" db:"description"`
	ThumbnailURL *string `json:"thumbnail_url" db:"thumbnail_url"`
	CoverPicture *string `json:"cover_picture" db:"cover_picture"`
	Type         Type    `json:"type" db:"type"`
	SeriesID     *shared.UUID `json:"series_id" db:"series_id"`
}

type Series struct {
	shared.BaseEntity
	Title       string  `json:"title" db:"title"`
	Description *string `json:"description" db:"description"`
	CoverURL    *string `json:"cover_url" db:"cover_url"`
	Type        Type    `json:"type" db:"type"`
}

type Genre struct {
	shared.BaseEntity
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
}

type ContentGenre struct {
	ContentID shared.UUID `json:"content_id" db:"content_id"`
	GenreID   shared.UUID `json:"genre_id" db:"genre_id"`
}

type Tag struct {
	shared.BaseEntity
	Name string `json:"name" db:"name"`
}

type ContentTag struct {
	ContentID shared.UUID `json:"content_id" db:"content_id"`
	TagID     shared.UUID `json:"tag_id" db:"tag_id"`
}

// ========== ANIME ==========

type Anime struct {
	shared.BaseEntity
	ContentID        shared.UUID `json:"content_id" db:"content_id"`
	NumberOfEpisodes *int        `json:"number_of_episodes" db:"number_of_episodes"`
	Season           *Season     `json:"season" db:"season"`
	Year             *int        `json:"year" db:"year"`
	Status           Status      `json:"status" db:"status"`
	Rating           *float64    `json:"rating" db:"rating"`
}

func (a *Anime) IsCompleted() bool { return a.Status == StatusCompleted }

type Episode struct {
	shared.BaseEntity
	ContentID          shared.UUID `json:"content_id" db:"content_id"`
	EpisodeNumber      int         `json:"episode_number" db:"episode_number"`
	Title              string      `json:"title" db:"title"`
	Description        *string     `json:"description" db:"description"`
	VideoURL           *string     `json:"video_url" db:"video_url"`
	Duration           *int        `json:"duration" db:"duration"`
	ReleaseDate        *time.Time  `json:"release_date" db:"release_date"`
	Subtitles          []string    `json:"subtitles" db:"subtitles"`
	LanguagesSupported []string    `json:"languages_supported" db:"languages_supported"`
	Views              int         `json:"views" db:"views"`
}

func (e *Episode) IncrementViews() {
	e.Views++
	e.UpdatedAt = time.Now()
}

// ========== MANGA ==========

type Manga struct {
	shared.BaseEntity
	ContentID        shared.UUID `json:"content_id" db:"content_id"`
	NumberOfChapters *int        `json:"number_of_chapters" db:"number_of_chapters"`
	Status           Status      `json:"status" db:"status"`
	Rating           *float64    `json:"rating" db:"rating"`
}

func (m *Manga) IsCompleted() bool { return m.Status == StatusCompleted }

type Chapter struct {
	shared.BaseEntity
	ContentID     shared.UUID `json:"content_id" db:"content_id"`
	ChapterNumber int         `json:"chapter_number" db:"chapter_number"`
	Title         string      `json:"title" db:"title"`
	Description   *string     `json:"description" db:"description"`
	Pages         []string    `json:"pages" db:"pages"`
	ReleaseDate   *time.Time  `json:"release_date" db:"release_date"`
	Views         int         `json:"views" db:"views"`
}

func (c *Chapter) IncrementViews() {
	c.Views++
	c.UpdatedAt = time.Now()
}

// ========== MOVIE ==========

type Movie struct {
	shared.BaseEntity
	ContentID   shared.UUID `json:"content_id" db:"content_id"`
	Duration    *int        `json:"duration" db:"duration"`
	VideoURL    *string     `json:"video_url" db:"video_url"`
	ReleaseDate *time.Time  `json:"release_date" db:"release_date"`
	Rating      *float64    `json:"rating" db:"rating"`
	Views       int         `json:"views" db:"views"`
}

func (m *Movie) IncrementViews() {
	m.Views++
	m.UpdatedAt = time.Now()
}

// ========== MUSIC ==========

type Music struct {
	shared.BaseEntity
	ContentID   shared.UUID `json:"content_id" db:"content_id"`
	MusicType   MusicType   `json:"music_type" db:"music_type"`
	AudioURL    *string     `json:"audio_url" db:"audio_url"`
	Duration    *int        `json:"duration" db:"duration"`
	ReleaseDate *time.Time  `json:"release_date" db:"release_date"`
	Plays       int         `json:"plays" db:"plays"`
}

func (m *Music) IncrementPlays() {
	m.Plays++
	m.UpdatedAt = time.Now()
}

type Artist struct {
	shared.BaseEntity
	Name      string  `json:"name" db:"name"`
	Bio       *string `json:"bio" db:"bio"`
	AvatarURL *string `json:"avatar_url" db:"avatar_url"`
}

type MusicArtist struct {
	MusicID  shared.UUID `json:"music_id" db:"music_id"`
	ArtistID shared.UUID `json:"artist_id" db:"artist_id"`
	Role     *string     `json:"role" db:"role"`
}
