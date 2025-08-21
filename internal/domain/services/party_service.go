package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"nekosync/internal/domain/entities"
	"nekosync/internal/domain/repositories"
	"time"
)

// PartyService contains business logic for watch party operations
type PartyService struct {
	partyRepo   repositories.PartyRepository
	contentRepo repositories.ContentRepository
	userRepo    repositories.UserRepository
}

// NewPartyService creates a new PartyService
func NewPartyService(partyRepo repositories.PartyRepository, contentRepo repositories.ContentRepository, userRepo repositories.UserRepository) *PartyService {
	return &PartyService{
		partyRepo:   partyRepo,
		contentRepo: contentRepo,
		userRepo:    userRepo,
	}
}

// CreateWatchParty creates a new watch party
func (s *PartyService) CreateWatchParty(ctx context.Context, hostUserID, contentID entities.UUID, title string, maxUsers int, isPrivate bool, password *string) (*entities.WatchParty, error) {
	// Validate host user exists
	_, err := s.userRepo.GetByID(ctx, hostUserID)
	if err != nil {
		return nil, fmt.Errorf("host user not found: %w", err)
	}

	// Validate content exists
	_, err = s.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("content not found: %w", err)
	}

	// Generate unique room code
	roomCode := s.generateRoomCode()

	party := &entities.WatchParty{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(s.generateUUID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		HostUserID: hostUserID,
		ContentID:  contentID,
		RoomCode:   roomCode,
		Title:      title,
		MaxUsers:   maxUsers,
		IsPrivate:  isPrivate,
		Password:   password,
		IsActive:   true,
	}

	err = s.partyRepo.Create(ctx, party)
	if err != nil {
		return nil, fmt.Errorf("failed to create watch party: %w", err)
	}

	// Add host as the first user
	partyUser := &entities.WatchPartyUser{
		PartyID:  party.ID,
		UserID:   hostUserID,
		JoinedAt: time.Now(),
		IsActive: true,
	}

	err = s.partyRepo.AddUser(ctx, partyUser)
	if err != nil {
		return nil, fmt.Errorf("failed to add host to party: %w", err)
	}

	return party, nil
}

// JoinWatchParty allows a user to join an existing watch party
func (s *PartyService) JoinWatchParty(ctx context.Context, userID entities.UUID, roomCode string, password *string) (*entities.WatchParty, error) {
	// Get party by room code
	party, err := s.partyRepo.GetByRoomCode(ctx, roomCode)
	if err != nil {
		return nil, fmt.Errorf("party not found")
	}

	if !party.IsActive {
		return nil, fmt.Errorf("party is not active")
	}

	// Check password if party is private
	if party.IsPrivate {
		if password == nil || party.Password == nil || *password != *party.Password {
			return nil, fmt.Errorf("invalid password")
		}
	}

	// Check if user is already in the party
	isInParty, err := s.partyRepo.IsUserInParty(ctx, party.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check party membership: %w", err)
	}

	if isInParty {
		return nil, fmt.Errorf("user is already in the party")
	}

	// Check if party is full
	userCount, err := s.partyRepo.GetPartyUserCount(ctx, party.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get party user count: %w", err)
	}

	if userCount >= party.MaxUsers {
		return nil, fmt.Errorf("party is full")
	}

	// Add user to party
	partyUser := &entities.WatchPartyUser{
		PartyID:  party.ID,
		UserID:   userID,
		JoinedAt: time.Now(),
		IsActive: true,
	}

	err = s.partyRepo.AddUser(ctx, partyUser)
	if err != nil {
		return nil, fmt.Errorf("failed to join party: %w", err)
	}

	return party, nil
}

// LeaveWatchParty allows a user to leave a watch party
func (s *PartyService) LeaveWatchParty(ctx context.Context, userID, partyID entities.UUID) error {
	party, err := s.partyRepo.GetByID(ctx, partyID)
	if err != nil {
		return fmt.Errorf("party not found: %w", err)
	}

	// Check if user is in the party
	isInParty, err := s.partyRepo.IsUserInParty(ctx, partyID, userID)
	if err != nil {
		return fmt.Errorf("failed to check party membership: %w", err)
	}

	if !isInParty {
		return fmt.Errorf("user is not in the party")
	}

	// Remove user from party
	err = s.partyRepo.RemoveUser(ctx, partyID, userID)
	if err != nil {
		return fmt.Errorf("failed to leave party: %w", err)
	}

	// If the host leaves, end the party
	if party.HostUserID == userID {
		party.End()
		err = s.partyRepo.Update(ctx, party)
		if err != nil {
			return fmt.Errorf("failed to end party: %w", err)
		}
	}

	return nil
}

// UpdatePartyState updates the playback state of a watch party
func (s *PartyService) UpdatePartyState(ctx context.Context, partyID, userID entities.UUID, currentTime int, isPlaying bool, playbackSpeed float64) error {
	party, err := s.partyRepo.GetByID(ctx, partyID)
	if err != nil {
		return fmt.Errorf("party not found: %w", err)
	}

	if !party.IsActive {
		return fmt.Errorf("party is not active")
	}

	// Check if user is in the party
	isInParty, err := s.partyRepo.IsUserInParty(ctx, partyID, userID)
	if err != nil {
		return fmt.Errorf("failed to check party membership: %w", err)
	}

	if !isInParty {
		return fmt.Errorf("user is not in the party")
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

// SendMessage sends a chat message in a watch party
func (s *PartyService) SendMessage(ctx context.Context, partyID, userID entities.UUID, message string, timestamp int) error {
	// Check if user is in the party
	isInParty, err := s.partyRepo.IsUserInParty(ctx, partyID, userID)
	if err != nil {
		return fmt.Errorf("failed to check party membership: %w", err)
	}

	if !isInParty {
		return fmt.Errorf("user is not in the party")
	}

	partyMessage := &entities.PartyMessage{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(s.generateUUID()),
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

// CreateDeviceTransfer initiates a content transfer between devices
func (s *PartyService) CreateDeviceTransfer(ctx context.Context, userID, fromDeviceID, toDeviceID, contentID entities.UUID, position int) (*entities.DeviceTransfer, error) {
	// Validate devices belong to the user
	devices, err := s.userRepo.GetDevicesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user devices: %w", err)
	}

	fromDeviceExists := false
	toDeviceExists := false
	for _, device := range devices {
		if device.ID == fromDeviceID {
			fromDeviceExists = true
		}
		if device.ID == toDeviceID {
			toDeviceExists = true
		}
	}

	if !fromDeviceExists || !toDeviceExists {
		return nil, fmt.Errorf("invalid device IDs")
	}

	// Get content to determine type
	content, err := s.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("content not found: %w", err)
	}

	transfer := &entities.DeviceTransfer{
		BaseEntity: entities.BaseEntity{
			ID:        entities.UUID(s.generateUUID()),
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

	err = s.partyRepo.CreateTransfer(ctx, transfer)
	if err != nil {
		return nil, fmt.Errorf("failed to create device transfer: %w", err)
	}

	return transfer, nil
}

// Private helper methods
func (s *PartyService) generateRoomCode() string {
	// Generate a 6-character room code
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[s.randomInt(len(chars))]
	}
	return string(code)
}

func (s *PartyService) generateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *PartyService) randomInt(max int) int {
	bytes := make([]byte, 1)
	rand.Read(bytes)
	return int(bytes[0]) % max
}
