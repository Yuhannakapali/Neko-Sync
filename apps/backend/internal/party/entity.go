package party

import (
	"nekosync/internal/platform/entity"
	"nekosync/internal/progress"
	"time"
)

// ========== WATCH PARTY ==========

// WatchParty coordinates a shared *timeline* over a Work, not a shared stream: each
// member's client resolves its own source (via domain/reference) and applies the
// broadcast play/pause/seek locally.
type WatchParty struct {
	entity.BaseEntity
	HostUserID  entity.UUID  `json:"host_user_id" db:"host_user_id"`
	WorkID      entity.UUID  `json:"work_id" db:"work_id"`
	ChildID     *entity.UUID `json:"child_id" db:"child_id"` // specific episode/chapter/track, if applicable
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

func (wp *WatchParty) IsHost(userID entity.UUID) bool {
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
	PartyID  entity.UUID `json:"party_id" db:"party_id"`
	UserID   entity.UUID `json:"user_id" db:"user_id"`
	JoinedAt time.Time   `json:"joined_at" db:"joined_at"`
	LeftAt   *time.Time  `json:"left_at" db:"left_at"`
	IsActive bool        `json:"is_active" db:"is_active"`
}

type PlaybackState struct {
	PartyID       entity.UUID       `json:"party_id" db:"party_id"`
	Position      progress.Progress `json:"position"`
	IsPlaying     bool              `json:"is_playing" db:"is_playing"`
	PlaybackSpeed float64           `json:"playback_speed" db:"playback_speed"`
	UpdatedAt     time.Time         `json:"updated_at" db:"updated_at"`
	UpdatedBy     entity.UUID       `json:"updated_by" db:"updated_by"`
}

type Message struct {
	entity.BaseEntity
	PartyID   entity.UUID `json:"party_id" db:"party_id"`
	UserID    entity.UUID `json:"user_id" db:"user_id"`
	Message   string      `json:"message" db:"message"`
	Timestamp int         `json:"timestamp" db:"timestamp"`
}

// ========== DEVICE TRANSFER ==========

// DeviceTransfer hands off in-progress playback from one device to another.
type DeviceTransfer struct {
	entity.BaseEntity
	UserID       entity.UUID       `json:"user_id" db:"user_id"`
	FromDeviceID entity.UUID       `json:"from_device_id" db:"from_device_id"`
	ToDeviceID   entity.UUID       `json:"to_device_id" db:"to_device_id"`
	WorkID       entity.UUID       `json:"work_id" db:"work_id"`
	ChildID      *entity.UUID      `json:"child_id" db:"child_id"`
	Position     progress.Progress `json:"position"`
	IsCompleted  bool              `json:"is_completed" db:"is_completed"`
}

func (dt *DeviceTransfer) Complete() {
	dt.IsCompleted = true
	dt.UpdatedAt = time.Now()
}
