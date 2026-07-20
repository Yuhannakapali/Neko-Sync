package party

import (
	"nekosync/internal/domain/shared"
	"time"
)

// ========== WATCH PARTY ==========

// WatchParty coordinates a shared *timeline* over a Work, not a shared stream: each
// member's client resolves its own source (via domain/reference) and applies the
// broadcast play/pause/seek locally.
type WatchParty struct {
	shared.BaseEntity
	HostUserID  shared.UUID  `json:"host_user_id" db:"host_user_id"`
	WorkID      shared.UUID  `json:"work_id" db:"work_id"`
	ChildID     *shared.UUID `json:"child_id" db:"child_id"` // specific episode/chapter/track, if applicable
	RoomCode    string       `json:"room_code" db:"room_code"`
	Title       string       `json:"title" db:"title"`
	Description *string      `json:"description" db:"description"`
	MaxUsers    int          `json:"max_users" db:"max_users"`
	IsPrivate   bool         `json:"is_private" db:"is_private"`
	Password    *string      `json:"password" db:"password"`
	StartedAt   *time.Time   `json:"started_at" db:"started_at"`
	EndedAt     *time.Time   `json:"ended_at" db:"ended_at"`
	IsActive    bool         `json:"is_active" db:"is_active"`
}

func (wp *WatchParty) IsFull(currentCount int) bool {
	return currentCount >= wp.MaxUsers
}

func (wp *WatchParty) CanJoin(currentCount int) bool {
	return wp.IsActive && !wp.IsFull(currentCount)
}

func (wp *WatchParty) IsHost(userID shared.UUID) bool {
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

type PartyMember struct {
	PartyID  shared.UUID `json:"party_id" db:"party_id"`
	UserID   shared.UUID `json:"user_id" db:"user_id"`
	JoinedAt time.Time   `json:"joined_at" db:"joined_at"`
	LeftAt   *time.Time  `json:"left_at" db:"left_at"`
	IsActive bool        `json:"is_active" db:"is_active"`
}

type PlaybackState struct {
	PartyID       shared.UUID `json:"party_id" db:"party_id"`
	CurrentTime   int         `json:"current_time" db:"current_time"`
	IsPlaying     bool        `json:"is_playing" db:"is_playing"`
	PlaybackSpeed float64     `json:"playback_speed" db:"playback_speed"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	UpdatedBy     shared.UUID `json:"updated_by" db:"updated_by"`
}

type Message struct {
	shared.BaseEntity
	PartyID   shared.UUID `json:"party_id" db:"party_id"`
	UserID    shared.UUID `json:"user_id" db:"user_id"`
	Message   string      `json:"message" db:"message"`
	Timestamp int         `json:"timestamp" db:"timestamp"`
}

// ========== DEVICE TRANSFER ==========

// DeviceTransfer hands off in-progress playback from one device to another.
// Position is kept as-is here; it folds into the unified Progress{Fraction, Locator}
// model in a later step.
type DeviceTransfer struct {
	shared.BaseEntity
	UserID       shared.UUID  `json:"user_id" db:"user_id"`
	FromDeviceID shared.UUID  `json:"from_device_id" db:"from_device_id"`
	ToDeviceID   shared.UUID  `json:"to_device_id" db:"to_device_id"`
	WorkID       shared.UUID  `json:"work_id" db:"work_id"`
	ChildID      *shared.UUID `json:"child_id" db:"child_id"`
	Position     int          `json:"position" db:"position"`
	IsCompleted  bool         `json:"is_completed" db:"is_completed"`
}

func (dt *DeviceTransfer) Complete() {
	dt.IsCompleted = true
	dt.UpdatedAt = time.Now()
}
