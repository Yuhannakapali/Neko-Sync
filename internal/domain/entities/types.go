package entities

import "time"

// BaseEntity contains common fields for all entities
type BaseEntity struct {
	ID        UUID      `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UUID represents a unique identifier
type UUID string

// ========== ENUMS ==========

type ContentType string

const (
	ContentTypeAnime ContentType = "anime"
	ContentTypeManga ContentType = "manga"
	ContentTypeMovie ContentType = "movie"
	ContentTypeMusic ContentType = "music"
)

type Season string

const (
	SeasonWinter Season = "winter"
	SeasonSpring Season = "spring"
	SeasonSummer Season = "summer"
	SeasonFall   Season = "fall"
)

type ContentStatus string

const (
	ContentStatusOngoing   ContentStatus = "ongoing"
	ContentStatusCompleted ContentStatus = "completed"
	ContentStatusUpcoming  ContentStatus = "upcoming"
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type MusicType string

const (
	MusicTypeSong       MusicType = "song"
	MusicTypeAlbum      MusicType = "album"
	MusicTypeSoundtrack MusicType = "soundtrack"
)

type PlatformType string

const (
	PlatformTypeWeb     PlatformType = "web"
	PlatformTypeDesktop PlatformType = "desktop"
	PlatformTypeMobile  PlatformType = "mobile"
	PlatformTypeTablet  PlatformType = "tablet"
)

type NotificationType string

const (
	NotificationTypeFollow        NotificationType = "follow"
	NotificationTypeLike          NotificationType = "like"
	NotificationTypeComment       NotificationType = "comment"
	NotificationTypePartyInvite   NotificationType = "party_invite"
	NotificationTypeContentUpdate NotificationType = "content_update"
)
