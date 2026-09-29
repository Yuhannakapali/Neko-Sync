package work

import "errors"

var (
	ErrNotFound      = errors.New("work not found")
	ErrChildNotFound = errors.New("work child not found")
)
