package user

import "testing"

func TestRole_Values(t *testing.T) {
	tests := []struct {
		name     string
		value    Role
		expected string
	}{
		{"user", RoleUser, "user"},
		{"admin", RoleAdmin, "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(tt.value))
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
		{"web", PlatformWeb, "web"},
		{"desktop", PlatformDesktop, "desktop"},
		{"mobile", PlatformMobile, "mobile"},
		{"tablet", PlatformTablet, "tablet"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(tt.value))
			}
		})
	}
}
