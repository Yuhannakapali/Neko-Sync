// Package mediafile holds the Instance's playable unit: one file on disk, its probed
// technical metadata, and its mapping to a Hub Work through provider IDs.
//
// MediaFile never leaves the Instance. The Hub only ever sees a ContentReference
// whose Locator points back here (instanceID + MediaFile ID).
package mediafile

import (
	"path/filepath"
	"strings"
	"time"

	"nekosync-instance/internal/domain/shared"
)

// Kind is the library medium. The string values match the Hub's work.Kind so the
// Hub connector can pass them through unchanged.
type Kind string

const (
	KindMovie  Kind = "movie"
	KindSeries Kind = "series"
	KindAnime  Kind = "anime"
)

// Valid reports whether k is a supported library kind.
func (k Kind) Valid() bool {
	switch k {
	case KindMovie, KindSeries, KindAnime:
		return true
	}
	return false
}

// StreamType is the ffprobe codec_type of a stream that the Instance cares about.
type StreamType string

const (
	StreamVideo    StreamType = "video"
	StreamAudio    StreamType = "audio"
	StreamSubtitle StreamType = "subtitle"
)

// Stream is one track inside a container.
type Stream struct {
	Index    int        `json:"index"` // ffprobe stream index, used later by ffmpeg -map
	Type     StreamType `json:"type"`
	Codec    string     `json:"codec"`
	Language string     `json:"language,omitempty"` // ISO 639-2 from the container tags, e.g. "jpn"
	Title    string     `json:"title,omitempty"`
	Width    int        `json:"width,omitempty"`
	Height   int        `json:"height,omitempty"`
	Channels int        `json:"channels,omitempty"`
	Default  bool       `json:"default"`
	Forced   bool       `json:"forced"`
}

// Technical is what a probe returns for a file.
type Technical struct {
	Container string   `json:"container"`
	Duration  float64  `json:"duration"` // seconds
	Streams   []Stream `json:"streams"`
}

// MatchStatus records how the file was mapped to a Work.
type MatchStatus string

const (
	MatchUnmatched MatchStatus = "unmatched" // no provider ID yet
	MatchMatched   MatchStatus = "matched"   // matched automatically
	MatchManual    MatchStatus = "manual"    // fixed by the user; the scanner must never overwrite it
)

// MediaFile is one playable file in a library.
type MediaFile struct {
	shared.BaseEntity
	LibraryID shared.UUID `json:"library_id"`
	Path      string      `json:"path"` // absolute path on the Instance's disk
	Size      int64       `json:"size"`
	ModTime   time.Time   `json:"mod_time"`
	Kind      Kind        `json:"kind"`
	Technical

	// Parsed from the path. Input for matching, and shown for unmatched files.
	ParsedTitle string   `json:"parsed_title"`
	ParsedYear  *int     `json:"parsed_year,omitempty"`
	Season      *int     `json:"season,omitempty"`
	Episode     *float64 `json:"episode,omitempty"` // float64 to match WorkChild.Ordinal (episode 7.5)

	ProviderIDs map[string]string `json:"provider_ids"` // e.g. "tmdb" -> "11615"
	MatchStatus MatchStatus       `json:"match_status"`
}

// Unchanged reports whether the file on disk still has the size and modification
// time recorded at the last scan. The scanner skips re-probing unchanged files.
func (m *MediaFile) Unchanged(size int64, modTime time.Time) bool {
	return m.Size == size && m.ModTime.Equal(modTime)
}

// videoExts lists the container extensions the scanner picks up.
var videoExts = map[string]bool{
	".mkv": true, ".mp4": true, ".m4v": true, ".avi": true, ".mov": true,
	".webm": true, ".ts": true, ".m2ts": true, ".wmv": true, ".flv": true,
}

// IsVideoFile reports whether path has a supported video extension.
func IsVideoFile(path string) bool {
	return videoExts[strings.ToLower(filepath.Ext(path))]
}
