package sync

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"anidesk/backend/api"
	"anidesk/backend/models"
	"anidesk/backend/repository"
	"anidesk/backend/services"
)

type Engine struct {
	syncRepo  *repository.SyncQueueRepo
	userRepo  *repository.UserRepo
	authSvc   *services.AuthService
	apiClient *api.Client
	running   bool
	stopCh    chan struct{}
}

func NewEngine(
	syncRepo *repository.SyncQueueRepo,
	userRepo *repository.UserRepo,
	authSvc *services.AuthService,
) *Engine {
	return &Engine{
		syncRepo: syncRepo,
		userRepo: userRepo,
		authSvc:  authSvc,
		stopCh:   make(chan struct{}),
	}
}

func (e *Engine) Start(interval time.Duration) {
	if e.running {
		return
	}
	e.running = true

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := e.Sync(); err != nil {
					log.Printf("sync error: %v", err)
				}
			case <-e.stopCh:
				e.running = false
				return
			}
		}
	}()

	log.Println("sync engine started")
}

func (e *Engine) Stop() {
	if e.running {
		close(e.stopCh)
	}
}

func (e *Engine) Sync() error {
	user, err := e.userRepo.GetActiveUser()
	if err != nil {
		return fmt.Errorf("no active user: %w", err)
	}

	if err := e.authSvc.EnsureValidToken(user); err != nil {
		return fmt.Errorf("token refresh failed: %w", err)
	}

	e.apiClient = api.New(user.AccessToken)

	items, err := e.syncRepo.GetPending()
	if err != nil {
		return fmt.Errorf("failed to get pending items: %w", err)
	}

	for _, item := range items {
		if err := e.processItem(item); err != nil {
			e.syncRepo.MarkFailed(item.ID, err.Error())
			log.Printf("sync item %d failed: %v", item.ID, err)
			continue
		}
		e.syncRepo.MarkDone(item.ID)
	}

	e.syncRepo.Cleanup()
	return nil
}

func (e *Engine) processItem(item models.SyncQueueItem) error {
	switch item.Action {
	case "update_anime_status":
		var payload struct {
			AnimeID         int64  `json:"anime_id"`
			Status          string `json:"status"`
			Score           int    `json:"score"`
			EpisodesWatched int    `json:"episodes_watched"`
		}
		if err := json.Unmarshal([]byte(item.Payload), &payload); err != nil {
			return err
		}
		return e.apiClient.UpdateAnimeStatus(payload.AnimeID, payload.Status, payload.Score, payload.EpisodesWatched)

	case "delete_anime":
		var payload struct {
			AnimeID int64 `json:"anime_id"`
		}
		if err := json.Unmarshal([]byte(item.Payload), &payload); err != nil {
			return err
		}
		return e.apiClient.DeleteAnimeFromList(payload.AnimeID)

	default:
		return fmt.Errorf("unknown action: %s", item.Action)
	}
}
