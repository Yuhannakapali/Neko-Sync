package user

import "nekosync/internal/domain/entities"

type Repository interface {
	// User operations
	Create(user *User) error
	GetByID(id entities.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id entities.UUID) error

	// User profile operations
	CreateProfile(profile *UserProfile) error
	GetProfile(userID entities.UUID) (*UserProfile, error)
	UpdateProfile(profile *UserProfile) error

	// User device operations
	CreateDevice(device *UserDevice) error
	GetDevicesByUserID(userID entities.UUID) ([]*UserDevice, error)
	UpdateDevice(device *UserDevice) error
	DeleteDevice(id entities.UUID) error

	// User follow operations
	CreateFollow(follow *UserFollow) error
	DeleteFollow(followerID, followingID entities.UUID) error
	GetFollowers(userID entities.UUID) ([]*User, error)
	GetFollowing(userID entities.UUID) ([]*User, error)
	IsFollowing(followerID, followingID entities.UUID) (bool, error)

	// Notification operations
	CreateNotification(notification *Notification) error
	GetNotificationsByUserID(userID entities.UUID) ([]*Notification, error)
	MarkNotificationAsRead(id entities.UUID) error
	DeleteNotification(id entities.UUID) error
}
