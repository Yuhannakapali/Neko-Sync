package party

import "errors"

var (
	ErrNotFound       = errors.New("party not found")
	ErrNotActive      = errors.New("party is not active")
	ErrFull           = errors.New("party is full")
	ErrWrongPassword  = errors.New("invalid party password")
	ErrAlreadyMember  = errors.New("user is already in the party")
	ErrNotMember      = errors.New("user is not in the party")
	ErrInvalidDevices = errors.New("invalid device IDs")
)
