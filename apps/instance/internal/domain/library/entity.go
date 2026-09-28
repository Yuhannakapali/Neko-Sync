// Package library holds a library root: one folder of one kind (movies, series,
// or anime) that the scanner walks.
package library

import (
	"nekosync-instance/internal/domain/mediafile"
	"nekosync-instance/internal/domain/shared"
)

// Library is a scanned folder.
type Library struct {
	shared.BaseEntity
	Name string         `json:"name"`
	Kind mediafile.Kind `json:"kind"`
	Root string         `json:"root"` // absolute path
}
