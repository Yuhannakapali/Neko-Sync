package mediafile

import "testing"

func TestParse(t *testing.T) {
	i := func(n int) *int { return &n }
	f := func(n float64) *float64 { return &n }

	tests := []struct {
		name string
		path string
		kind Kind
		want Parsed
	}{
		// Movies
		{"folder with year", "The Life of David Gale (2003)/movie.mkv", KindMovie, Parsed{Title: "The Life of David Gale", Year: i(2003)}},
		{"release name", "The.Matrix.1999.1080p.BluRay.x264.mkv", KindMovie, Parsed{Title: "The Matrix", Year: i(1999)}},
		{"number in title", "Blade Runner 2049 (2017).mkv", KindMovie, Parsed{Title: "Blade Runner 2049", Year: i(2017)}},
		{"title starts with year", "2001 A Space Odyssey (1968).mkv", KindMovie, Parsed{Title: "2001 A Space Odyssey", Year: i(1968)}},
		{"title is a year", "1917.mkv", KindMovie, Parsed{Title: "1917"}},
		{"title is a year plus year", "1917 (2019).mkv", KindMovie, Parsed{Title: "1917", Year: i(2019)}},
		{"no year junk", "Heat.1080p.WEB-DL.mkv", KindMovie, Parsed{Title: "Heat"}},
		{"underscores", "Spirited_Away_(2001).mp4", KindMovie, Parsed{Title: "Spirited Away", Year: i(2001)}},

		// Series
		{"SxxEyy release", "Slow.Horses.S06E02.1080p.WEB.mkv", KindSeries, Parsed{Title: "Slow Horses", Season: i(6), Episode: f(2)}},
		{"season folder", "Breaking Bad (2008)/Season 02/Breaking Bad - S02E05 - Breakage.mkv", KindSeries, Parsed{Title: "Breaking Bad", Season: i(2), Episode: f(5)}},
		{"episode only in season folder", "Futurama/Season 14/E09.mkv", KindSeries, Parsed{Title: "Futurama", Season: i(14), Episode: f(9)}},
		{"1x02", "Dutton Ranch 1x02.mkv", KindSeries, Parsed{Title: "Dutton Ranch", Season: i(1), Episode: f(2)}},
		{"year in show folder", "Doctor Who (2005)/Season 1/S01E01.mkv", KindSeries, Parsed{Title: "Doctor Who", Year: i(2005), Season: i(1), Episode: f(1)}},
		{"specials folder", "Sherlock/Specials/S00E01.mkv", KindSeries, Parsed{Title: "Sherlock", Season: i(0), Episode: f(1)}},
		{"multi episode keeps first", "Show.S01E01E02.mkv", KindSeries, Parsed{Title: "Show", Season: i(1), Episode: f(1)}},

		// Anime
		{"fansub", "[SubsPlease] Sousou no Frieren - 07 (1080p) [ABCD1234].mkv", KindAnime, Parsed{Title: "Sousou no Frieren", Episode: f(7)}},
		{"fansub v2", "[Group] Dandadan - 12v2 [1080p].mkv", KindAnime, Parsed{Title: "Dandadan", Episode: f(12)}},
		{"half episode", "Gintama - 7.5.mkv", KindAnime, Parsed{Title: "Gintama", Episode: f(7.5)}},
		{"anime SxxEyy", "Mob Psycho 100 S02E03.mkv", KindAnime, Parsed{Title: "Mob Psycho 100", Season: i(2), Episode: f(3)}},
		{"anime trailing number in folder", "Hunter x Hunter (2011)/Hunter x Hunter 045.mkv", KindAnime, Parsed{Title: "Hunter x Hunter", Episode: f(45)}},
		{"loose anime movie", "Mob Psycho 100.mkv", KindAnime, Parsed{Title: "Mob Psycho 100"}},
		{"anime title from folder", "Code Geass R2/[Grp] 05 [720p].mkv", KindAnime, Parsed{Title: "Code Geass R2", Episode: f(5)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Parse(tc.path, tc.kind)
			if got.Title != tc.want.Title {
				t.Errorf("Title = %q, want %q", got.Title, tc.want.Title)
			}
			if !eqInt(got.Year, tc.want.Year) {
				t.Errorf("Year = %v, want %v", deref(got.Year), deref(tc.want.Year))
			}
			if !eqInt(got.Season, tc.want.Season) {
				t.Errorf("Season = %v, want %v", deref(got.Season), deref(tc.want.Season))
			}
			if !eqFloat(got.Episode, tc.want.Episode) {
				t.Errorf("Episode = %v, want %v", derefF(got.Episode), derefF(tc.want.Episode))
			}
		})
	}
}

func TestIsVideoFile(t *testing.T) {
	for path, want := range map[string]bool{
		"a.mkv": true, "a.MP4": true, "a.m2ts": true,
		"a.srt": false, "a.nfo": false, "a.jpg": false, "mkv": false,
	} {
		if got := IsVideoFile(path); got != want {
			t.Errorf("IsVideoFile(%q) = %v, want %v", path, got, want)
		}
	}
}

func eqInt(a, b *int) bool       { return (a == nil && b == nil) || (a != nil && b != nil && *a == *b) }
func eqFloat(a, b *float64) bool { return (a == nil && b == nil) || (a != nil && b != nil && *a == *b) }
func deref(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}
func derefF(p *float64) any {
	if p == nil {
		return nil
	}
	return *p
}

func TestNormalize(t *testing.T) {
	for in, want := range map[string]string{
		"The.Matrix.1999.1080p": "The Matrix 1999 1080p",
		"Gintama - 7.5":         "Gintama - 7.5",
		"Show.S01E07.5.mkv":     "Show S01E07.5 mkv",
		"Movie.DDP5.1.H.264":    "Movie DDP5.1 H 264",
		"a__b  c":               "a b c",
	} {
		if got := normalize(in); got != want {
			t.Errorf("normalize(%q) = %q, want %q", in, got, want)
		}
	}
}
