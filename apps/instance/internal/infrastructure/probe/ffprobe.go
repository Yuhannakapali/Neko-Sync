// Package probe reads technical metadata from media files with ffprobe.
package probe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"nekosync-instance/internal/domain/mediafile"
)

// FFprobe runs the ffprobe binary. The zero value uses "ffprobe" from PATH and a
// 30-second timeout per file.
type FFprobe struct {
	Bin     string
	Timeout time.Duration
}

// Probe returns the container, duration, and streams of the file at path.
func (p FFprobe) Probe(ctx context.Context, path string) (mediafile.Technical, error) {
	bin := p.Bin
	if bin == "" {
		bin = "ffprobe"
	}
	timeout := p.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin,
		"-v", "error",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"--", path, // "--" stops a file name that starts with "-" being read as a flag
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return mediafile.Technical{}, fmt.Errorf("ffprobe %q: %w: %s", path, err, strings.TrimSpace(stderr.String()))
	}
	return parse(stdout.Bytes())
}

type ffprobeOutput struct {
	Format struct {
		FormatName string `json:"format_name"`
		Duration   string `json:"duration"`
	} `json:"format"`
	Streams []struct {
		Index       int               `json:"index"`
		CodecType   string            `json:"codec_type"`
		CodecName   string            `json:"codec_name"`
		Width       int               `json:"width"`
		Height      int               `json:"height"`
		Channels    int               `json:"channels"`
		Tags        map[string]string `json:"tags"`
		Disposition struct {
			Default     int `json:"default"`
			Forced      int `json:"forced"`
			AttachedPic int `json:"attached_pic"`
		} `json:"disposition"`
	} `json:"streams"`
}

func parse(data []byte) (mediafile.Technical, error) {
	var out ffprobeOutput
	if err := json.Unmarshal(data, &out); err != nil {
		return mediafile.Technical{}, fmt.Errorf("decode ffprobe output: %w", err)
	}

	t := mediafile.Technical{Container: out.Format.FormatName}
	if out.Format.Duration != "" {
		d, err := strconv.ParseFloat(out.Format.Duration, 64)
		if err != nil {
			return mediafile.Technical{}, fmt.Errorf("parse duration %q: %w", out.Format.Duration, err)
		}
		t.Duration = d
	}

	for _, s := range out.Streams {
		typ := mediafile.StreamType(s.CodecType)
		switch typ {
		case mediafile.StreamVideo:
			if s.Disposition.AttachedPic == 1 {
				continue // embedded cover art, not a playable video track
			}
		case mediafile.StreamAudio, mediafile.StreamSubtitle:
		default:
			continue // data and attachment streams (fonts in anime MKVs, chapters)
		}
		t.Streams = append(t.Streams, mediafile.Stream{
			Index:    s.Index,
			Type:     typ,
			Codec:    s.CodecName,
			Language: tag(s.Tags, "language"),
			Title:    tag(s.Tags, "title"),
			Width:    s.Width,
			Height:   s.Height,
			Channels: s.Channels,
			Default:  s.Disposition.Default == 1,
			Forced:   s.Disposition.Forced == 1,
		})
	}
	return t, nil
}

// tag reads a stream tag case-insensitively; MKV and MP4 muxers disagree on case.
func tag(tags map[string]string, key string) string {
	for k, v := range tags {
		if strings.EqualFold(k, key) {
			return v
		}
	}
	return ""
}
