package dto

import "time"

// CreateUserRequest represents the request to create a new user.
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// CreateUserResponse represents the response after creating a user.
type CreateUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// LoginRequest represents the request for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response after successful login.
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// UserResponse represents a user in API responses.
type UserResponse struct {
	ID         string  `json:"id"`
	Username   string  `json:"username"`
	Email      string  `json:"email"`
	AvatarURL  *string `json:"avatar_url"`
	Role       string  `json:"role"`
	IsVerified bool    `json:"is_verified"`
	CreatedAt  string  `json:"created_at"`
}

// UpdateProfileRequest represents the request to update user profile.
type UpdateProfileRequest struct {
	About     *string    `json:"about"`
	Location  *string    `json:"location"`
	Website   *string    `json:"website"`
	BannerURL *string    `json:"banner_url"`
	Birthdate *time.Time `json:"birthdate"`
}

// UserProfileResponse represents a user profile in API responses.
type UserProfileResponse struct {
	UserID    string     `json:"user_id"`
	About     *string    `json:"about"`
	Location  *string    `json:"location"`
	Website   *string    `json:"website"`
	BannerURL *string    `json:"banner_url"`
	Birthdate *time.Time `json:"birthdate"`
}

// FollowUserRequest represents the request to follow a user.
type FollowUserRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// NotificationResponse represents a notification in API responses.
type NotificationResponse struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	IsRead    bool                   `json:"is_read"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt string                 `json:"created_at"`
}

// RegisterDeviceRequest represents the request to register a device.
type RegisterDeviceRequest struct {
	DeviceName string `json:"device_name" validate:"required"`
	Platform   string `json:"platform" validate:"required"`
}

// DeviceResponse represents a device in API responses.
type DeviceResponse struct {
	ID          string `json:"id"`
	DeviceName  string `json:"device_name"`
	Platform    string `json:"platform"`
	LastSeen    string `json:"last_seen"`
	WebsocketID string `json:"websocket_id,omitempty"`
	IsActive    bool   `json:"is_active"`
}
