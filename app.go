package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"

	"anidesk/backend/api"
	"anidesk/backend/auth"
	"anidesk/backend/cache"
	"anidesk/backend/config"
	"anidesk/backend/database"
	"anidesk/backend/models"
	"anidesk/backend/repository"
	"anidesk/backend/services"
	"anidesk/backend/sync"
	"anidesk/backend/workers"
)

type App struct {
	ctx            context.Context
	cfg            *config.Config
	db             *sqlx.DB
	auth           *auth.Authenticator
	apiClient      *api.Client
	cacheMgr       *cache.Manager
	syncEngine     *sync.Engine

	authSvc        *services.AuthService
	animeSvc       *services.AnimeService
	librarySvc     *services.LibraryService
	settingsSvc    *services.SettingsService

	imageWorker    *workers.Worker
	syncWorker     *workers.Worker
	cacheWorker    *workers.Worker
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	cfg, err := config.Load("")
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	a.cfg = cfg

	db, err := database.New(cfg.DBPath)
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
	a.db = db

	userRepo := repository.NewUserRepo(db)
	animeRepo := repository.NewAnimeRepo(db)
	libraryRepo := repository.NewLibraryRepo(db)
	syncRepo := repository.NewSyncQueueRepo(db)
	settingsRepo := repository.NewSettingsRepo(db)
	imageRepo := repository.NewImageRepo(db)

	a.auth = auth.New(cfg)
	a.cacheMgr = cache.New(cfg.CachePath+"/images", imageRepo)

	a.authSvc = services.NewAuthService(a.auth, userRepo)
	a.animeSvc = services.NewAnimeService(animeRepo, a.cacheMgr)
	a.librarySvc = services.NewLibraryService(libraryRepo, syncRepo)
	a.settingsSvc = services.NewSettingsService(settingsRepo)

	a.syncEngine = sync.NewEngine(syncRepo, userRepo, a.authSvc)

	if activeUser, err := userRepo.GetActiveUser(); err == nil {
		a.apiClient = api.New(activeUser.AccessToken)
	}

	a.startWorkers(db)

	if activeUser, err := userRepo.GetActiveUser(); err == nil {
		go func() {
			if err := a.authSvc.EnsureValidToken(activeUser); err == nil {
				a.apiClient = api.New(activeUser.AccessToken)
			}
		}()
	}
}

func (a *App) startWorkers(db *sqlx.DB) {
	a.imageWorker = workers.New("image-downloader", 1*time.Hour, func() error {
		_ = a.cacheMgr.Cleanup(24 * time.Hour * 7)
		return nil
	})

	a.syncWorker = workers.New("sync", 5*time.Minute, func() error {
		return a.syncEngine.Sync()
	})

	a.cacheWorker = workers.New("cache-cleaner", 24*time.Hour, func() error {
		return a.cacheMgr.Cleanup(24 * time.Hour * 14)
	})

	a.imageWorker.Start()
	a.syncWorker.Start()
	a.cacheWorker.Start()
}

func (a *App) shutdown(ctx context.Context) {
	a.syncEngine.Stop()
	if a.imageWorker != nil {
		a.imageWorker.Stop()
	}
	if a.syncWorker != nil {
		a.syncWorker.Stop()
	}
	if a.cacheWorker != nil {
		a.cacheWorker.Stop()
	}
	if a.db != nil {
		a.db.Close()
	}
}

func (a *App) setMALClient(token string) {
	a.apiClient = api.New(token)
}

func (a *App) ensureAPI() *api.Client {
	if a.apiClient == nil {
		if user, err := repository.NewUserRepo(a.db).GetActiveUser(); err == nil {
			a.apiClient = api.New(user.AccessToken)
		}
	}
	return a.apiClient
}

func (a *App) ensureAuth() error {
	user, err := repository.NewUserRepo(a.db).GetActiveUser()
	if err != nil {
		return err
	}
	return a.authSvc.EnsureValidToken(user)
}

func (a *App) GetSeasonalAnime(year int, season string) string {
	client := a.ensureAPI()
	if client == nil {
		return `[]`
	}
	result, err := client.GetSeasonalAnime(year, season, 20, 0)
	if err != nil {
		return `[]`
	}
	data, _ := json.Marshal(result)
	return string(data)
}

func (a *App) GetAnimeRanking(rankingType string) string {
	client := a.ensureAPI()
	if client == nil {
		return `[]`
	}
	result, err := client.GetRanking(rankingType, 20, 0)
	if err != nil {
		return `[]`
	}
	data, _ := json.Marshal(result)
	return string(data)
}

