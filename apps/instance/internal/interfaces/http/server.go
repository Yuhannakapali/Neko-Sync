// Package http exposes the Instance API.
//
// It uses only net/http (Go 1.22+ method and path patterns), so the open-source
// Instance has no web framework dependency. The Hub uses Echo; the two never share
// handlers.
//
// Routes (all JSON):
//
//	GET  /health
//	GET  /api/libraries
//	GET  /api/libraries/{id}/files
//	POST /api/libraries/{id}/scan      202 Accepted, scan runs in the background
//	GET  /api/files/{id}
//
// There is no auth yet. Bind to localhost or your LAN only until signed stream URLs
// and API tokens land (see apps/instance/CLAUDE.md, "Remaining work").
package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"nekosync-instance/internal/application/scan"
	"nekosync-instance/internal/domain/library"
	"nekosync-instance/internal/domain/mediafile"
	"nekosync-instance/internal/domain/shared"
)

// Server holds the HTTP dependencies.
type Server struct {
	libs    library.Repository
	files   mediafile.Repository
	scanner *scan.Service
	log     *slog.Logger
	// bg is the context for background scans; it is cancelled on shutdown.
	bg context.Context
}

// NewServer builds the API handler. bg is cancelled when the process shuts down.
func NewServer(bg context.Context, libs library.Repository, files mediafile.Repository, scanner *scan.Service, log *slog.Logger) *Server {
	return &Server{libs: libs, files: files, scanner: scanner, log: log, bg: bg}
}

// Handler returns the routed http.Handler.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /api/libraries", s.listLibraries)
	mux.HandleFunc("GET /api/libraries/{id}/files", s.listFiles)
	mux.HandleFunc("POST /api/libraries/{id}/scan", s.startScan)
	mux.HandleFunc("GET /api/files/{id}", s.getFile)
	return s.recoverer(mux)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "nekosync-instance"})
}

func (s *Server) listLibraries(w http.ResponseWriter, r *http.Request) {
	libs, err := s.libs.List(r.Context())
	if err != nil {
		s.fail(w, err)
		return
	}
	type item struct {
		*library.Library
		Scanning bool `json:"scanning"`
	}
	out := make([]item, 0, len(libs))
	for _, l := range libs {
		out = append(out, item{l, s.scanner.Running(l.ID)})
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) listFiles(w http.ResponseWriter, r *http.Request) {
	id := shared.UUID(r.PathValue("id"))
	if _, err := s.libs.Get(r.Context(), id); err != nil {
		s.fail(w, err)
		return
	}
	files, err := s.files.ListByLibrary(r.Context(), id)
	if err != nil {
		s.fail(w, err)
		return
	}
	if files == nil {
		files = []*mediafile.MediaFile{}
	}
	writeJSON(w, http.StatusOK, files)
}

func (s *Server) startScan(w http.ResponseWriter, r *http.Request) {
	lib, err := s.libs.Get(r.Context(), shared.UUID(r.PathValue("id")))
	if err != nil {
		s.fail(w, err)
		return
	}
	if s.scanner.Running(lib.ID) {
		writeError(w, http.StatusConflict, scan.ErrScanRunning.Error())
		return
	}
	go func() {
		if _, err := s.scanner.ScanLibrary(s.bg, lib); err != nil {
			s.log.Error("scan failed", "library", lib.Name, "err", err)
		}
	}()
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "scanning", "library_id": string(lib.ID)})
}

func (s *Server) getFile(w http.ResponseWriter, r *http.Request) {
	f, err := s.files.Get(r.Context(), shared.UUID(r.PathValue("id")))
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, f)
}

func (s *Server) fail(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, library.ErrNotFound), errors.Is(err, mediafile.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		s.log.Error("request failed", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func (s *Server) recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			if v := recover(); v != nil {
				s.log.Error("panic", "path", r.URL.Path, "panic", v)
				writeError(w, http.StatusInternalServerError, "internal error")
			}
			s.log.Debug("request", "method", r.Method, "path", r.URL.Path, "took", time.Since(start))
		}()
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
