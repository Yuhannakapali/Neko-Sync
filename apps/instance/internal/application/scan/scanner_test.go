package scan

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"nekosync-instance/internal/domain/library"
	"nekosync-instance/internal/domain/mediafile"
	"nekosync-instance/internal/domain/shared"
)

// memRepo is an in-memory mediafile.Repository.
type memRepo struct {
	mu    sync.Mutex
	files map[shared.UUID]*mediafile.MediaFile
}

func newMemRepo() *memRepo { return &memRepo{files: map[shared.UUID]*mediafile.MediaFile{}} }

func (r *memRepo) Get(_ context.Context, id shared.UUID) (*mediafile.MediaFile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if f, ok := r.files[id]; ok {
		c := *f
		return &c, nil
	}
	return nil, mediafile.ErrNotFound
}

func (r *memRepo) GetByPath(_ context.Context, path string) (*mediafile.MediaFile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, f := range r.files {
		if f.Path == path {
			c := *f
			return &c, nil
		}
	}
	return nil, mediafile.ErrNotFound
}

func (r *memRepo) ListByLibrary(_ context.Context, id shared.UUID) ([]*mediafile.MediaFile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []*mediafile.MediaFile
	for _, f := range r.files {
		if f.LibraryID == id {
			c := *f
			out = append(out, &c)
		}
	}
	return out, nil
}

func (r *memRepo) Save(_ context.Context, f *mediafile.MediaFile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c := *f
	r.files[f.ID] = &c
	return nil
}

func (r *memRepo) Delete(_ context.Context, id shared.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.files, id)
	return nil
}

// fakeProber counts calls and fails for paths in failOn.
type fakeProber struct {
	mu     sync.Mutex
	calls  int
	failOn map[string]bool
}

func (p *fakeProber) Probe(_ context.Context, path string) (mediafile.Technical, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.calls++
	if p.failOn[path] {
		return mediafile.Technical{}, errors.New("corrupt file")
	}
	return mediafile.Technical{Container: "matroska,webm", Duration: 60}, nil
}

