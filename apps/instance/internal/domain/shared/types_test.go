package shared

import (
	"regexp"
	"testing"
)

var uuidV4 = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestNewUUID(t *testing.T) {
	seen := map[UUID]bool{}
	for i := 0; i < 1000; i++ {
		id := NewUUID()
		if !uuidV4.MatchString(string(id)) {
			t.Fatalf("not a v4 UUID: %q", id)
		}
		if seen[id] {
			t.Fatalf("duplicate UUID: %q", id)
		}
		seen[id] = true
	}
}
