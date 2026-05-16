package user

import (
	"context"
	"nekosync/internal/domain/shared"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id shared.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id shared.UUID) error
	List(ctx context.Context, limit, offset int) ([]*User, error)

	CreateProfile(ctx context.Context, p *Profile) error
	GetProfile(ctx context.Context, userID shared.UUID) (*Profile, error)
	UpdateProfile(ctx context.Context, p *Profile) error
}

type DeviceRepository interface {
	Create(ctx context.Context, d *Device) error
	GetByUserID(ctx context.Context, userID shared.UUID) ([]*Device, error)
	Update(ctx context.Context, d *Device) error
	Delete(ctx context.Context, id shared.UUID) error
	DeactivateAllForUser(ctx context.Context, userID shared.UUID) error
}

type FollowRepository interface {
	Follow(ctx context.Context, followerID, followingID shared.UUID) error
	Unfollow(ctx context.Context, followerID, followingID shared.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID shared.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID shared.UUID, limit, offset int) ([]*User, error)
	GetFollowing(ctx context.Context, userID shared.UUID, limit, offset int) ([]*User, error)
	GetFollowCount(ctx context.Context, userID shared.UUID) (followers int, following int, err error)
}

type NotificationRepository interface {
	Create(ctx context.Context, n *Notification) error
	GetByUserID(ctx context.Context, userID shared.UUID, limit, offset int) ([]*Notification, error)
	GetUnread(ctx context.Context, userID shared.UUID) ([]*Notification, error)
	MarkAsRead(ctx context.Context, id shared.UUID) error
	MarkAllAsRead(ctx context.Context, userID shared.UUID) error
	Delete(ctx context.Context, id shared.UUID) error
}
