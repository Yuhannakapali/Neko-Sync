package user

import "errors"

var (
	ErrNotFound          = errors.New("user not found")
	ErrAlreadyExists     = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCannotFollowSelf  = errors.New("users cannot follow themselves")
	ErrAlreadyFollowing  = errors.New("already following this user")
	ErrNotFollowing      = errors.New("not following this user")
	ErrDeviceNotFound    = errors.New("device not found")
	ErrInvalidDeviceID   = errors.New("invalid device IDs")
)
