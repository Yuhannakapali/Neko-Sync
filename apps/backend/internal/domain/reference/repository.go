package reference

import (
	"context"

	"nekosync/internal/domain/shared"
)

// Repository is the persistence contract for the ContentReference registry.
type Repository interface {
	Create(ctx context.Context, r *ContentReference) error
	// Upsert creates or updates a reference by its natural identity
	// (UserID, WorkID, ChildID, Source, Locator). A user's library changes on every
	// re-scan, so references must be reconcilable rather than only append-only.
	Upsert(ctx context.Context, r *ContentReference) error
	Update(ctx context.Context, r *ContentReference) error
	Delete(ctx context.Context, id shared.UUID) error
	// Resolve returns every reference a user has for a Work (optionally narrowed to a
	// specific child). Callers typically pass the result through Rank. A nil childID
	// matches references whose ChildID is also nil (Work-level playables).
	Resolve(ctx context.Context, userID, workID shared.UUID, childID *shared.UUID) ([]*ContentReference, error)
	ListByUser(ctx context.Context, userID shared.UUID, limit, offset int) ([]*ContentReference, error)
}
