package content

import (
	"testing"
	"time"

	"nekosync/internal/domain"
)

func TestSeries_Creation(t *testing.T) {
	description := "Test series description"
	coverURL := "https://example.com/cover.jpg"

	series := &Series{
		BaseEntity: domain.BaseEntity{
			ID:        domain.UUID("test-series-id"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Title:       "Test Series",
		Description: &description,
		CoverURL:    &coverURL,
		Type:        domain.ContentTypeAnime,
	}

	if series.Title != "Test Series" {
		t.Errorf("Expected title 'Test Series', got %s", series.Title)
	}

	if *series.Description != description {
		t.Errorf("Expected description '%s', got %s", description, *series.Description)
	}

	if series.Type != domain.ContentTypeAnime {
		t.Errorf("Expected type anime, got %s", series.Type)
	}
}

func TestSeries_OptionalFields(t *testing.T) {
	series := &Series{
		BaseEntity: domain.BaseEntity{
			ID: domain.UUID("test-id"),
		},
		Title: "Test Series",
		Type:  domain.ContentTypeManga,
	}

	if series.Description != nil {
		t.Error("Expected nil description")
	}

	if series.CoverURL != nil {
		t.Error("Expected nil cover URL")
	}
}
