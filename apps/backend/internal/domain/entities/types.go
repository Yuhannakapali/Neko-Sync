package entities

import "time"

// BaseEntity contains common fields for all entities.
type BaseEntity struct {
	ID        UUID      `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UUID represents a unique identifier.
type UUID string

// ========== ENUMS ==========

// ContentType represents the type of content in the system.
type ContentType string

// Content type constants.
const (
	ContentTypeAnime ContentType = "anime"
	ContentTypeManga ContentType = "manga"
	ContentTypeMovie ContentType = "movie"
	ContentTypeMusic ContentType = "music"
)

// Season represents the season when content was released.
type Season string

// Season constants.
const (
	SeasonWinter Season = "winter"
	SeasonSpring Season = "spring"
	SeasonSummer Season = "summer"
	SeasonFall   Season = "fall"
)

// ContentStatus represents the status of content.
type ContentStatus string

// Content status constants.
const (
	ContentStatusOngoing   ContentStatus = "ongoing"
	ContentStatusCompleted ContentStatus = "completed"
	ContentStatusUpcoming  ContentStatus = "upcoming"
)

// UserRole represents the role of a user in the system.
type UserRole string

// User role constants.
const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// MusicType represents the type of music content.
type MusicType string

// Music type constants.
const (
	MusicTypeSong       MusicType = "song"
	MusicTypeAlbum      MusicType = "album"
	MusicTypeSoundtrack MusicType = "soundtrack"
)

// PlatformType represents the platform type for device access.
type PlatformType string

// Platform type constants.
const (
	PlatformTypeWeb     PlatformType = "web"
	PlatformTypeDesktop PlatformType = "desktop"
	PlatformTypeMobile  PlatformType = "mobile"
	PlatformTypeTablet  PlatformType = "tablet"
)

// NotificationType represents the type of notification.
type NotificationType string

// Notification type constants.
const (
	NotificationTypeFollow        NotificationType = "follow"
	NotificationTypeLike          NotificationType = "like"
	NotificationTypeComment       NotificationType = "comment"
	NotificationTypePartyInvite   NotificationType = "party_invite"
	NotificationTypeContentUpdate NotificationType = "content_update"
)
