package entities

import "time"

// User represents a user in the system
type User struct {
	BaseEntity
	Username     string   `json:"username" db:"username"`
	Email        string   `json:"email" db:"email"`
	PasswordHash string   `json:"password_hash" db:"password_hash"`
	AvatarURL    *string  `json:"avatar_url" db:"avatar_url"`
	Role         UserRole `json:"role" db:"role"`
	IsVerified   bool     `json:"is_verified" db:"is_verified"`
}

// UserProfile represents extended user profile information
type UserProfile struct {
	UserID    UUID       `json:"user_id" db:"user_id"`
	About     *string    `json:"about" db:"about"`
	Location  *string    `json:"location" db:"location"`
	Website   *string    `json:"website" db:"website"`
	BannerURL *string    `json:"banner_url" db:"banner_url"`
	Birthdate *time.Time `json:"birthdate" db:"birthdate"`
}

// UserDevice represents a user's connected device
type UserDevice struct {
	BaseEntity
	UserID      UUID         `json:"user_id" db:"user_id"`
	DeviceName  string       `json:"device_name" db:"device_name"`
	Platform    PlatformType `json:"platform" db:"platform"`
	LastSeen    time.Time    `json:"last_seen" db:"last_seen"`
	WebsocketID *string      `json:"websocket_id" db:"websocket_id"`
	IsActive    bool         `json:"is_active" db:"is_active"`
}

// UserFollow represents a follow relationship between users
type UserFollow struct {
	FollowerID  UUID      `json:"follower_id" db:"follower_id"`
	FollowingID UUID      `json:"following_id" db:"following_id"`
	FollowedAt  time.Time `json:"followed_at" db:"followed_at"`
}

// Notification represents a user notification
type Notification struct {
	BaseEntity
	UserID  UUID                   `json:"user_id" db:"user_id"`
	Type    NotificationType       `json:"type" db:"type"`
	Title   string                 `json:"title" db:"title"`
	Message string                 `json:"message" db:"message"`
	IsRead  bool                   `json:"is_read" db:"is_read"`
	Data    map[string]interface{} `json:"data" db:"data"` // Additional notification data
}

// BusinessRules for User domain
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

func (u *User) CanModify(targetUser *User) bool {
	return u.ID == targetUser.ID || u.IsAdmin()
}

func (n *Notification) MarkAsRead() {
	n.IsRead = true
	n.UpdatedAt = time.Now()
}
