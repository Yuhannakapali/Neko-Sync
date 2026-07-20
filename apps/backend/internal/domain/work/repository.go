package work

import (
	"context"

	"nekosync/internal/domain/shared"
)

// Repository is the persistence contract for the canonical metadata catalog.
type Repository interface {
	Create(ctx context.Context, w *Work) error
	GetByID(ctx context.Context, id shared.UUID) (*Work, error)
	// GetByProviderID resolves a Work from an external catalog ID. This is the path an
	// Instance uses to map a scanned file to its canonical Work.
	GetByProviderID(ctx context.Context, provider, id string) (*Work, error)
	Update(ctx context.Context, w *Work) error
	Delete(ctx context.Context, id shared.UUID) error
	List(ctx context.Context, kind Kind, limit, offset int) ([]*Work, error)
	Search(ctx context.Context, query string, kind Kind, limit, offset int) ([]*Work, error)

	CreateChild(ctx context.Context, c *WorkChild) error
	GetChild(ctx context.Context, id shared.UUID) (*WorkChild, error)
	ListChildren(ctx context.Context, workID shared.UUID) ([]*WorkChild, error)
	// GetChildByOrdinal looks up a child by its position within a parent. parentID is
	// nil for top-level children; ordinals are unique only within (workID, parentID, kind).
	GetChildByOrdinal(ctx context.Context, workID shared.UUID, parentID *shared.UUID, kind ChildKind, ordinal float64) (*WorkChild, error)
	UpdateChild(ctx context.Context, c *WorkChild) error
	DeleteChild(ctx context.Context, id shared.UUID) error
}
