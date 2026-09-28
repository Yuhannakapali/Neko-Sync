package probe

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"

	"nekosync-instance/internal/domain/mediafile"
)

const fixture = `{
  "streams": [
    {"index": 0, "codec_type": "video", "codec_name": "hevc", "width": 1920, "height": 1080,
     "disposition": {"default": 1, "forced": 0, "attached_pic": 0}},
    {"index": 1, "codec_type": "audio", "codec_name": "aac", "channels": 2,
     "tags": {"language": "jpn"}, "disposition": {"default": 1}},
    {"index": 2, "codec_type": "audio", "codec_name": "eac3", "channels": 6,
     "tags": {"LANGUAGE": "eng", "title": "English 5.1"}, "disposition": {"default": 0}},
    {"index": 3, "codec_type": "subtitle", "codec_name": "ass",
     "tags": {"language": "eng", "title": "Signs"}, "disposition": {"forced": 1}},
    {"index": 4, "codec_type": "attachment", "codec_name": "ttf"},
    {"index": 5, "codec_type": "video", "codec_name": "mjpeg", "disposition": {"attached_pic": 1}}
  ],
  "format": {"format_name": "matroska,webm", "duration": "1420.512000"}
}`

func TestParse(t *testing.T) {
	got, err := parse([]byte(fixture))
	if err != nil {
		t.Fatal(err)
	}
	if got.Container != "matroska,webm" || got.Duration != 1420.512 {
		t.Fatalf("format = %q %v", got.Container, got.Duration)
	}
	if len(got.Streams) != 4 {
		t.Fatalf("want 4 streams (attachment and cover art dropped), got %d: %+v", len(got.Streams), got.Streams)
	}
	want := []mediafile.Stream{
		{Index: 0, Type: mediafile.StreamVideo, Codec: "hevc", Width: 1920, Height: 1080, Default: true},
		{Index: 1, Type: mediafile.StreamAudio, Codec: "aac", Language: "jpn", Channels: 2, Default: true},
		{Index: 2, Type: mediafile.StreamAudio, Codec: "eac3", Language: "eng", Title: "English 5.1", Channels: 6},
		{Index: 3, Type: mediafile.StreamSubtitle, Codec: "ass", Language: "eng", Title: "Signs", Forced: true},
	}
	for i := range want {
		if got.Streams[i] != want[i] {
			t.Errorf("stream %d = %+v, want %+v", i, got.Streams[i], want[i])
		}
	}
}

func TestParseBadJSON(t *testing.T) {
	if _, err := parse([]byte("not json")); err == nil {
		t.Fatal("want error")
	}
}

// TestProbeRealFile generates a one-second clip with ffmpeg and probes it.
func TestProbeRealFile(t *testing.T) {
	if testing.Short() {
		t.Skip("needs ffmpeg")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed")
	}
	path := filepath.Join(t.TempDir(), "clip.mkv")
	gen := exec.Command("ffmpeg", "-v", "error",
		"-f", "lavfi", "-i", "testsrc=size=320x240:rate=10:duration=1",
		"-f", "lavfi", "-i", "sine=frequency=440:duration=1",
		"-metadata:s:a:0", "language=jpn",
		"-c:v", "mpeg4", "-c:a", "aac", path)
	if out, err := gen.CombinedOutput(); err != nil {
		t.Fatalf("ffmpeg: %v: %s", err, out)
	}

	got, err := FFprobe{}.Probe(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Duration < 0.9 || got.Duration > 1.2 {
		t.Errorf("duration = %v, want about 1s", got.Duration)
	}
	var video, audio int
	for _, s := range got.Streams {
		switch s.Type {
		case mediafile.StreamVideo:
			video++
			if s.Width != 320 || s.Height != 240 {
				t.Errorf("video size = %dx%d", s.Width, s.Height)
			}
		case mediafile.StreamAudio:
			audio++
			if s.Language != "jpn" {
				t.Errorf("audio language = %q", s.Language)
			}
		}
	}
	if video != 1 || audio != 1 {
		t.Errorf("streams: %d video, %d audio", video, audio)
	}
}

func TestProbeMissingFile(t *testing.T) {
	if _, err := exec.LookPath("ffprobe"); err != nil {
		t.Skip("ffprobe not installed")
	}
	if _, err := (FFprobe{}).Probe(context.Background(), "/does/not/exist.mkv"); err == nil {
		t.Fatal("want error")
	}
}
