package entity

import (
	"regexp"
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

func TestNewUUID(t *testing.T) {
	canonical := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	seen := map[UUID]bool{}
	for range 1000 {
		id := NewUUID()
		if !canonical.MatchString(string(id)) {
			t.Fatalf("%q is not a canonical v4 UUID", id)
		}
		if seen[id] {
			t.Fatalf("duplicate id %q", id)
		}
		seen[id] = true
	}
}
