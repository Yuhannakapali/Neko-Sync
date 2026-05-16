package services

import (
	"context"
	"nekosync/internal/domain/entities"
)

// UserServiceInterface defines the contract for user operations used by the application layer.
type UserServiceInterface interface {
	CreateUser(ctx context.Context, username, email, password string) (*entities.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (*entities.User, error)
	UpdateProfile(ctx context.Context, userID entities.UUID, profile *entities.UserProfile) error
	FollowUser(ctx context.Context, followerID, followingID entities.UUID) error
	UnfollowUser(ctx context.Context, followerID, followingID entities.UUID) error
	CreateNotification(ctx context.Context, userID entities.UUID, notificationType entities.NotificationType, title, message string, data map[string]interface{}) error
	RegisterDevice(ctx context.Context, userID entities.UUID, deviceName string, platform entities.PlatformType) (*entities.UserDevice, error)
}

// PartyServiceInterface defines the contract for watch party operations used by the application layer.
type PartyServiceInterface interface {
	CreateWatchParty(ctx context.Context, hostUserID, contentID entities.UUID, title string, maxUsers int, isPrivate bool, password *string) (*entities.WatchParty, error)
	JoinWatchParty(ctx context.Context, userID entities.UUID, roomCode string, password *string) (*entities.WatchParty, error)
	LeaveWatchParty(ctx context.Context, userID, partyID entities.UUID) error
	UpdatePartyState(ctx context.Context, partyID, userID entities.UUID, currentTime int, isPlaying bool, playbackSpeed float64) error
	SendMessage(ctx context.Context, partyID, userID entities.UUID, message string, timestamp int) error
	CreateDeviceTransfer(ctx context.Context, userID, fromDeviceID, toDeviceID, contentID entities.UUID, position int) (*entities.DeviceTransfer, error)
}
