package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

// UserRepository defines the contract for user data operations
type UserRepository interface {
	// User CRUD operations
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id entities.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id entities.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// User Profile operations
	CreateProfile(ctx context.Context, profile *entities.UserProfile) error
	GetProfile(ctx context.Context, userID entities.UUID) (*entities.UserProfile, error)
	UpdateProfile(ctx context.Context, profile *entities.UserProfile) error

	// User Device operations
	CreateDevice(ctx context.Context, device *entities.UserDevice) error
	GetDevicesByUserID(ctx context.Context, userID entities.UUID) ([]*entities.UserDevice, error)
	UpdateDevice(ctx context.Context, device *entities.UserDevice) error
	DeleteDevice(ctx context.Context, id entities.UUID) error
	DeactivateUserDevices(ctx context.Context, userID entities.UUID) error

	// User Follow operations
	Follow(ctx context.Context, followerID, followingID entities.UUID) error
	Unfollow(ctx context.Context, followerID, followingID entities.UUID) error
	GetFollowers(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.User, error)
	GetFollowing(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.User, error)
	IsFollowing(ctx context.Context, followerID, followingID entities.UUID) (bool, error)
	GetFollowCount(ctx context.Context, userID entities.UUID) (followers int, following int, err error)

	// Notification operations
	CreateNotification(ctx context.Context, notification *entities.Notification) error
	GetNotifications(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.Notification, error)
	GetUnreadNotifications(ctx context.Context, userID entities.UUID) ([]*entities.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID entities.UUID) error
	MarkAllNotificationsAsRead(ctx context.Context, userID entities.UUID) error
	DeleteNotification(ctx context.Context, id entities.UUID) error
}
