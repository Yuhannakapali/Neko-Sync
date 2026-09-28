// Command nekosync-instance is the self-hosted Neko-Sync media server.
//
// It scans library folders, stores MediaFiles in SQLite, and serves the Instance
// API. It runs fully without the Hub.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nekosync-instance/internal/application/scan"
	"nekosync-instance/internal/config"
	"nekosync-instance/internal/domain/library"
	"nekosync-instance/internal/domain/shared"
	"nekosync-instance/internal/infrastructure/probe"
	"nekosync-instance/internal/infrastructure/sqlite"
	api "nekosync-instance/internal/interfaces/http"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	if err := run(log); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := sqlite.Open(ctx, cfg.DBPath)
	if err != nil {
		return err
	}
	defer db.Close()

	libs := sqlite.NewLibraryRepository(db)
	files := sqlite.NewMediaFileRepository(db)
	scanner := scan.NewService(files, probe.FFprobe{Bin: cfg.FFprobeBin}, log)

	if err := syncLibraries(ctx, libs, cfg.Libraries, log); err != nil {
		return err
	}
	if cfg.ScanOnStart {
		go scanAll(ctx, libs, scanner, log)
	}

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           api.NewServer(ctx, libs, files, scanner, log).Handler(),
		ReadHeaderTimeout: 10 * time.Second,
		// No WriteTimeout: video responses stream for hours.
	}
	errc := make(chan error, 1)
	go func() {
		log.Info("instance listening", "addr", cfg.Addr, "db", cfg.DBPath, "libraries", len(cfg.Libraries))
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		log.Info("shutting down")
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

// syncLibraries makes sure every library in the config exists in the database.
// A library is identified by its root path, so renaming it in the config keeps
// its files and IDs. Libraries removed from the config are left in place.
func syncLibraries(ctx context.Context, repo library.Repository, want []config.Library, log *slog.Logger) error {
	now := time.Now().UTC()
	for _, c := range want {
		l, err := repo.GetByRoot(ctx, c.Root)
		switch {
		case errors.Is(err, library.ErrNotFound):
			l = &library.Library{Name: c.Name, Kind: c.Kind, Root: c.Root}
			l.ID, l.CreatedAt = shared.NewUUID(), now
			log.Info("library added", "name", c.Name, "kind", c.Kind, "root", c.Root)
		case err != nil:
			return err
		case l.Name == c.Name && l.Kind == c.Kind:
			continue
		default:
			l.Name, l.Kind = c.Name, c.Kind
		}
		l.UpdatedAt = now
		if err := repo.Save(ctx, l); err != nil {
			return err
		}
	}
	return nil
}

func scanAll(ctx context.Context, repo library.Repository, scanner *scan.Service, log *slog.Logger) {
	all, err := repo.List(ctx)
	if err != nil {
		log.Error("list libraries", "err", err)
		return
	}
	for _, l := range all {
		if _, err := scanner.ScanLibrary(ctx, l); err != nil {
			log.Error("scan failed", "library", l.Name, "err", err)
		}
	}
}
