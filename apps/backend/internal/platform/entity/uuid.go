package entity

import (
	"crypto/rand"
	"fmt"
)

// NewUUID returns a random RFC 4122 version 4 UUID in canonical lowercase form,
// the same form Postgres returns from a uuid column.
func NewUUID() UUID {
	var b [16]byte
	_, _ = rand.Read(b[:]) // crypto/rand.Read never returns an error since Go 1.24
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return UUID(fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]))
}
