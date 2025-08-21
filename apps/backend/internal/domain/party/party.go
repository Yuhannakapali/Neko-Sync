package party

import (
	"nekosync/internal/domain"
	"time"
)

// ========== WATCH PARTIES ==========

type WatchParty struct {
	domain.BaseEntity
	HostUserID domain.UUID  `json:"host_user_id" db:"host_user_id"`
	ContentID  domain.UUID  `json:"content_id" db:"content_id"`
	EpisodeID  *domain.UUID `json:"episode_id" db:"episode_id"`
	RoomCode   string       `json:"room_code" db:"room_code"`
	StartedAt  *time.Time   `json:"started_at" db:"started_at"`
	EndedAt    *time.Time   `json:"ended_at" db:"ended_at"`
	IsActive   bool         `json:"is_active" db:"is_active"`
}

type WatchPartyUser struct {
	PartyID  domain.UUID `json:"party_id" db:"party_id"`
	UserID   domain.UUID `json:"user_id" db:"user_id"`
	JoinedAt time.Time   `json:"joined_at" db:"joined_at"`
}

// ========== DEVICE TRANSFER ==========

type DeviceTransfer struct {
	domain.BaseEntity
	FromDeviceID domain.UUID        `json:"from_device_id" db:"from_device_id"`
	ToDeviceID   domain.UUID        `json:"to_device_id" db:"to_device_id"`
	ContentType  domain.ContentType `json:"content_type" db:"content_type"`
	ContentID    domain.UUID        `json:"content_id" db:"content_id"`
	Position     int                `json:"position" db:"position"`
}
