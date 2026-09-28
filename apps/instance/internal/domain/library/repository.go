package library

import (
	"context"
	"errors"

	"nekosync-instance/internal/domain/shared"
)

// ErrNotFound is returned by a Repository when no Library matches.
var ErrNotFound = errors.New("library not found")

// Repository persists Libraries.
type Repository interface {
	Get(ctx context.Context, id shared.UUID) (*Library, error)
	GetByRoot(ctx context.Context, root string) (*Library, error)
	List(ctx context.Context) ([]*Library, error)
	Save(ctx context.Context, l *Library) error
}
