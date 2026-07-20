// Package work holds the canonical metadata identity for a title — the "what is
// this" that is shared across all users and legal to centralize on the Hub.
//
// A Work never carries media bytes or links to them: no video/audio URLs, no page
// lists. Where a given user's actual playable copy lives is expressed separately by
// domain/reference.ContentReference, and the bytes themselves live only on a
// self-hosted Instance (domain not present in this layer).
package work

import (
	"time"

	"nekosync/internal/domain/shared"
)

// Kind is the top-level medium of a Work.
type Kind string

const (
	KindAnime  Kind = "anime"
	KindManga  Kind = "manga"
	KindMovie  Kind = "movie"
	KindSeries Kind = "series"
	KindMusic  Kind = "music"
	KindBook   Kind = "book"
)

// ChildKind identifies a structural child of a Work.
type ChildKind string

const (
	ChildSeason  ChildKind = "season"
	ChildEpisode ChildKind = "episode"
	ChildVolume  ChildKind = "volume"
	ChildChapter ChildKind = "chapter"
	ChildTrack   ChildKind = "track"
)

// LocalizedTitle is a title in a specific language.
type LocalizedTitle struct {
	Lang  string `json:"lang" db:"lang"` // BCP-47, e.g. "en", "ja", "ja-Latn"
	Title string `json:"title" db:"title"`
}

// ArtworkKind classifies a metadata image (poster art, not media bytes).
type ArtworkKind string

const (
	ArtworkPoster   ArtworkKind = "poster"
	ArtworkBackdrop ArtworkKind = "backdrop"
	ArtworkBanner   ArtworkKind = "banner"
	ArtworkLogo     ArtworkKind = "logo"
)

// Artwork is a metadata image reference. The URL points at cover art, never at the
// playable media.
type Artwork struct {
	Kind   ArtworkKind `json:"kind" db:"kind"`
	URL    string      `json:"url" db:"url"`
	Width  *int        `json:"width" db:"width"`
	Height *int        `json:"height" db:"height"`
}

// Work is the canonical metadata identity — one per title, shared across all users.
//
// ProviderIDs maps an external catalog name (anilist, mal, tmdb, musicbrainz,
// hardcover, …) to that catalog's ID for this Work. It MAY be empty: a self-hosted
// user can have files that match no known provider (obscure or local-only works),
// and such a Work is still valid.
//
// The composite fields (Titles, Artwork, Genres, Tags, ProviderIDs) are stored as
// JSONB (or join tables) when the persistence layer is built; their db tags here are
// nominal and do not imply single-column mapping.
type Work struct {
	shared.BaseEntity
	Kind        Kind              `json:"kind" db:"kind"`
	Title       string            `json:"title" db:"title"` // primary display title
	Titles      []LocalizedTitle  `json:"titles" db:"titles"`
	Year        *int              `json:"year" db:"year"`
	Synopsis    *string           `json:"synopsis" db:"synopsis"`
	Artwork     []Artwork         `json:"artwork" db:"artwork"`
	Genres      []string          `json:"genres" db:"genres"`
	Tags        []string          `json:"tags" db:"tags"`
	ProviderIDs map[string]string `json:"provider_ids" db:"provider_ids"`
}

// ProviderID returns the external ID for the given provider and whether it is present.
func (w *Work) ProviderID(provider string) (string, bool) {
	if w.ProviderIDs == nil {
		return "", false
	}
	id, ok := w.ProviderIDs[provider]
	return id, ok
}

// WorkChild is a structural child of a Work: a season/episode/volume/chapter/track.
//
// Playable level: a childless Work (a movie, a standalone track) is itself the
// playable — a ContentReference points at the Work with a nil ChildID. Episodic and
// serial works are played at the child level.
//
// Ordinal is a float64 so that half-numbered entries (episode 7.5, chapter 10.5) are
// representable. Ordinals are unique per (WorkID, ParentID, Kind), not per Work —
// season 1 and season 2 can both contain an "episode 1".
type WorkChild struct {
	shared.BaseEntity
	WorkID      shared.UUID       `json:"work_id" db:"work_id"`
	ParentID    *shared.UUID      `json:"parent_id" db:"parent_id"` // e.g. an episode's season; nil for top-level children
	Kind        ChildKind         `json:"kind" db:"kind"`
	Ordinal     float64           `json:"ordinal" db:"ordinal"`
	Title       string            `json:"title" db:"title"`
	Synopsis    *string           `json:"synopsis" db:"synopsis"`
	Duration    *int              `json:"duration" db:"duration"` // seconds; nil for non-time-based children (chapters/volumes)
	ReleaseDate *time.Time        `json:"release_date" db:"release_date"`
	ProviderIDs map[string]string `json:"provider_ids" db:"provider_ids"`
}
