// Package config reads Instance settings from the environment.
//
// Load returns an error instead of calling log.Fatal, so tests can run it
// (the Hub's config.go kills the test binary on a missing variable).
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nekosync-instance/internal/domain/mediafile"
)

// Library is one library declared in NEKO_LIBRARIES.
type Library struct {
	Name string
	Kind mediafile.Kind
	Root string
}

// Config holds all Instance settings.
type Config struct {
	Addr        string    // NEKO_ADDR, default ":8787"
	DBPath      string    // NEKO_DB_PATH, default "./data/instance.db"
	FFprobeBin  string    // NEKO_FFPROBE, default "ffprobe"
	ScanOnStart bool      // NEKO_SCAN_ON_START, default true
	Libraries   []Library // NEKO_LIBRARIES
}

// Load reads the environment.
//
// NEKO_LIBRARIES format: entries separated by ";", each "Name:kind:/absolute/path".
// Example: "Movies:movie:/media/movies;TV:series:/media/tv;Anime:anime:/media/anime"
// The path may contain ":" (only the first two are separators).
func Load() (*Config, error) {
	return load(os.Getenv)
}

func load(getenv func(string) string) (*Config, error) {
	cfg := &Config{
		Addr:        or(getenv("NEKO_ADDR"), ":8787"),
		DBPath:      or(getenv("NEKO_DB_PATH"), "./data/instance.db"),
		FFprobeBin:  or(getenv("NEKO_FFPROBE"), "ffprobe"),
		ScanOnStart: true,
	}
	switch strings.ToLower(getenv("NEKO_SCAN_ON_START")) {
	case "", "1", "true", "yes":
	case "0", "false", "no":
		cfg.ScanOnStart = false
	default:
		return nil, fmt.Errorf("NEKO_SCAN_ON_START: want true or false, got %q", getenv("NEKO_SCAN_ON_START"))
	}

	libs, err := parseLibraries(getenv("NEKO_LIBRARIES"))
	if err != nil {
		return nil, fmt.Errorf("NEKO_LIBRARIES: %w", err)
	}
	cfg.Libraries = libs
	return cfg, nil
}

func parseLibraries(s string) ([]Library, error) {
	var out []Library
	roots := map[string]bool{}
	for _, entry := range strings.Split(s, ";") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		parts := strings.SplitN(entry, ":", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("%q: want Name:kind:/path", entry)
		}
		lib := Library{
			Name: strings.TrimSpace(parts[0]),
			Kind: mediafile.Kind(strings.ToLower(strings.TrimSpace(parts[1]))),
			Root: filepath.Clean(strings.TrimSpace(parts[2])),
		}
		if lib.Name == "" {
			return nil, fmt.Errorf("%q: name is empty", entry)
		}
		if !lib.Kind.Valid() {
			return nil, fmt.Errorf("%q: kind must be movie, series, or anime", entry)
		}
		if !filepath.IsAbs(lib.Root) {
			return nil, fmt.Errorf("%q: path must be absolute", entry)
		}
		if roots[lib.Root] {
			return nil, fmt.Errorf("%q: path is used by another library", entry)
		}
		roots[lib.Root] = true
		out = append(out, lib)
	}
	return out, nil
}

func or(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
