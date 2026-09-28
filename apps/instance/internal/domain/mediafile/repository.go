package mediafile

import (
	"context"

	"nekosync-instance/internal/domain/shared"
)

// Repository persists MediaFiles.
type Repository interface {
	Get(ctx context.Context, id shared.UUID) (*MediaFile, error)
	GetByPath(ctx context.Context, path string) (*MediaFile, error)
	ListByLibrary(ctx context.Context, libraryID shared.UUID) ([]*MediaFile, error)
	// Save inserts the file, or replaces it if a file with the same ID exists.
	Save(ctx context.Context, f *MediaFile) error
	Delete(ctx context.Context, id shared.UUID) error
}
