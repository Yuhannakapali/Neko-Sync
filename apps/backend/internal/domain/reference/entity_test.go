package reference

import "testing"

func srcs(refs []*ContentReference) []SourceKind {
	out := make([]SourceKind, len(refs))
	for i, r := range refs {
		out[i] = r.Source
	}
	return out
}

func TestRank_Ordering(t *testing.T) {
	in := []*ContentReference{
		{Source: SourceOfficialLink},
		{Source: SourceInstance},
		{Source: SourcePublicDomain},
		{Source: SourceLocalFile},
	}
	got := srcs(Rank(in))
	want := []SourceKind{SourceInstance, SourceLocalFile, SourcePublicDomain, SourceOfficialLink}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Rank order = %v, want %v", got, want)
		}
	}
}

func TestRank_StableWithinSource(t *testing.T) {
	// Two instance refs (e.g. 4K then 1080p) must keep their input order so a
	// quality pre-sort survives ranking.
	first := &ContentReference{Source: SourceInstance, Locator: "4k"}
	second := &ContentReference{Source: SourceInstance, Locator: "1080p"}
	got := Rank([]*ContentReference{first, second})

	if got[0] != first || got[1] != second {
		t.Fatalf("Rank is not stable within a source kind")
	}
}

func TestRank_EmptyAndDoesNotMutateInput(t *testing.T) {
	if len(Rank(nil)) != 0 {
		t.Fatalf("Rank(nil) should be empty")
	}

	in := []*ContentReference{{Source: SourceOfficialLink}, {Source: SourceInstance}}
	_ = Rank(in)
	if in[0].Source != SourceOfficialLink {
		t.Fatalf("Rank must not mutate its input slice")
	}
}
