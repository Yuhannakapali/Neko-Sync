package content

import "testing"

func TestContentType_Values(t *testing.T) {
	tests := []struct {
		name     string
		value    Type
		expected string
	}{
		{"anime", TypeAnime, "anime"},
		{"manga", TypeManga, "manga"},
		{"movie", TypeMovie, "movie"},
		{"music", TypeMusic, "music"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(tt.value))
			}
		})
	}
}
