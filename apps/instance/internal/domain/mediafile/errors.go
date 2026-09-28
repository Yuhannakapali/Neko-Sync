package mediafile

import "errors"

// ErrNotFound is returned by a Repository when no MediaFile matches.
var ErrNotFound = errors.New("media file not found")
