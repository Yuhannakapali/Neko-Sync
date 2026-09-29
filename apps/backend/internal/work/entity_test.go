package work

import "testing"

func TestProviderID(t *testing.T) {
	w := &Work{ProviderIDs: map[string]string{"anilist": "154587", "mal": "52991"}}

	if id, ok := w.ProviderID("anilist"); !ok || id != "154587" {
		t.Fatalf("anilist: got (%q, %v), want (154587, true)", id, ok)
	}
	if _, ok := w.ProviderID("tmdb"); ok {
		t.Fatalf("tmdb: expected missing provider to return ok=false")
	}
}

func TestProviderID_NilMap(t *testing.T) {
	w := &Work{} // local-only work: no providers matched

	if _, ok := w.ProviderID("anilist"); ok {
		t.Fatalf("expected ok=false for a Work with nil ProviderIDs")
	}
}

func TestWorkChild_FractionalOrdinal(t *testing.T) {
	// Half-numbered entries (episode 7.5, chapter 10.5) must be representable.
	c := &WorkChild{Kind: ChildEpisode, Ordinal: 7.5}
	if c.Ordinal != 7.5 {
		t.Fatalf("ordinal: got %v, want 7.5", c.Ordinal)
	}
}
