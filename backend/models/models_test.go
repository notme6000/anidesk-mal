package models

import (
	"encoding/json"
	"testing"
)

func TestMALAnimeJSON(t *testing.T) {
	data := `{
		"id": 1,
		"title": "Test Anime",
		"mean": 8.5,
		"rank": 10,
		"media_type": "tv",
		"num_episodes": 24,
		"status": "finished_airing",
		"genres": [{"id": 1, "name": "Action"}]
	}`

	var anime MALAnime
	if err := json.Unmarshal([]byte(data), &anime); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if anime.ID != 1 {
		t.Errorf("expected ID 1, got %d", anime.ID)
	}
	if anime.Title != "Test Anime" {
		t.Errorf("expected 'Test Anime', got %s", anime.Title)
	}
	if anime.Mean != 8.5 {
		t.Errorf("expected 8.5, got %f", anime.Mean)
	}
	if len(anime.Genres) != 1 || anime.Genres[0].Name != "Action" {
		t.Error("genres not parsed correctly")
	}
}

func TestAnimeSearchResultPaging(t *testing.T) {
	data := `{
		"data": [
			{"id": 1, "title": "Anime 1"},
			{"id": 2, "title": "Anime 2"}
		],
		"paging": {
			"next": "https://api.myanimelist.net/v2/anime?offset=2"
		}
	}`

	var result AnimeSearchResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(result.Data) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Data))
	}
	if result.Paging == nil {
		t.Fatal("expected paging")
	}
	if result.Paging.Next != "https://api.myanimelist.net/v2/anime?offset=2" {
		t.Errorf("unexpected next URL: %s", result.Paging.Next)
	}
}

func TestLibraryEntryDefaults(t *testing.T) {
	entry := LibraryEntry{
		UserID:  1,
		AnimeID: 123,
		Status:  "watching",
		Score:   8,
	}

	if entry.EpisodesWatched != 0 {
		t.Errorf("expected 0, got %d", entry.EpisodesWatched)
	}
	if entry.Priority != 0 {
		t.Errorf("expected 0, got %d", entry.Priority)
	}
}

func TestUserTokenSecurity(t *testing.T) {
	user := User{
		AccessToken:  "secret-access-token",
		RefreshToken: "secret-refresh-token",
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	if _, ok := result["access_token"]; ok {
		t.Error("access_token should be hidden in JSON")
	}
	if _, ok := result["refresh_token"]; ok {
		t.Error("refresh_token should be hidden in JSON")
	}
}
