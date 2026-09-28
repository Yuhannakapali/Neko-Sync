package mediafile

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Parsed is what the file path tells us before any metadata lookup.
type Parsed struct {
	Title   string
	Year    *int
	Season  *int
	Episode *float64
}

var (
	// A year needs a delimiter on both sides so "1080p" and "x264" never match.
	yearRe = regexp.MustCompile(`(?:^|[\s(\[])((?:19|20)\d{2})(?:[\s)\]]|$)`)

	// Release-group noise. Everything from the first match on is dropped from a title.
	junkRe = regexp.MustCompile(`(?i)(?:^|\s)(?:2160p|1080p|720p|576p|480p|4k|uhd|blu-?ray|bdrip|brrip|web-?dl|webrip|web|hdtv|dvdrip|hdrip|x264|x265|h 264|h 265|h264|h265|hevc|avc|remux|proper|repack|extended|unrated|hdr10|hdr|10bit|8bit|aac|ac3|eac3|dts|atmos|multi|dual audio)(?:\s|$)`)

	// S01E02, s1e2, S01 E02, S01E07.5
	seRe = regexp.MustCompile(`(?i)(?:^|\s)s(\d{1,2})\s?e(\d{1,4}(?:\.\d)?)(?:\s|$|e\d)`)
	// 1x02
	xRe = regexp.MustCompile(`(?i)(?:^|\s)(\d{1,2})x(\d{2,3})(?:\s|$)`)
	// E05, Ep 05, Episode 05 (only trusted inside a season folder)
	epOnlyRe = regexp.MustCompile(`(?i)(?:^|\s)(?:e|ep|episode)\s?(\d{1,4}(?:\.\d)?)(?:\s|$)`)
	// "Season 1", "Season 01", "S01", "Series 2"
	seasonDirRe = regexp.MustCompile(`(?i)^(?:season|series|s)\s?(\d{1,2})$`)

	// Fansub style: "Title - 07", "Title - 07v2", "Title - 7.5"
	animeDashRe = regexp.MustCompile(`^(.*?)\s+-\s+(\d{1,4}(?:\.\d)?)(?:v\d+)?(?:\s|$)`)
	// Fallback: "Title 07" at the very end.
	animeTailRe = regexp.MustCompile(`^(.*\S)\s+(\d{1,4})(?:v\d+)?$`)
	// A file named only by its number: "05.mkv", "[Grp] 05 [720p].mkv".
	animeBareRe = regexp.MustCompile(`^(\d{1,4}(?:\.\d)?)(?:v\d+)?$`)
	bracketRe   = regexp.MustCompile(`\[[^\]]*\]|\([^)]*\)|\{[^}]*\}`)

	spaceRe = regexp.MustCompile(`\s+`)
)

// Parse extracts title, year, season, and episode from a path. relPath is the path
// relative to the library root, so parent folders can be used as a fallback title.
// Parse always returns a title (it falls back to the raw file name). Season and
// Episode are nil when the path does not contain them.
func Parse(relPath string, kind Kind) Parsed {
	relPath = filepath.ToSlash(relPath)
	dirs := strings.Split(filepath.ToSlash(filepath.Dir(relPath)), "/")
	if len(dirs) == 1 && dirs[0] == "." {
		dirs = nil
	}
	stem := strings.TrimSuffix(filepath.Base(relPath), filepath.Ext(relPath))

	var p Parsed
	switch kind {
	case KindMovie:
		p = parseMovie(stem, dirs)
	case KindAnime:
		p = parseAnime(stem, dirs)
	default:
		p = parseSeries(stem, dirs)
	}
	if p.Title == "" {
		p.Title = normalize(stem)
	}
	return p
}

func parseMovie(stem string, dirs []string) Parsed {
	// Jellyfin/Plex style "Title (Year)/anything.mkv": the folder is more reliable.
	if len(dirs) > 0 {
		if t, y := titleYear(dirs[len(dirs)-1]); y != nil && t != "" {
			return Parsed{Title: t, Year: y}
		}
	}
	t, y := titleYear(stem)
	return Parsed{Title: t, Year: y}
}

func parseSeries(stem string, dirs []string) Parsed {
	name := normalize(stem)
	var p Parsed
	var titlePart string

	if m := seRe.FindStringSubmatchIndex(name); m != nil {
		p.Season = intPtr(name[m[2]:m[3]])
		p.Episode = floatPtr(name[m[4]:m[5]])
		titlePart = name[:m[0]]
	} else if m := xRe.FindStringSubmatchIndex(name); m != nil {
		p.Season = intPtr(name[m[2]:m[3]])
		p.Episode = floatPtr(name[m[4]:m[5]])
		titlePart = name[:m[0]]
	} else if s, ok := seasonFromDirs(dirs); ok {
		if m := epOnlyRe.FindStringSubmatchIndex(name); m != nil {
			p.Season = &s
			p.Episode = floatPtr(name[m[2]:m[3]])
			titlePart = name[:m[0]]
		}
	}

	p.Title, p.Year = titleYear(titlePart)
	if p.Title == "" {
		p.Title, p.Year = showFromDirs(dirs)
	}
	if p.Season == nil {
		if s, ok := seasonFromDirs(dirs); ok && p.Episode != nil {
			p.Season = &s
		}
	}
	return p
}

