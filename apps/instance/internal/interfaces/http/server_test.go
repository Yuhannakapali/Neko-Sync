package http

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"nekosync-instance/internal/application/scan"
	"nekosync-instance/internal/domain/library"
	"nekosync-instance/internal/domain/mediafile"
	"nekosync-instance/internal/domain/shared"
)

type libRepo struct {
	libs map[shared.UUID]*library.Library
}

func (r *libRepo) Get(_ context.Context, id shared.UUID) (*library.Library, error) {
	if l, ok := r.libs[id]; ok {
		return l, nil
	}
	return nil, library.ErrNotFound
}
func (r *libRepo) GetByRoot(context.Context, string) (*library.Library, error) {
	return nil, library.ErrNotFound
}
func (r *libRepo) List(context.Context) ([]*library.Library, error) {
	var out []*library.Library
	for _, l := range r.libs {
		out = append(out, l)
	}
	return out, nil
}
func (r *libRepo) Save(_ context.Context, l *library.Library) error { r.libs[l.ID] = l; return nil }

type fileRepo struct {
	mu    sync.Mutex
	files map[shared.UUID]*mediafile.MediaFile
}

func (r *fileRepo) Get(_ context.Context, id shared.UUID) (*mediafile.MediaFile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if f, ok := r.files[id]; ok {
		return f, nil
	}
	return nil, mediafile.ErrNotFound
}
func (r *fileRepo) GetByPath(context.Context, string) (*mediafile.MediaFile, error) {
	return nil, mediafile.ErrNotFound
}
func (r *fileRepo) ListByLibrary(_ context.Context, id shared.UUID) ([]*mediafile.MediaFile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []*mediafile.MediaFile
	for _, f := range r.files {
		if f.LibraryID == id {
			out = append(out, f)
		}
	}
	return out, nil
}
func (r *fileRepo) Save(_ context.Context, f *mediafile.MediaFile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.files[f.ID] = f
	return nil
}
func (r *fileRepo) Delete(context.Context, shared.UUID) error { return nil }

type noProbe struct{}

func (noProbe) Probe(context.Context, string) (mediafile.Technical, error) {
	return mediafile.Technical{}, nil
}

func newTestServer(t *testing.T) (*httptest.Server, *library.Library, *fileRepo) {
	t.Helper()
	lib := &library.Library{Name: "Movies", Kind: mediafile.KindMovie, Root: t.TempDir()}
	lib.ID = shared.NewUUID()
	libs := &libRepo{libs: map[shared.UUID]*library.Library{lib.ID: lib}}
	files := &fileRepo{files: map[shared.UUID]*mediafile.MediaFile{}}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := NewServer(context.Background(), libs, files, scan.NewService(files, noProbe{}, log), log)
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return ts, lib, files
}

func get(t *testing.T, url string, into any) int {
	t.Helper()
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if into != nil {
		if err := json.NewDecoder(res.Body).Decode(into); err != nil {
			t.Fatal(err)
		}
	}
	return res.StatusCode
}

func TestRoutes(t *testing.T) {
	ts, lib, files := newTestServer(t)

	var health map[string]string
	if code := get(t, ts.URL+"/health", &health); code != 200 || health["service"] != "nekosync-instance" {
		t.Fatalf("health = %d %v", code, health)
	}

	var libs []map[string]any
	if code := get(t, ts.URL+"/api/libraries", &libs); code != 200 || len(libs) != 1 || libs[0]["scanning"] != false {
		t.Fatalf("libraries = %d %v", code, libs)
	}

	var list []any
	if code := get(t, ts.URL+"/api/libraries/"+string(lib.ID)+"/files", &list); code != 200 || len(list) != 0 {
		t.Fatalf("empty files = %d %v", code, list)
	}

	f := &mediafile.MediaFile{LibraryID: lib.ID, ParsedTitle: "Heat"}
	f.ID = shared.NewUUID()
	_ = files.Save(context.Background(), f)
	var got mediafile.MediaFile
	if code := get(t, ts.URL+"/api/files/"+string(f.ID), &got); code != 200 || got.ParsedTitle != "Heat" {
		t.Fatalf("file = %d %+v", code, got)
	}

	if code := get(t, ts.URL+"/api/files/nope", nil); code != 404 {
		t.Fatalf("missing file = %d", code)
	}
	if code := get(t, ts.URL+"/api/libraries/nope/files", nil); code != 404 {
		t.Fatalf("missing library = %d", code)
	}
	if res, _ := http.Post(ts.URL+"/health", "", nil); res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("POST /health = %d", res.StatusCode)
	}
}

func TestStartScan(t *testing.T) {
	ts, lib, _ := newTestServer(t)
	res, err := http.Post(ts.URL+"/api/libraries/"+string(lib.ID)+"/scan", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("scan = %d", res.StatusCode)
	}
	time.Sleep(50 * time.Millisecond) // let the background scan finish before cleanup
}