func write(t *testing.T, root, rel string) string {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func setup(t *testing.T, kind mediafile.Kind) (*Service, *memRepo, *fakeProber, *library.Library) {
	t.Helper()
	repo := newMemRepo()
	prober := &fakeProber{failOn: map[string]bool{}}
	lib := &library.Library{Name: "test", Kind: kind, Root: t.TempDir()}
	lib.ID = shared.NewUUID()
	return NewService(repo, prober, nil), repo, prober, lib
}

func TestScanAddsParsesAndSkips(t *testing.T) {
	svc, repo, prober, lib := setup(t, mediafile.KindSeries)
	ctx := context.Background()

	write(t, lib.Root, "Slow Horses/Season 06/Slow.Horses.S06E02.1080p.mkv")
	write(t, lib.Root, "Slow Horses/Season 06/Slow.Horses.S06E02.en.srt") // not video
	write(t, lib.Root, ".hidden/secret.mkv")                              // hidden folder
	write(t, lib.Root, "Slow Horses/._S06E03.mkv")                        // macOS resource fork

	rep, err := svc.ScanLibrary(ctx, lib)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Seen != 1 || rep.Added != 1 {
		t.Fatalf("report = %+v", rep)
	}
	files, _ := repo.ListByLibrary(ctx, lib.ID)
	f := files[0]
	if f.ParsedTitle != "Slow Horses" || *f.Season != 6 || *f.Episode != 2 {
		t.Errorf("parsed = %q S%v E%v", f.ParsedTitle, *f.Season, *f.Episode)
	}
	if f.MatchStatus != mediafile.MatchUnmatched || f.Container != "matroska,webm" {
		t.Errorf("file = %+v", f)
	}

	// Second scan: nothing changed, so no probe.
	rep, err = svc.ScanLibrary(ctx, lib)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Unchanged != 1 || prober.calls != 1 {
		t.Fatalf("second scan report = %+v, probe calls = %d", rep, prober.calls)
	}
}

func TestScanUpdatesChangedFileAndKeepsManualMatch(t *testing.T) {
	svc, repo, _, lib := setup(t, mediafile.KindMovie)
	ctx := context.Background()
	path := write(t, lib.Root, "Heat (1995)/Heat.mkv")

	if _, err := svc.ScanLibrary(ctx, lib); err != nil {
		t.Fatal(err)
	}
	f, _ := repo.GetByPath(ctx, path)
	id := f.ID
	f.MatchStatus = mediafile.MatchManual
	f.ProviderIDs = map[string]string{"tmdb": "949"}
	_ = repo.Save(ctx, f)

	// Replace the file with a new encode.
	if err := os.WriteFile(path, []byte("bigger file"), 0o644); err != nil {
		t.Fatal(err)
	}
	later := time.Now().Add(time.Minute)
	_ = os.Chtimes(path, later, later)

	rep, err := svc.ScanLibrary(ctx, lib)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Updated != 1 {
		t.Fatalf("report = %+v", rep)
	}
	f, _ = repo.GetByPath(ctx, path)
	if f.ID != id {
		t.Error("ID changed on update")
	}
	if f.MatchStatus != mediafile.MatchManual || f.ProviderIDs["tmdb"] != "949" {
		t.Errorf("manual match lost: %v %v", f.MatchStatus, f.ProviderIDs)
	}
}

func TestScanRemovesDeletedFiles(t *testing.T) {
	svc, _, _, lib := setup(t, mediafile.KindMovie)
	ctx := context.Background()
	keep := write(t, lib.Root, "A (2001).mkv")
	gone := write(t, lib.Root, "B (2002).mkv")
	_ = keep

	if _, err := svc.ScanLibrary(ctx, lib); err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(gone)
	rep, err := svc.ScanLibrary(ctx, lib)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Removed != 1 || rep.Unchanged != 1 {
		t.Fatalf("report = %+v", rep)
	}
}

func TestScanRefusesToWipeEmptyLibrary(t *testing.T) {
	svc, repo, _, lib := setup(t, mediafile.KindMovie)
	ctx := context.Background()
	path := write(t, lib.Root, "A (2001).mkv")
	if _, err := svc.ScanLibrary(ctx, lib); err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(path) // looks like an unmounted drive: folder exists, but empty
	if _, err := svc.ScanLibrary(ctx, lib); !errors.Is(err, ErrLibraryEmpty) {
		t.Fatalf("err = %v, want ErrLibraryEmpty", err)
	}
	if files, _ := repo.ListByLibrary(ctx, lib.ID); len(files) != 1 {
		t.Fatalf("records were deleted: %d left", len(files))
	}
}

func TestScanMissingRootRemovesNothing(t *testing.T) {
	svc, repo, _, lib := setup(t, mediafile.KindMovie)
	ctx := context.Background()
	write(t, lib.Root, "A (2001).mkv")
	if _, err := svc.ScanLibrary(ctx, lib); err != nil {
		t.Fatal(err)
	}

	lib.Root = filepath.Join(lib.Root, "unmounted")
	if _, err := svc.ScanLibrary(ctx, lib); err == nil {
		t.Fatal("want error for missing root")
	}
	if files, _ := repo.ListByLibrary(ctx, lib.ID); len(files) != 1 {
		t.Fatalf("records were deleted: %d left", len(files))
	}
}

func TestScanCountsProbeFailures(t *testing.T) {
	svc, _, prober, lib := setup(t, mediafile.KindMovie)
	bad := write(t, lib.Root, "Broken (2001).mkv")
	write(t, lib.Root, "Good (2002).mkv")
	prober.failOn[bad] = true

	rep, err := svc.ScanLibrary(context.Background(), lib)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Failed != 1 || rep.Added != 1 {
		t.Fatalf("report = %+v", rep)
	}
}

func TestScanOneAtATime(t *testing.T) {
	svc, _, _, lib := setup(t, mediafile.KindMovie)
	if !svc.acquire(lib.ID) {
		t.Fatal("first acquire failed")
	}
	if _, err := svc.ScanLibrary(context.Background(), lib); !errors.Is(err, ErrScanRunning) {
		t.Fatalf("err = %v, want ErrScanRunning", err)
	}
	svc.release(lib.ID)
	if svc.Running(lib.ID) {
		t.Fatal("still running after release")
	}
}