func (a *App) SearchAnime(query string, limit int) string {
	client := a.ensureAPI()
	if client == nil {
		return `[]`
	}
	result, err := client.SearchAnime(query, limit, 0)
	if err != nil {
		return `[]`
	}
	data, _ := json.Marshal(result)
	return string(data)
}

func (a *App) SearchManga(query string, limit int) string {
	client := a.ensureAPI()
	if client == nil {
		return `[]`
	}
	result, err := client.SearchManga(query, limit, 0)
	if err != nil {
		return `[]`
	}
	data, _ := json.Marshal(result)
	return string(data)
}

func (a *App) GetAnimeDetails(id int64) string {
	client := a.ensureAPI()
	if client == nil {
		return `{"error":"not_authenticated"}`
	}
	result, err := client.GetAnime(id, "")
	if err != nil {
		return `{"error":"` + err.Error() + `"}`
	}
	data, _ := json.Marshal(result)
	return string(data)
}

func (a *App) Login() string {
	user, err := a.authSvc.Login(a.ctx)
	if err != nil {
		return `{"error":"` + err.Error() + `"}`
	}
	a.setMALClient(user.AccessToken)
	a.syncEngine.Start(5 * time.Minute)
	data, _ := json.Marshal(user)
	return string(data)
}

func (a *App) Logout() string {
	a.authSvc.Logout()
	repo := repository.NewUserRepo(a.db)
	if user, err := repo.GetActiveUser(); err == nil {
		repo.Delete(user.ID)
	}
	a.apiClient = nil
	return `{"ok":true}`
}

func (a *App) GetActiveUser() string {
	user, err := a.authSvc.GetActiveUser()
	if err != nil {
		return `{"error":"not_logged_in"}`
	}
	data, _ := json.Marshal(user)
	return string(data)
}

func (a *App) GetLibrary() string {
	user, err := repository.NewUserRepo(a.db).GetActiveUser()
	if err != nil {
		return `{"error":"not_logged_in"}`
	}
	entries, err := a.librarySvc.GetLibrary(user.ID)
	if err != nil {
		return `{"error":"` + err.Error() + `"}`
	}
	data, _ := json.Marshal(entries)
	return string(data)
}

func (a *App) GetLibraryByStatus(status string) string {
	user, err := repository.NewUserRepo(a.db).GetActiveUser()
	if err != nil {
		return `{"error":"not_logged_in"}`
	}
	entries, err := a.librarySvc.GetByStatus(user.ID, status)
	if err != nil {
		return `{"error":"` + err.Error() + `"}`
	}
	data, _ := json.Marshal(entries)
	return string(data)
}

func (a *App) UpdateAnimeStatus(animeID int64, status string, score int, episodesWatched int) string {
	user, err := repository.NewUserRepo(a.db).GetActiveUser()
	if err != nil {
		return `{"error":"not_logged_in"}`
	}
	entry := &models.LibraryEntry{
		UserID:          user.ID,
		AnimeID:         animeID,
		Status:          status,
		Score:           score,
		EpisodesWatched: episodesWatched,
	}
	if err := a.librarySvc.UpdateEntry(entry); err != nil {
		return `{"error":"` + err.Error() + `"}`
	}
	return `{"ok":true}`
}

func (a *App) RemoveAnimeFromLibrary(animeID int64) string {
	user, err := repository.NewUserRepo(a.db).GetActiveUser()
	if err != nil {
		return `{"error":"not_logged_in"}`
	}
	if err := a.librarySvc.RemoveFromLibrary(user.ID, animeID); err != nil {
		return `{"error":"` + err.Error() + `"}`
	}
	return `{"ok":true}`
}

func (a *App) GetSettings() string {
	settings, err := a.settingsSvc.GetAll()
	if err != nil {
		return `{}`
	}
	data, _ := json.Marshal(settings)
	return string(data)
}

func (a *App) SetSetting(key, value string) string {
	if err := a.settingsSvc.Set(key, value); err != nil {
		return `{"error":"` + err.Error() + `"}`
	}
	return `{"ok":true}`
}

func (a *App) SetMALClientID(clientID string) string {
	a.cfg.MALClientID = clientID
	a.cfg.Save("")
	a.auth = auth.New(a.cfg)
	return `{"ok":true}`
}
