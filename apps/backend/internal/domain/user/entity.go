package user

import (
	"nekosync/internal/domain/shared"
	"time"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type PlatformType string

const (
	PlatformWeb     PlatformType = "web"
	PlatformDesktop PlatformType = "desktop"
	PlatformMobile  PlatformType = "mobile"
	PlatformTablet  PlatformType = "tablet"
)

type NotificationType string

const (
	NotificationFollow        NotificationType = "follow"
	NotificationLike          NotificationType = "like"
	NotificationComment       NotificationType = "comment"
	NotificationPartyInvite   NotificationType = "party_invite"
	NotificationContentUpdate NotificationType = "content_update"
)

type User struct {
	shared.BaseEntity
	Username     string  `json:"username" db:"username"`
	Email        string  `json:"email" db:"email"`
	PasswordHash string  `json:"password_hash" db:"password_hash"`
	AvatarURL    *string `json:"avatar_url" db:"avatar_url"`
	Role         Role    `json:"role" db:"role"`
	IsVerified   bool    `json:"is_verified" db:"is_verified"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) CanModify(target *User) bool {
	return u.ID == target.ID || u.IsAdmin()
}

type Profile struct {
	UserID    shared.UUID `json:"user_id" db:"user_id"`
	About     *string     `json:"about" db:"about"`
	Location  *string     `json:"location" db:"location"`
	Website   *string     `json:"website" db:"website"`
	BannerURL *string     `json:"banner_url" db:"banner_url"`
	Birthdate *time.Time  `json:"birthdate" db:"birthdate"`
}

type Device struct {
	shared.BaseEntity
	UserID      shared.UUID  `json:"user_id" db:"user_id"`
	DeviceName  string       `json:"device_name" db:"device_name"`
	Platform    PlatformType `json:"platform" db:"platform"`
	LastSeen    time.Time    `json:"last_seen" db:"last_seen"`
	WebsocketID *string      `json:"websocket_id" db:"websocket_id"`
	IsActive    bool         `json:"is_active" db:"is_active"`
}

type Follow struct {
	FollowerID  shared.UUID `json:"follower_id" db:"follower_id"`
	FollowingID shared.UUID `json:"following_id" db:"following_id"`
	FollowedAt  time.Time   `json:"followed_at" db:"followed_at"`
}

type Notification struct {
	shared.BaseEntity
	UserID  shared.UUID            `json:"user_id" db:"user_id"`
	Type    NotificationType       `json:"type" db:"type"`
	Title   string                 `json:"title" db:"title"`
	Message string                 `json:"message" db:"message"`
	IsRead  bool                   `json:"is_read" db:"is_read"`
	Data    map[string]interface{} `json:"data" db:"data"`
}

func (n *Notification) MarkAsRead() {
	n.IsRead = true
	n.UpdatedAt = time.Now()
}
