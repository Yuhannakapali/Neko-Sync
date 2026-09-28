// Package shared holds cross-aggregate types for the Instance.
//
// The Instance is a separate Go module from the Hub (apps/backend), so it cannot
// import nekosync/internal/... . These types mirror the Hub's shared.UUID and
// shared.BaseEntity so that IDs cross the Hub connector as plain strings.
package shared

import (
	"crypto/rand"
	"fmt"
	"time"
)

// UUID is a string-encoded UUID, the same representation the Hub uses.
type UUID string

// NewUUID returns a random (version 4) UUID.
func NewUUID() UUID {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Sprintf("shared: crypto/rand failed: %v", err))
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant
	return UUID(fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]))
}

// BaseEntity carries identity and timestamps for every persisted entity.
type BaseEntity struct {
	ID        UUID      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
