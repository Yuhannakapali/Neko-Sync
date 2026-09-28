// Package scan walks a library folder and keeps MediaFile records in step with the
// files on disk: new files are probed and added, changed files are re-probed,
// unchanged files are skipped, and deleted files are removed.
package scan

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"nekosync-instance/internal/domain/library"
	"nekosync-instance/internal/domain/mediafile"
	"nekosync-instance/internal/domain/shared"
)

// Prober reads technical metadata from a file. probe.FFprobe implements it.
type Prober interface {
	Probe(ctx context.Context, path string) (mediafile.Technical, error)
}

// ErrScanRunning is returned when a scan of the same library is already running.
var ErrScanRunning = errors.New("scan already running for this library")

// ErrLibraryEmpty is returned when the walk finds no video files but the database
// still has some. This usually means the drive is not mounted, so the scanner
// refuses to delete every record.
var ErrLibraryEmpty = errors.New("library folder has no video files but records exist; refusing to remove them (is the drive mounted?)")

// Report summarizes one scan.
type Report struct {
	LibraryID shared.UUID   `json:"library_id"`
	Seen      int           `json:"seen"`
	Added     int           `json:"added"`
	Updated   int           `json:"updated"`
	Unchanged int           `json:"unchanged"`
	Removed   int           `json:"removed"`
	Failed    int           `json:"failed"`
	Took      time.Duration `json:"took"`
}

// Service runs library scans.
type Service struct {
	files  mediafile.Repository
	prober Prober
	log    *slog.Logger
	now    func() time.Time

	mu      sync.Mutex
	running map[shared.UUID]bool
}

// NewService builds a scan Service.
func NewService(files mediafile.Repository, prober Prober, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	return &Service{
		files:   files,
		prober:  prober,
		log:     log,
		now:     time.Now,
		running: map[shared.UUID]bool{},
	}
}

// Running reports whether a scan of the library is in progress.
func (s *Service) Running(id shared.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running[id]
}

// ScanLibrary scans one library. Only one scan per library runs at a time.
func (s *Service) ScanLibrary(ctx context.Context, lib *library.Library) (Report, error) {
	if !s.acquire(lib.ID) {
		return Report{}, ErrScanRunning
	}
	defer s.release(lib.ID)

	start := s.now()
	rep := Report{LibraryID: lib.ID}
	seen := map[string]bool{}

	walkErr := filepath.WalkDir(lib.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if path == lib.Root {
				return err // the root itself is unreadable: stop, and do not remove anything
			}
			s.log.Warn("scan: skip unreadable path", "path", path, "err", err)
			if d != nil && d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		name := d.Name()
		if d.IsDir() {
			if path != lib.Root && (strings.HasPrefix(name, ".") || strings.EqualFold(name, "@eaDir")) {
				return fs.SkipDir // hidden folders and Synology thumbnail folders
			}
			return nil
		}
		if strings.HasPrefix(name, ".") || !d.Type().IsRegular() || !mediafile.IsVideoFile(name) {
			return nil
		}

		seen[path] = true
		rep.Seen++
		switch outcome, err := s.scanFile(ctx, lib, path); {
		case err != nil:
			rep.Failed++
			s.log.Warn("scan: file failed", "path", path, "err", err)
		case outcome == added:
			rep.Added++
		case outcome == updated:
			rep.Updated++
		default:
			rep.Unchanged++
		}
		return nil
	})
	if walkErr != nil {
		return rep, fmt.Errorf("walk %q: %w", lib.Root, walkErr)
	}

	removed, err := s.removeMissing(ctx, lib, seen)
	rep.Removed = removed
	rep.Took = s.now().Sub(start)
	if err != nil {
		return rep, err
	}
	s.log.Info("scan: done", "library", lib.Name, "seen", rep.Seen, "added", rep.Added,
		"updated", rep.Updated, "unchanged", rep.Unchanged, "removed", rep.Removed,
		"failed", rep.Failed, "took", rep.Took)
	return rep, nil
}

type outcome int

const (
	unchanged outcome = iota
	added
	updated
)

func (s *Service) scanFile(ctx context.Context, lib *library.Library, path string) (outcome, error) {
	info, err := statFile(path)
	if err != nil {
		return unchanged, err
	}

	existing, err := s.files.GetByPath(ctx, path)
	if err != nil && !errors.Is(err, mediafile.ErrNotFound) {
		return unchanged, err
	}
	if existing != nil && existing.Unchanged(info.Size(), info.ModTime()) {
		return unchanged, nil
	}

	tech, err := s.prober.Probe(ctx, path)
	if err != nil {
		return unchanged, err
	}

	rel, err := filepath.Rel(lib.Root, path)
	if err != nil {
		rel = filepath.Base(path)
	}
	parsed := mediafile.Parse(rel, lib.Kind)
	now := s.now().UTC()

	f := &mediafile.MediaFile{
		LibraryID:   lib.ID,
		Path:        path,
		Size:        info.Size(),
		ModTime:     info.ModTime(),
		Kind:        lib.Kind,
		Technical:   tech,
		ParsedTitle: parsed.Title,
		ParsedYear:  parsed.Year,
		Season:      parsed.Season,
		Episode:     parsed.Episode,
		ProviderIDs: map[string]string{},
		MatchStatus: mediafile.MatchUnmatched,
	}
	result := added
	if existing != nil {
		result = updated
		f.ID = existing.ID
		f.CreatedAt = existing.CreatedAt
		// A manual match is the user's decision. A re-encode of the same file must
		// not undo it. Automatic matches are redone by the matcher.
		if existing.MatchStatus == mediafile.MatchManual {
			f.ProviderIDs = existing.ProviderIDs
			f.MatchStatus = mediafile.MatchManual
		}
	} else {
		f.ID = shared.NewUUID()
		f.CreatedAt = now
	}
	f.UpdatedAt = now

	if err := s.files.Save(ctx, f); err != nil {
		return unchanged, err
	}
	return result, nil
}

func (s *Service) removeMissing(ctx context.Context, lib *library.Library, seen map[string]bool) (int, error) {
	existing, err := s.files.ListByLibrary(ctx, lib.ID)
	if err != nil {
		return 0, err
	}
	if len(seen) == 0 && len(existing) > 0 {
		return 0, ErrLibraryEmpty
	}
	removed := 0
	for _, f := range existing {
		if seen[f.Path] {
			continue
		}
		if err := s.files.Delete(ctx, f.ID); err != nil {
			return removed, err
		}
		removed++
	}
	return removed, nil
}

func (s *Service) acquire(id shared.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running[id] {
		return false
	}
	s.running[id] = true
	return true
}

func (s *Service) release(id shared.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.running, id)
}