func parseAnime(stem string, dirs []string) Parsed {
	// Many anime releases also use SxxEyy. Prefer that when present.
	if p := parseSeries(stem, dirs); p.Episode != nil {
		return p
	}

	name := normalize(bracketRe.ReplaceAllString(stem, " "))
	var p Parsed
	if m := animeDashRe.FindStringSubmatch(name); m != nil {
		p.Title = cleanTitle(m[1])
		p.Episode = floatPtr(m[2])
	} else if m := animeTailRe.FindStringSubmatch(name); m != nil && len(dirs) > 0 {
		// Only trust a bare trailing number when the file sits in a show folder,
		// so a loose "Mob Psycho 100.mkv" is not read as episode 100.
		p.Title = cleanTitle(m[1])
		p.Episode = floatPtr(m[2])
	} else if m := animeBareRe.FindStringSubmatch(name); m != nil && len(dirs) > 0 {
		p.Episode = floatPtr(m[1])
	}
	if p.Title == "" {
		p.Title, p.Year = showFromDirs(dirs)
	}
	if s, ok := seasonFromDirs(dirs); ok && p.Episode != nil {
		p.Season = &s
	}
	return p
}

// titleYear splits "Blade Runner 2049 (2017)" into ("Blade Runner 2049", 2017).
// The last year that has text before it wins, so a title that starts or ends with
// a number keeps it.
func titleYear(s string) (string, *int) {
	name := normalize(s)
	matches := yearRe.FindAllStringSubmatchIndex(name, -1)
	for i := len(matches) - 1; i >= 0; i-- {
		m := matches[i]
		if title := cleanTitle(name[:m[0]]); title != "" {
			return title, intPtr(name[m[2]:m[3]])
		}
	}
	if loc := junkRe.FindStringIndex(name); loc != nil {
		name = name[:loc[0]]
	}
	return cleanTitle(name), nil
}

// showFromDirs returns the show title from the folder path, skipping a trailing
// season folder ("Show (2019)/Season 01/file.mkv" gives "Show", 2019).
func showFromDirs(dirs []string) (string, *int) {
	for i := len(dirs) - 1; i >= 0; i-- {
		d := normalize(dirs[i])
		if seasonDirRe.MatchString(d) || strings.EqualFold(d, "specials") {
			continue
		}
		return titleYear(bracketRe.ReplaceAllStringFunc(dirs[i], keepYear))
	}
	return "", nil
}

func seasonFromDirs(dirs []string) (int, bool) {
	if len(dirs) == 0 {
		return 0, false
	}
	d := normalize(dirs[len(dirs)-1])
	if strings.EqualFold(d, "specials") {
		return 0, true
	}
	if m := seasonDirRe.FindStringSubmatch(d); m != nil {
		n, _ := strconv.Atoi(m[1])
		return n, true
	}
	return 0, false
}

// keepYear drops bracketed tags such as "[1080p]" but keeps "(2019)".
func keepYear(tag string) string {
	inner := strings.Trim(tag, "[](){}")
	if len(inner) == 4 && (strings.HasPrefix(inner, "19") || strings.HasPrefix(inner, "20")) {
		if _, err := strconv.Atoi(inner); err == nil {
			return "(" + inner + ")"
		}
	}
	return " "
}

// normalize turns release-style separators into spaces: "The.Matrix_1999" -> "The Matrix 1999".
// A dot is kept when it is a decimal point in a half episode ("7.5"): a digit before
// it and exactly one digit after it.
func normalize(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	b := []byte(s)
	for i, c := range b {
		if c != '.' {
			continue
		}
		decimal := i > 0 && isDigit(b[i-1]) &&
			i+1 < len(b) && isDigit(b[i+1]) &&
			(i+2 == len(b) || !isDigit(b[i+2]))
		if !decimal {
			b[i] = ' '
		}
	}
	return strings.TrimSpace(spaceRe.ReplaceAllString(string(b), " "))
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func cleanTitle(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimRight(s, " -([{")
	s = strings.TrimLeft(s, " -)]}")
	return strings.TrimSpace(spaceRe.ReplaceAllString(s, " "))
}

func intPtr(s string) *int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &n
}

func floatPtr(s string) *float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &f
}
