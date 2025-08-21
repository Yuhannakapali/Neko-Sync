package user

import (
	"nekosync/internal/domain"
	"time"
)

type User struct {
	domain.BaseEntity
	Username     string          `json:"username" db:"username"`
	Email        string          `json:"email" db:"email"`
	PasswordHash string          `json:"password_hash" db:"password_hash"`
	AvatarURL    *string         `json:"avatar_url" db:"avatar_url"`
	Role         domain.UserRole `json:"role" db:"role"`
	IsVerified   bool            `json:"is_verified" db:"is_verified"`
}

type UserProfile struct {
	UserID    domain.UUID `json:"user_id" db:"user_id"`
	About     *string     `json:"about" db:"about"`
	Location  *string     `json:"location" db:"location"`
	Website   *string     `json:"website" db:"website"`
	BannerURL *string     `json:"banner_url" db:"banner_url"`
	Birthdate *time.Time  `json:"birthdate" db:"birthdate"`
}

type UserDevice struct {
	domain.BaseEntity
	UserID      domain.UUID         `json:"user_id" db:"user_id"`
	DeviceName  string              `json:"device_name" db:"device_name"`
	Platform    domain.PlatformType `json:"platform" db:"platform"`
	LastSeen    time.Time           `json:"last_seen" db:"last_seen"`
	WebsocketID *string             `json:"websocket_id" db:"websocket_id"`
	IsActive    bool                `json:"is_active" db:"is_active"`
}

type UserFollow struct {
	FollowerID  domain.UUID `json:"follower_id" db:"follower_id"`
	FollowingID domain.UUID `json:"following_id" db:"following_id"`
	FollowedAt  time.Time   `json:"followed_at" db:"followed_at"`
}

type Notification struct {
	domain.BaseEntity
	UserID  domain.UUID `json:"user_id" db:"user_id"`
	Type    string      `json:"type" db:"type"`
	Message string      `json:"message" db:"message"`
	IsRead  bool        `json:"is_read" db:"is_read"`
}
