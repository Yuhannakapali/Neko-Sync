package userdevice

// UserDevice represents a device associated with a user.
type UserDevice struct {
	ID          string `json:"id" db:"id"`
	UserID      string `json:"user_id" db:"user_id"`
	DeviceName  string `json:"device_name" db:"device_name"`
	Platform    string `json:"platform" db:"platform"`
	LastSeen    string `json:"last_seen" db:"last_seen"`
	WebsocketID string `json:"websocket_id" db:"websocket_id"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}
