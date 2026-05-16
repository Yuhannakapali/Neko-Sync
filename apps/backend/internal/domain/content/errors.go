package content

import "errors"

var (
	ErrNotFound       = errors.New("content not found")
	ErrSeriesNotFound = errors.New("series not found")
)
