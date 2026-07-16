package services

import (
	"context"
	"fmt"
	"time"

	"anidesk/backend/api"
	"anidesk/backend/auth"
	"anidesk/backend/cache"
	"anidesk/backend/models"
	"anidesk/backend/repository"
)

var timeNow = time.Now

type AuthService struct {
	auth     *auth.Authenticator
	userRepo *repository.UserRepo
}

func NewAuthService(auth *auth.Authenticator, userRepo *repository.UserRepo) *AuthService {
	return &AuthService{auth: auth, userRepo: userRepo}
}

func (s *AuthService) Login(ctx context.Context) (*models.User, error) {
	user, err := s.auth.Login(ctx)
	if err != nil {
		return nil, err
	}

	apiClient := api.New(user.AccessToken)
	malUser, err := apiClient.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile: %w", err)
	}

	user.MALID = malUser.ID
	user.Username = malUser.Name
	user.AvatarURL = malUser.Picture

	if err := s.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Logout() error {
	s.auth.Logout()
	return nil
}

func (s *AuthService) GetActiveUser() (*models.User, error) {
	return s.userRepo.GetActiveUser()
}

func (s *AuthService) IsLoggedIn() bool {
	user, err := s.userRepo.GetActiveUser()
	return err == nil && user != nil && user.AccessToken != ""
}

func (s *AuthService) EnsureValidToken(user *models.User) error {
	if timeNow().After(user.TokenExpiry) {
		if err := s.auth.RefreshToken(context.Background(), user); err != nil {
			return err
		}
		return s.userRepo.UpdateTokens(user)
	}
	return nil
}

type AnimeService struct {
	animeRepo *repository.AnimeRepo
	cache     *cache.Manager
}

func NewAnimeService(animeRepo *repository.AnimeRepo, cache *cache.Manager) *AnimeService {
	return &AnimeService{animeRepo: animeRepo, cache: cache}
}

func (s *AnimeService) GetAnime(id int64) (*models.Anime, error) {
	return s.animeRepo.GetByMALID(id)
}

func (s *AnimeService) Search(query string, limit int) ([]models.Anime, error) {
	return s.animeRepo.Search(query, limit)
}

func (s *AnimeService) CacheAnime(anime *models.Anime) error {
	return s.animeRepo.Save(anime)
}

type LibraryService struct {
	libraryRepo *repository.LibraryRepo
	syncRepo    *repository.SyncQueueRepo
}

func NewLibraryService(libraryRepo *repository.LibraryRepo, syncRepo *repository.SyncQueueRepo) *LibraryService {
	return &LibraryService{libraryRepo: libraryRepo, syncRepo: syncRepo}
}

func (s *LibraryService) GetLibrary(userID int64) ([]models.LibraryEntry, error) {
	return s.libraryRepo.GetByUserID(userID)
}

func (s *LibraryService) GetByStatus(userID int64, status string) ([]models.LibraryEntry, error) {
	return s.libraryRepo.GetByStatus(userID, status)
}

func (s *LibraryService) UpdateEntry(entry *models.LibraryEntry) error {
	if err := s.libraryRepo.Upsert(entry); err != nil {
		return err
	}
	payload := fmt.Sprintf(`{"anime_id":%d,"status":"%s","score":%d,"episodes_watched":%d}`,
		entry.AnimeID, entry.Status, entry.Score, entry.EpisodesWatched)
	return s.syncRepo.Enqueue(&models.SyncQueueItem{
		UserID:  entry.UserID,
		Action:  "update_anime_status",
		Payload: payload,
	})
}

func (s *LibraryService) RemoveFromLibrary(userID, animeID int64) error {
	if err := s.libraryRepo.Delete(userID, animeID); err != nil {
		return err
	}
	return s.syncRepo.Enqueue(&models.SyncQueueItem{
		UserID:  userID,
		Action:  "delete_anime",
		Payload: fmt.Sprintf(`{"anime_id":%d}`, animeID),
	})
}

type SettingsService struct {
	settingsRepo *repository.SettingsRepo
}

func NewSettingsService(settingsRepo *repository.SettingsRepo) *SettingsService {
	return &SettingsService{settingsRepo: settingsRepo}
}

func (s *SettingsService) Get(key string) (string, error) {
	return s.settingsRepo.Get(key)
}

func (s *SettingsService) Set(key, value string) error {
	return s.settingsRepo.Set(key, value)
}

func (s *SettingsService) GetAll() (map[string]string, error) {
	return s.settingsRepo.GetAll()
}
