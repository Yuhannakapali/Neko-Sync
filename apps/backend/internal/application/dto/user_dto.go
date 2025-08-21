package dto

import "time"

// User DTOs
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UserResponse struct {
	ID         string  `json:"id"`
	Username   string  `json:"username"`
	Email      string  `json:"email"`
	AvatarURL  *string `json:"avatar_url"`
	Role       string  `json:"role"`
	IsVerified bool    `json:"is_verified"`
	CreatedAt  string  `json:"created_at"`
}

type UpdateProfileRequest struct {
	About     *string    `json:"about"`
	Location  *string    `json:"location"`
	Website   *string    `json:"website"`
	BannerURL *string    `json:"banner_url"`
	Birthdate *time.Time `json:"birthdate"`
}

type UserProfileResponse struct {
	UserID    string     `json:"user_id"`
	About     *string    `json:"about"`
	Location  *string    `json:"location"`
	Website   *string    `json:"website"`
	BannerURL *string    `json:"banner_url"`
	Birthdate *time.Time `json:"birthdate"`
}

type FollowUserRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type NotificationResponse struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	IsRead    bool                   `json:"is_read"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt string                 `json:"created_at"`
}

// Device DTOs
type RegisterDeviceRequest struct {
	DeviceName string `json:"device_name" validate:"required"`
	Platform   string `json:"platform" validate:"required"`
}

type DeviceResponse struct {
	ID          string `json:"id"`
	DeviceName  string `json:"device_name"`
	Platform    string `json:"platform"`
	LastSeen    string `json:"last_seen"`
	WebsocketID string `json:"websocket_id,omitempty"`
	IsActive    bool   `json:"is_active"`
}
