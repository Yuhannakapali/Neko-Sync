package entities

import (
	"testing"
	"time"
)

func TestBaseEntity_Creation(t *testing.T) {
	entity := BaseEntity{
		ID:        UUID("test-id"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if entity.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", entity.ID)
	}
}

func TestContentType_Values(t *testing.T) {
	tests := []struct {
		name     string
		value    ContentType
		expected string
	}{
		{"anime", ContentTypeAnime, "anime"},
		{"manga", ContentTypeManga, "manga"},
		{"movie", ContentTypeMovie, "movie"},
		{"music", ContentTypeMusic, "music"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.value))
			}
		})
	}
}

func TestUserRole_Values(t *testing.T) {
	tests := []struct {
		name     string
		value    UserRole
		expected string
	}{
		{"user", UserRoleUser, "user"},
		{"admin", UserRoleAdmin, "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.value))
			}
		})
	}
}

func TestPlatformType_Values(t *testing.T) {
	tests := []struct {
		name     string
		value    PlatformType
		expected string
	}{
		{"web", PlatformTypeWeb, "web"},
		{"desktop", PlatformTypeDesktop, "desktop"},
		{"mobile", PlatformTypeMobile, "mobile"},
		{"tablet", PlatformTypeTablet, "tablet"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.value))
			}
		})
	}
}
