package domain

import "time"

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
	PlatformTypeAndroid PlatformType = "android"
	PlatformTypeIOS     PlatformType = "ios"
	PlatformTypeWeb     PlatformType = "web"
	PlatformTypeDesktop PlatformType = "desktop"
)

type ReportStatus string

const (
	ReportStatusPending  ReportStatus = "pending"
	ReportStatusReviewed ReportStatus = "reviewed"
	ReportStatusResolved ReportStatus = "resolved"
)

// ========== BASE TYPES ==========

type UUID string

type BaseEntity struct {
	ID        UUID      `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
