package history

import (
	"context"
	"nekosync/internal/platform/entity"
)

type Repository interface {
	// SaveEntry upserts on (UserID, WorkID, ChildID).
	SaveEntry(ctx context.Context, e *Entry) error
	GetEntry(ctx context.Context, userID, workID entity.UUID, childID *entity.UUID) (*Entry, error)
	ListEntries(ctx context.Context, userID entity.UUID, status Status, limit, offset int) ([]*Entry, error)

	AddFavorite(ctx context.Context, f *Favorite) error
	RemoveFavorite(ctx context.Context, userID, workID entity.UUID) error
	GetFavorites(ctx context.Context, userID entity.UUID, limit, offset int) ([]*Favorite, error)

	CreatePlaylist(ctx context.Context, p *Playlist) error
	GetPlaylistsByUserID(ctx context.Context, userID entity.UUID) ([]*Playlist, error)
	UpdatePlaylist(ctx context.Context, p *Playlist) error
	DeletePlaylist(ctx context.Context, id entity.UUID) error

	AddToPlaylist(ctx context.Context, item *PlaylistItem) error
	RemoveFromPlaylist(ctx context.Context, playlistID, workID entity.UUID, childID *entity.UUID) error
	GetPlaylistItems(ctx context.Context, playlistID entity.UUID) ([]*PlaylistItem, error)
}
