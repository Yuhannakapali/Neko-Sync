package config

import (
	"testing"

	"nekosync-instance/internal/domain/mediafile"
)

func env(m map[string]string) func(string) string { return func(k string) string { return m[k] } }

func TestLoadDefaults(t *testing.T) {
	cfg, err := load(env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Addr != ":8787" || cfg.DBPath != "./data/instance.db" || !cfg.ScanOnStart || len(cfg.Libraries) != 0 {
		t.Fatalf("defaults = %+v", cfg)
	}
}

func TestLoadLibraries(t *testing.T) {
	cfg, err := load(env(map[string]string{
		"NEKO_LIBRARIES":     "Movies:movie:/media/movies ; TV:Series:/media/tv/;Anime:anime:/mnt/d:/anime",
		"NEKO_SCAN_ON_START": "false",
	}))
	if err != nil {
		t.Fatal(err)
	}
	want := []Library{
		{"Movies", mediafile.KindMovie, "/media/movies"},
		{"TV", mediafile.KindSeries, "/media/tv"},
		{"Anime", mediafile.KindAnime, "/mnt/d:/anime"},
	}
	if len(cfg.Libraries) != len(want) {
		t.Fatalf("got %+v", cfg.Libraries)
	}
	for i := range want {
		if cfg.Libraries[i] != want[i] {
			t.Errorf("library %d = %+v, want %+v", i, cfg.Libraries[i], want[i])
		}
	}
	if cfg.ScanOnStart {
		t.Error("ScanOnStart should be false")
	}
}

func TestLoadErrors(t *testing.T) {
	for name, libs := range map[string]string{
		"missing parts": "Movies:/media",
		"bad kind":      "Music:music:/media/music",
		"relative path": "Movies:movie:media/movies",
		"empty name":    ":movie:/media/movies",
		"duplicate":     "A:movie:/m;B:series:/m/",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := load(env(map[string]string{"NEKO_LIBRARIES": libs})); err == nil {
				t.Fatal("want error")
			}
		})
	}
	if _, err := load(env(map[string]string{"NEKO_SCAN_ON_START": "maybe"})); err == nil {
		t.Fatal("want error for bad bool")
	}
}
