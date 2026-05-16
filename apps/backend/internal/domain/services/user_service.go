package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"nekosync/internal/domain/entities"
	domainerrors "nekosync/internal/domain/errors"
	"nekosync/internal/domain/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo         repositories.UserRepository
	deviceRepo       repositories.DeviceRepository
	followRepo       repositories.FollowRepository
	notificationRepo repositories.NotificationRepository
}

func NewUserService(
	userRepo repositories.UserRepository,
	deviceRepo repositories.DeviceRepository,
	followRepo repositories.FollowRepository,
	notificationRepo repositories.NotificationRepository,
) *UserService {
	return &UserService{
		userRepo:         userRepo,
		deviceRepo:       deviceRepo,
		followRepo:       followRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, username, email, password string) (*entities.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, fmt.Errorf("username, email, and password are required")
	}

	if existing, _ := s.userRepo.GetByEmail(ctx, email); existing != nil {
		return nil, domainerrors.ErrUserAlreadyExists
	}

	if existing, _ := s.userRepo.GetByUsername(ctx, username); existing != nil {
		return nil, domainerrors.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entities.User{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(generateID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         entities.UserRoleUser,
		IsVerified:   false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, domainerrors.ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, domainerrors.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID entities.UUID, profile *entities.UserProfile) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domainerrors.ErrUserNotFound
	}

	profile.UserID = user.ID

	existing, _ := s.userRepo.GetProfile(ctx, userID)
	if existing == nil {
		return s.userRepo.CreateProfile(ctx, profile)
	}

	return s.userRepo.UpdateProfile(ctx, profile)
}

func (s *UserService) FollowUser(ctx context.Context, followerID, followingID entities.UUID) error {
	if followerID == followingID {
		return domainerrors.ErrCannotFollowSelf
	}

	if _, err := s.userRepo.GetByID(ctx, followerID); err != nil {
		return domainerrors.ErrUserNotFound
	}

	if _, err := s.userRepo.GetByID(ctx, followingID); err != nil {
		return domainerrors.ErrUserNotFound
	}

	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to check follow status: %w", err)
	}

	if isFollowing {
		return domainerrors.ErrAlreadyFollowing
	}

	return s.followRepo.Follow(ctx, followerID, followingID)
}

func (s *UserService) UnfollowUser(ctx context.Context, followerID, followingID entities.UUID) error {
	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to check follow status: %w", err)
	}

	if !isFollowing {
		return domainerrors.ErrNotFollowing
	}

	return s.followRepo.Unfollow(ctx, followerID, followingID)
}

func (s *UserService) CreateNotification(ctx context.Context, userID entities.UUID, notificationType entities.NotificationType, title, message string, data map[string]interface{}) error {
	notification := &entities.Notification{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(generateID()),
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

	return s.notificationRepo.Create(ctx, notification)
}

func (s *UserService) RegisterDevice(ctx context.Context, userID entities.UUID, deviceName string, platform entities.PlatformType) (*entities.UserDevice, error) {
	if err := s.deviceRepo.DeactivateAllForUser(ctx, userID); err != nil {
		return nil, fmt.Errorf("failed to deactivate existing devices: %w", err)
	}

	device := &entities.UserDevice{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(generateID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:     userID,
		DeviceName: deviceName,
		Platform:   platform,
		LastSeen:   time.Now(),
		IsActive:   true,
	}

	if err := s.deviceRepo.Create(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	return device, nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
