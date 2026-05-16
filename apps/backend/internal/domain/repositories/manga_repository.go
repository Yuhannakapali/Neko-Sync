package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

type MangaRepository interface {
	CreateManga(ctx context.Context, manga *entities.Manga) error
	GetMangaByContentID(ctx context.Context, contentID entities.UUID) (*entities.Manga, error)
	UpdateManga(ctx context.Context, manga *entities.Manga) error

	CreateChapter(ctx context.Context, chapter *entities.Chapter) error
	GetChaptersByContentID(ctx context.Context, contentID entities.UUID) ([]*entities.Chapter, error)
	GetChapterByNumber(ctx context.Context, contentID entities.UUID, chapterNumber int) (*entities.Chapter, error)
	UpdateChapter(ctx context.Context, chapter *entities.Chapter) error
	DeleteChapter(ctx context.Context, id entities.UUID) error
}
