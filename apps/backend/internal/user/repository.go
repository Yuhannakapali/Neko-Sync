package user

import (
	"context"
	"nekosync/internal/platform/entity"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id entity.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id entity.UUID) error
	List(ctx context.Context, limit, offset int) ([]*User, error)

	CreateProfile(ctx context.Context, p *Profile) error
	GetProfile(ctx context.Context, userID entity.UUID) (*Profile, error)
	UpdateProfile(ctx context.Context, p *Profile) error
}

type DeviceRepository interface {
	Create(ctx context.Context, d *Device) error
	GetByUserID(ctx context.Context, userID entity.UUID) ([]*Device, error)
	Update(ctx context.Context, d *Device) error
	Delete(ctx context.Context, id entity.UUID) error
	DeactivateAllForUser(ctx context.Context, userID entity.UUID) error
	// ReplaceActive atomically deactivates the user's devices and inserts d.
	ReplaceActive(ctx context.Context, d *Device) error
}

type FollowRepository interface {
	Follow(ctx context.Context, followerID, followingID entity.UUID) error
	Unfollow(ctx context.Context, followerID, followingID entity.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID entity.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID entity.UUID, limit, offset int) ([]*User, error)
	GetFollowing(ctx context.Context, userID entity.UUID, limit, offset int) ([]*User, error)
	GetFollowCount(ctx context.Context, userID entity.UUID) (followers int, following int, err error)
}

type NotificationRepository interface {
	Create(ctx context.Context, n *Notification) error
	GetByUserID(ctx context.Context, userID entity.UUID, limit, offset int) ([]*Notification, error)
	GetUnread(ctx context.Context, userID entity.UUID) ([]*Notification, error)
	MarkAsRead(ctx context.Context, id entity.UUID) error
	MarkAllAsRead(ctx context.Context, userID entity.UUID) error
	Delete(ctx context.Context, id entity.UUID) error
}
