package user

import (
	"errors"

	"nekosync/internal/domain/entities"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// User operations
func (s *Service) RegisterUser(u *User) error {
	// business logic here - validate email, hash password, etc.
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Username == "" {
		return errors.New("username is required")
	}

	// Check if user already exists
	existingUser, _ := s.repo.GetByEmail(u.Email)
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	existingUser, _ = s.repo.GetByUsername(u.Username)
	if existingUser != nil {
		return errors.New("user with this username already exists")
	}

	return s.repo.Create(u)
}

func (s *Service) GetUser(id entities.UUID) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetUserByEmail(email string) (*User, error) {
	return s.repo.GetByEmail(email)
}

func (s *Service) UpdateUser(u *User) error {
	// business logic for updating user
	return s.repo.Update(u)
}

// User profile operations
func (s *Service) CreateUserProfile(profile *UserProfile) error {
	return s.repo.CreateProfile(profile)
}

func (s *Service) GetUserProfile(userID entities.UUID) (*UserProfile, error) {
	return s.repo.GetProfile(userID)
}

func (s *Service) UpdateUserProfile(profile *UserProfile) error {
	return s.repo.UpdateProfile(profile)
}

// User device operations
func (s *Service) RegisterDevice(device *UserDevice) error {
	return s.repo.CreateDevice(device)
}

func (s *Service) GetUserDevices(userID entities.UUID) ([]*UserDevice, error) {
	return s.repo.GetDevicesByUserID(userID)
}

// Follow operations
func (s *Service) FollowUser(followerID, followingID entities.UUID) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	isFollowing, err := s.repo.IsFollowing(followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return errors.New("already following this user")
	}

	follow := &UserFollow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return s.repo.CreateFollow(follow)
}

func (s *Service) UnfollowUser(followerID, followingID entities.UUID) error {
	return s.repo.DeleteFollow(followerID, followingID)
}

func (s *Service) GetFollowers(userID entities.UUID) ([]*User, error) {
	return s.repo.GetFollowers(userID)
}

func (s *Service) GetFollowing(userID entities.UUID) ([]*User, error) {
	return s.repo.GetFollowing(userID)
}

// Notification operations
func (s *Service) CreateNotification(notification *Notification) error {
	return s.repo.CreateNotification(notification)
}

func (s *Service) GetUserNotifications(userID entities.UUID) ([]*Notification, error) {
	return s.repo.GetNotificationsByUserID(userID)
}

func (s *Service) MarkNotificationAsRead(id entities.UUID) error {
	return s.repo.MarkNotificationAsRead(id)
}
