package anime

import "testing"

func TestAnime_Creation(t *testing.T) {
	anime := &Anime{
		ID:          "test-anime-id",
		NoOfEpisode: 12,
		Season:      "Spring",
		Year:        "2024",
		Status:      "completed",
		Rating:      "8.5",
	}

	if anime.ID != "test-anime-id" {
		t.Errorf("Expected ID 'test-anime-id', got %s", anime.ID)
	}

	if anime.NoOfEpisode != 12 {
		t.Errorf("Expected 12 episodes, got %d", anime.NoOfEpisode)
	}

	if anime.Rating != "8.5" {
		t.Errorf("Expected rating '8.5', got %s", anime.Rating)
	}

	if anime.Season != "Spring" {
		t.Errorf("Expected season 'Spring', got %s", anime.Season)
	}

	if anime.Year != "2024" {
		t.Errorf("Expected year '2024', got %s", anime.Year)
	}
}

func TestAnime_ValidStatus(t *testing.T) {
	validStatuses := []string{"ongoing", "completed", "upcoming", "cancelled"}

	for _, status := range validStatuses {
		anime := &Anime{
			ID:     "test-id",
			Status: status,
		}

		if anime.Status != status {
			t.Errorf("Expected status %s, got %s", status, anime.Status)
		}
	}
}
