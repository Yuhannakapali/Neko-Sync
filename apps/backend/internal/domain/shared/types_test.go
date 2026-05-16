package shared

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
		t.Errorf("expected ID 'test-id', got %s", entity.ID)
	}
}
