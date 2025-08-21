package user

import (
	"testing"
	"nekosync/internal/domain/entities"
)

func TestUser_Creation(t *testing.T) {
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     entities.UserRoleUser,
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}
	
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", user.Email)
	}
	
	if user.Role != entities.UserRoleUser {
		t.Errorf("Expected role %s, got %s", entities.UserRoleUser, user.Role)
	}
}

func TestUserProfile_Creation(t *testing.T) {
	profile := &UserProfile{
		UserID: entities.UUID("test-user-id"),
	}
	
	if profile.UserID != "test-user-id" {
		t.Errorf("Expected UserID 'test-user-id', got %s", profile.UserID)
	}
}

func TestUserDevice_Creation(t *testing.T) {
	device := &UserDevice{
		DeviceName: "iPhone",
		Platform:   entities.PlatformTypeMobile,
		IsActive:   true,
	}
	
	if device.DeviceName != "iPhone" {
		t.Errorf("Expected device name 'iPhone', got %s", device.DeviceName)
	}
	
	if device.Platform != entities.PlatformTypeMobile {
		t.Errorf("Expected platform %s, got %s", entities.PlatformTypeMobile, device.Platform)
	}
	
	if !device.IsActive {
		t.Error("Expected device to be active")
	}
}

func TestUserFollow_Creation(t *testing.T) {
	follow := &UserFollow{
		FollowerID:  entities.UUID("follower-id"),
		FollowingID: entities.UUID("following-id"),
	}
	
	if follow.FollowerID != "follower-id" {
		t.Errorf("Expected FollowerID 'follower-id', got %s", follow.FollowerID)
	}
	
	if follow.FollowingID != "following-id" {
		t.Errorf("Expected FollowingID 'following-id', got %s", follow.FollowingID)
	}
}
