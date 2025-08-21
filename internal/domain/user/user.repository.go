package user

import "nekosync/internal/domain"

type Repository interface {
	// User operations
	Create(user *User) error
	GetByID(id domain.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id domain.UUID) error

	// User profile operations
	CreateProfile(profile *UserProfile) error
	GetProfile(userID domain.UUID) (*UserProfile, error)
	UpdateProfile(profile *UserProfile) error

	// User device operations
	CreateDevice(device *UserDevice) error
	GetDevicesByUserID(userID domain.UUID) ([]*UserDevice, error)
	UpdateDevice(device *UserDevice) error
	DeleteDevice(id domain.UUID) error

	// User follow operations
	CreateFollow(follow *UserFollow) error
	DeleteFollow(followerID, followingID domain.UUID) error
	GetFollowers(userID domain.UUID) ([]*User, error)
	GetFollowing(userID domain.UUID) ([]*User, error)
	IsFollowing(followerID, followingID domain.UUID) (bool, error)

	// Notification operations
	CreateNotification(notification *Notification) error
	GetNotificationsByUserID(userID domain.UUID) ([]*Notification, error)
	MarkNotificationAsRead(id domain.UUID) error
	DeleteNotification(id domain.UUID) error
}
