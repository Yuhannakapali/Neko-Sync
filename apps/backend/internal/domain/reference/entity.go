// Package reference holds the ContentReference registry — the per-user pointer to
// where that user's copy of a Work actually lives. The Hub stores the pointer, never
// the payload; resolving a reference to real bytes happens edge-to-edge (client ↔
// Instance / local file / official platform), never through the Hub.
package reference

import (
	"nekosync/internal/domain/shared"
)

// SourceKind is where a referenced copy lives, in descending playback preference.
type SourceKind string

const (
	SourceInstance     SourceKind = "instance"      // on the user's self-hosted instance
	SourceLocalFile    SourceKind = "local_file"    // on a specific client device
	SourcePublicDomain SourceKind = "public_domain" // Internet Archive / Gutenberg (legal to centralize)
	SourceOfficialLink SourceKind = "official_link" // deep-link to a licensed platform
)

// rankOrder gives each source a preference weight (lower = preferred).
var rankOrder = map[SourceKind]int{
	SourceInstance:     0,
	SourceLocalFile:    1,
	SourcePublicDomain: 2,
	SourceOfficialLink: 3,
}

// ContentReference maps a user + Work (and optionally a specific child) to where that
// user can actually play it.
//
// A nil ChildID means the Work itself is the playable (a movie, a standalone track).
// A non-nil ChildID points at a specific episode/chapter/track.
//
// Locator is source-specific and opaque to the Hub: "instanceID+nodeID" for an
// instance, "deviceID+path" for a local file, or an external URL for a link.
type ContentReference struct {
	shared.BaseEntity
	UserID  shared.UUID  `json:"user_id" db:"user_id"`
	WorkID  shared.UUID  `json:"work_id" db:"work_id"`
	ChildID *shared.UUID `json:"child_id" db:"child_id"`
	Source  SourceKind   `json:"source" db:"source"`
	Locator string       `json:"locator" db:"locator"`
	Quality *string      `json:"quality" db:"quality"` // free-form for now; grows into language/codec/resolution (Instance concern)
}

// Rank returns refs ordered by playback preference: instance > local_file >
// public_domain > official_link. Order within the same source is preserved (stable),
// so a caller that pre-sorts by quality keeps that sub-order. Ranking is advisory —
// the client makes the final choice and may override it.
func Rank(refs []*ContentReference) []*ContentReference {
	out := make([]*ContentReference, len(refs))
	copy(out, refs)
	// insertion sort keeps it stable and dependency-free for small reference sets.
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && rankOrder[out[j].Source] < rankOrder[out[j-1].Source]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}
