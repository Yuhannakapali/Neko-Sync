package entities

import "time"

// ========== WATCH PARTIES ==========

// WatchParty represents a synchronized viewing session
type WatchParty struct {
	BaseEntity
	HostUserID  UUID       `json:"host_user_id" db:"host_user_id"`
	ContentID   UUID       `json:"content_id" db:"content_id"`
	EpisodeID   *UUID      `json:"episode_id" db:"episode_id"`
	ChapterID   *UUID      `json:"chapter_id" db:"chapter_id"`
	RoomCode    string     `json:"room_code" db:"room_code"`
	Title       string     `json:"title" db:"title"`
	Description *string    `json:"description" db:"description"`
	MaxUsers    int        `json:"max_users" db:"max_users"`
	IsPrivate   bool       `json:"is_private" db:"is_private"`
	Password    *string    `json:"password" db:"password"`
	StartedAt   *time.Time `json:"started_at" db:"started_at"`
	EndedAt     *time.Time `json:"ended_at" db:"ended_at"`
	IsActive    bool       `json:"is_active" db:"is_active"`
}

// WatchPartyUser represents a user's participation in a watch party
type WatchPartyUser struct {
	PartyID  UUID       `json:"party_id" db:"party_id"`
	UserID   UUID       `json:"user_id" db:"user_id"`
	JoinedAt time.Time  `json:"joined_at" db:"joined_at"`
	LeftAt   *time.Time `json:"left_at" db:"left_at"`
	IsActive bool       `json:"is_active" db:"is_active"`
}

// WatchPartyState represents the current state of a watch party
type WatchPartyState struct {
	PartyID       UUID      `json:"party_id" db:"party_id"`
	CurrentTime   int       `json:"current_time" db:"current_time"` // in seconds
	IsPlaying     bool      `json:"is_playing" db:"is_playing"`
	PlaybackSpeed float64   `json:"playback_speed" db:"playback_speed"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy     UUID      `json:"updated_by" db:"updated_by"`
}

// ========== DEVICE TRANSFER ==========

// DeviceTransfer represents content transfer between devices
type DeviceTransfer struct {
	BaseEntity
	UserID       UUID        `json:"user_id" db:"user_id"`
	FromDeviceID UUID        `json:"from_device_id" db:"from_device_id"`
	ToDeviceID   UUID        `json:"to_device_id" db:"to_device_id"`
	ContentType  ContentType `json:"content_type" db:"content_type"`
	ContentID    UUID        `json:"content_id" db:"content_id"`
	EpisodeID    *UUID       `json:"episode_id" db:"episode_id"`
	ChapterID    *UUID       `json:"chapter_id" db:"chapter_id"`
	Position     int         `json:"position" db:"position"` // Current position in seconds/pages
	IsCompleted  bool        `json:"is_completed" db:"is_completed"`
}

// ========== PARTY CHAT ==========

// PartyMessage represents a chat message in a watch party
type PartyMessage struct {
	BaseEntity
	PartyID   UUID   `json:"party_id" db:"party_id"`
	UserID    UUID   `json:"user_id" db:"user_id"`
	Message   string `json:"message" db:"message"`
	Timestamp int    `json:"timestamp" db:"timestamp"` // Video timestamp when message was sent
}

// Business Rules for Party domain
func (wp *WatchParty) CanJoin(userID UUID) bool {
	// Logic to check if user can join the party
	return wp.IsActive && !wp.IsFull()
}

func (wp *WatchParty) IsFull() bool {
	// This would need to be implemented with a repository to count current users
	return false // Placeholder
}

func (wp *WatchParty) IsHost(userID UUID) bool {
	return wp.HostUserID == userID
}

func (wp *WatchParty) Start() {
	now := time.Now()
	wp.StartedAt = &now
	wp.IsActive = true
	wp.UpdatedAt = now
}

func (wp *WatchParty) End() {
	now := time.Now()
	wp.EndedAt = &now
	wp.IsActive = false
	wp.UpdatedAt = now
}

func (dt *DeviceTransfer) Complete() {
	dt.IsCompleted = true
	dt.UpdatedAt = time.Now()
}
