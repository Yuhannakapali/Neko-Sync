package user

import (
	"errors"
	"strings"
	"time"
	"nekosync/internal/domain/entities"
)

type User struct {
	entities.BaseEntity
	Username     string               `json:"username" db:"username"`
	Email        string               `json:"email" db:"email"`
	PasswordHash string               `json:"password_hash" db:"password_hash"`
	AvatarURL    *string              `json:"avatar_url" db:"avatar_url"`
	Role         entities.UserRole    `json:"role" db:"role"`
	IsVerified   bool                 `json:"is_verified" db:"is_verified"`
}

type UserProfile struct {
	UserID    entities.UUID `json:"user_id" db:"user_id"`
	About     *string       `json:"about" db:"about"`
	Location  *string       `json:"location" db:"location"`
	Website   *string       `json:"website" db:"website"`
	BannerURL *string       `json:"banner_url" db:"banner_url"`
	Birthdate *time.Time    `json:"birthdate" db:"birthdate"`
}

type UserDevice struct {
	entities.BaseEntity
	UserID      entities.UUID         `json:"user_id" db:"user_id"`
	DeviceName  string                `json:"device_name" db:"device_name"`
	Platform    entities.PlatformType `json:"platform" db:"platform"`
	LastSeen    time.Time             `json:"last_seen" db:"last_seen"`
	WebsocketID *string               `json:"websocket_id" db:"websocket_id"`
	IsActive    bool                  `json:"is_active" db:"is_active"`
}

type UserFollow struct {
	FollowerID  entities.UUID `json:"follower_id" db:"follower_id"`
	FollowingID entities.UUID `json:"following_id" db:"following_id"`
	FollowedAt  time.Time     `json:"followed_at" db:"followed_at"`
}

type Notification struct {
	entities.BaseEntity
	UserID  entities.UUID `json:"user_id" db:"user_id"`
	Type    string        `json:"type" db:"type"`
	Message string        `json:"message" db:"message"`
	IsRead  bool          `json:"is_read" db:"is_read"`
}

// User validation methods
func (u *User) ValidateEmail() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	// Basic email validation
	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func (u *User) ValidateUsername() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(u.Username) > 20 {
		return errors.New("username must be less than 20 characters")
	}
	return nil
}

// UserProfile validation methods
func (up *UserProfile) Validate() error {
	if up.UserID == "" {
		return errors.New("user ID is required")
	}
	return nil
}
