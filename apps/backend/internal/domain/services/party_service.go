package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"nekosync/internal/domain/entities"
	domainerrors "nekosync/internal/domain/errors"
	"nekosync/internal/domain/repositories"
	"time"
)

type PartyService struct {
	partyRepo    repositories.WatchPartyRepository
	transferRepo repositories.DeviceTransferRepository
	contentRepo  repositories.ContentRepository
	userRepo     repositories.UserRepository
	deviceRepo   repositories.DeviceRepository
}

func NewPartyService(
	partyRepo repositories.WatchPartyRepository,
	transferRepo repositories.DeviceTransferRepository,
	contentRepo repositories.ContentRepository,
	userRepo repositories.UserRepository,
	deviceRepo repositories.DeviceRepository,
) *PartyService {
	return &PartyService{
		partyRepo:    partyRepo,
		transferRepo: transferRepo,
		contentRepo:  contentRepo,
		userRepo:     userRepo,
		deviceRepo:   deviceRepo,
	}
}

func (s *PartyService) CreateWatchParty(ctx context.Context, hostUserID, contentID entities.UUID, title string, maxUsers int, isPrivate bool, password *string) (*entities.WatchParty, error) {
	if _, err := s.userRepo.GetByID(ctx, hostUserID); err != nil {
		return nil, fmt.Errorf("host user not found: %w", err)
	}

	if _, err := s.contentRepo.GetByID(ctx, contentID); err != nil {
		return nil, fmt.Errorf("content not found: %w", err)
	}

	party := &entities.WatchParty{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(generateID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		HostUserID: hostUserID,
		ContentID:  contentID,
		RoomCode:   generateRoomCode(),
		Title:      title,
		MaxUsers:   maxUsers,
		IsPrivate:  isPrivate,
		Password:   password,
		IsActive:   true,
	}

	if err := s.partyRepo.Create(ctx, party); err != nil {
		return nil, fmt.Errorf("failed to create watch party: %w", err)
	}

	partyUser := &entities.WatchPartyUser{
		PartyID:  party.ID,
		UserID:   hostUserID,
		JoinedAt: time.Now(),
		IsActive: true,
	}

	if err := s.partyRepo.AddUser(ctx, partyUser); err != nil {
		return nil, fmt.Errorf("failed to add host to party: %w", err)
	}

	return party, nil
}

func (s *PartyService) JoinWatchParty(ctx context.Context, userID entities.UUID, roomCode string, password *string) (*entities.WatchParty, error) {
	party, err := s.partyRepo.GetByRoomCode(ctx, roomCode)
	if err != nil {
		return nil, domainerrors.ErrPartyNotFound
	}

	if !party.IsActive {
		return nil, domainerrors.ErrPartyNotActive
	}

	if party.IsPrivate {
		if password == nil || party.Password == nil || *password != *party.Password {
			return nil, domainerrors.ErrPartyWrongPassword
		}
	}

	isInParty, err := s.partyRepo.IsUserInParty(ctx, party.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check party membership: %w", err)
	}

	if isInParty {
		return nil, domainerrors.ErrUserAlreadyInParty
	}

	count, err := s.partyRepo.GetPartyUserCount(ctx, party.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get party user count: %w", err)
	}

	if party.IsFull(count) {
		return nil, domainerrors.ErrPartyFull
	}

	partyUser := &entities.WatchPartyUser{
		PartyID:  party.ID,
		UserID:   userID,
		JoinedAt: time.Now(),
		IsActive: true,
	}

	if err := s.partyRepo.AddUser(ctx, partyUser); err != nil {
		return nil, fmt.Errorf("failed to join party: %w", err)
	}

	return party, nil
}

func (s *PartyService) LeaveWatchParty(ctx context.Context, userID, partyID entities.UUID) error {
	party, err := s.partyRepo.GetByID(ctx, partyID)
	if err != nil {
		return domainerrors.ErrPartyNotFound
	}

	isInParty, err := s.partyRepo.IsUserInParty(ctx, partyID, userID)
	if err != nil {
		return fmt.Errorf("failed to check party membership: %w", err)
	}

	if !isInParty {
		return domainerrors.ErrUserNotInParty
	}

	if err := s.partyRepo.RemoveUser(ctx, partyID, userID); err != nil {
		return fmt.Errorf("failed to leave party: %w", err)
	}

	if party.HostUserID == userID {
		party.End()
		if err := s.partyRepo.Update(ctx, party); err != nil {
			return fmt.Errorf("failed to end party: %w", err)
		}
	}

	return nil
}

func (s *PartyService) UpdatePartyState(ctx context.Context, partyID, userID entities.UUID, currentTime int, isPlaying bool, playbackSpeed float64) error {
	party, err := s.partyRepo.GetByID(ctx, partyID)
	if err != nil {
		return domainerrors.ErrPartyNotFound
	}

	if !party.IsActive {
		return domainerrors.ErrPartyNotActive
	}

	isInParty, err := s.partyRepo.IsUserInParty(ctx, partyID, userID)
	if err != nil {
		return fmt.Errorf("failed to check party membership: %w", err)
	}

	if !isInParty {
		return domainerrors.ErrUserNotInParty
	}

	state := &entities.WatchPartyState{
		PartyID:       partyID,
		CurrentTime:   currentTime,
		IsPlaying:     isPlaying,
		PlaybackSpeed: playbackSpeed,
		UpdatedAt:     time.Now(),
		UpdatedBy:     userID,
	}

	return s.partyRepo.UpdateState(ctx, state)
}

func (s *PartyService) SendMessage(ctx context.Context, partyID, userID entities.UUID, message string, timestamp int) error {
	isInParty, err := s.partyRepo.IsUserInParty(ctx, partyID, userID)
	if err != nil {
		return fmt.Errorf("failed to check party membership: %w", err)
	}

	if !isInParty {
		return domainerrors.ErrUserNotInParty
	}

	partyMessage := &entities.PartyMessage{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(generateID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PartyID:   partyID,
		UserID:    userID,
		Message:   message,
		Timestamp: timestamp,
	}

	return s.partyRepo.CreateMessage(ctx, partyMessage)
}

func (s *PartyService) CreateDeviceTransfer(ctx context.Context, userID, fromDeviceID, toDeviceID, contentID entities.UUID, position int) (*entities.DeviceTransfer, error) {
	devices, err := s.deviceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user devices: %w", err)
	}

	fromExists, toExists := false, false
	for _, d := range devices {
		if d.ID == fromDeviceID {
			fromExists = true
		}
		if d.ID == toDeviceID {
			toExists = true
		}
	}

	if !fromExists || !toExists {
		return nil, domainerrors.ErrInvalidDeviceID
	}

	content, err := s.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		return nil, domainerrors.ErrContentNotFound
	}

	transfer := &entities.DeviceTransfer{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(generateID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:       userID,
		FromDeviceID: fromDeviceID,
		ToDeviceID:   toDeviceID,
		ContentType:  content.Type,
		ContentID:    contentID,
		Position:     position,
		IsCompleted:  false,
	}

	if err := s.transferRepo.Create(ctx, transfer); err != nil {
		return nil, fmt.Errorf("failed to create device transfer: %w", err)
	}

	return transfer, nil
}

func generateRoomCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		randB := make([]byte, 1)
		rand.Read(randB)
		b[i] = chars[int(randB[0])%len(chars)]
	}
	return string(b)
}
