package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"nekosync/internal/domain/entities"
	"nekosync/internal/domain/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserService contains business logic for user operations
type UserService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with validation and password hashing
func (s *UserService) CreateUser(ctx context.Context, username, email, password string) (*entities.User, error) {
	// Validate input
	if username == "" || email == "" || password == "" {
		return nil, fmt.Errorf("username, email, and password are required")
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	existingUser, _ = s.userRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", username)
	}

	// Hash password
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &entities.User{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(s.generateUUID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         entities.UserRoleUser,
		IsVerified:   false,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// AuthenticateUser validates user credentials
func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !s.checkPassword(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// UpdateProfile updates a user's profile information
func (s *UserService) UpdateProfile(ctx context.Context, userID entities.UUID, profile *entities.UserProfile) error {
	// Validate that the user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Ensure the profile belongs to the user
	profile.UserID = user.ID

	// Check if profile exists
	existingProfile, _ := s.userRepo.GetProfile(ctx, userID)
	if existingProfile == nil {
		return s.userRepo.CreateProfile(ctx, profile)
	}

	return s.userRepo.UpdateProfile(ctx, profile)
}

// FollowUser creates a follow relationship between users
func (s *UserService) FollowUser(ctx context.Context, followerID, followingID entities.UUID) error {
	if followerID == followingID {
		return fmt.Errorf("users cannot follow themselves")
	}

	// Check if users exist
	_, err := s.userRepo.GetByID(ctx, followerID)
	if err != nil {
		return fmt.Errorf("follower not found: %w", err)
	}

	_, err = s.userRepo.GetByID(ctx, followingID)
	if err != nil {
		return fmt.Errorf("user to follow not found: %w", err)
	}

	// Check if already following
	isFollowing, err := s.userRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to check follow status: %w", err)
	}

	if isFollowing {
		return fmt.Errorf("already following this user")
	}

	return s.userRepo.Follow(ctx, followerID, followingID)
}

// UnfollowUser removes a follow relationship between users
func (s *UserService) UnfollowUser(ctx context.Context, followerID, followingID entities.UUID) error {
	isFollowing, err := s.userRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to check follow status: %w", err)
	}

	if !isFollowing {
		return fmt.Errorf("not following this user")
	}

	return s.userRepo.Unfollow(ctx, followerID, followingID)
}

// CreateNotification creates a new notification for a user
func (s *UserService) CreateNotification(ctx context.Context, userID entities.UUID, notificationType entities.NotificationType, title, message string, data map[string]interface{}) error {
	notification := &entities.Notification{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(s.generateUUID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Message: message,
		IsRead:  false,
		Data:    data,
	}

	return s.userRepo.CreateNotification(ctx, notification)
}

// RegisterDevice registers a new device for a user
func (s *UserService) RegisterDevice(ctx context.Context, userID entities.UUID, deviceName string, platform entities.PlatformType) (*entities.UserDevice, error) {
	// Deactivate all other devices for this user
	err := s.userRepo.DeactivateUserDevices(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to deactivate existing devices: %w", err)
	}

	device := &entities.UserDevice{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(s.generateUUID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:     userID,
		DeviceName: deviceName,
		Platform:   platform,
		LastSeen:   time.Now(),
		IsActive:   true,
	}

	err = s.userRepo.CreateDevice(ctx, device)
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	return device, nil
}

// Private helper methods
func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *UserService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserService) generateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
