package domainerrors

import "errors"

// User domain errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCannotFollowSelf  = errors.New("users cannot follow themselves")
	ErrAlreadyFollowing  = errors.New("already following this user")
	ErrNotFollowing      = errors.New("not following this user")
)

// Content domain errors
var (
	ErrContentNotFound = errors.New("content not found")
	ErrSeriesNotFound  = errors.New("series not found")
)

// Party domain errors
var (
	ErrPartyNotFound   = errors.New("party not found")
	ErrPartyNotActive  = errors.New("party is not active")
	ErrPartyFull       = errors.New("party is full")
	ErrPartyWrongPassword = errors.New("invalid party password")
	ErrUserAlreadyInParty = errors.New("user is already in the party")
	ErrUserNotInParty  = errors.New("user is not in the party")
)

// Device domain errors
var (
	ErrDeviceNotFound  = errors.New("device not found")
	ErrInvalidDeviceID = errors.New("invalid device IDs")
)
