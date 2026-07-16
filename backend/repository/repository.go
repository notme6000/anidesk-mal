package repository

import (
	"time"

	"github.com/jmoiron/sqlx"

	"anidesk/backend/models"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Save(user *models.User) error {
	query := `INSERT INTO users (mal_id, username, access_token, refresh_token, token_expiry, avatar_url, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT(mal_id) DO UPDATE SET
			access_token = excluded.access_token,
			refresh_token = excluded.refresh_token,
			token_expiry = excluded.token_expiry,
			avatar_url = excluded.avatar_url,
			updated_at = datetime('now')`
	_, err := r.db.Exec(query, user.MALID, user.Username, user.AccessToken, user.RefreshToken, user.TokenExpiry, user.AvatarURL)
	return err
}

func (r *UserRepo) GetActiveUser() (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM users ORDER BY id DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdateTokens(user *models.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET access_token = ?, refresh_token = ?, token_expiry = ?, updated_at = datetime('now') WHERE id = ?",
		user.AccessToken, user.RefreshToken, user.TokenExpiry, user.ID)
	return err
}

func (r *UserRepo) Delete(userID int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", userID)
	return err
}

type AnimeRepo struct {
	db *sqlx.DB
}

func NewAnimeRepo(db *sqlx.DB) *AnimeRepo {
	return &AnimeRepo{db: db}
}

func (r *AnimeRepo) Save(anime *models.Anime) error {
	query := `INSERT INTO anime_cache (
		mal_id, title, title_english, title_japanese, synopsis, type, episodes, status,
		score, scored_by, rank, popularity, season, year, airing, poster_url, banner_url,
		trailer_url, genres, studios, rating, source, cached_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
	ON CONFLICT(mal_id) DO UPDATE SET
		title = excluded.title, title_english = excluded.title_english,
		title_japanese = excluded.title_japanese, synopsis = excluded.synopsis,
		type = excluded.type, episodes = excluded.episodes, status = excluded.status,
		score = excluded.score, scored_by = excluded.scored_by, rank = excluded.rank,
		popularity = excluded.popularity, season = excluded.season, year = excluded.year,
		airing = excluded.airing, poster_url = excluded.poster_url,
		banner_url = excluded.banner_url, trailer_url = excluded.trailer_url,
		genres = excluded.genres, studios = excluded.studios, rating = excluded.rating,
		source = excluded.source, cached_at = datetime('now')`
	_, err := r.db.Exec(query,
		anime.MALID, anime.Title, anime.TitleEnglish, anime.TitleJapanese,
		anime.Synopsis, anime.Type, anime.Episodes, anime.Status,
		anime.Score, anime.ScoredBy, anime.Rank, anime.Popularity,
		anime.Season, anime.Year, boolToInt(anime.Airing),
		anime.PosterURL, anime.BannerURL, anime.TrailerURL,
		anime.Genres, anime.Studios, anime.Rating, anime.Source)
	return err
}

func (r *AnimeRepo) GetByMALID(malID int64) (*models.Anime, error) {
	var anime models.Anime
	err := r.db.Get(&anime, "SELECT * FROM anime_cache WHERE mal_id = ?", malID)
	if err != nil {
		return nil, err
	}
	return &anime, nil
}

func (r *AnimeRepo) GetCachedIDs() ([]int64, error) {
	var ids []int64
	err := r.db.Select(&ids, "SELECT mal_id FROM anime_cache")
	return ids, err
}

func (r *AnimeRepo) Search(query string, limit int) ([]models.Anime, error) {
	var results []models.Anime
	err := r.db.Select(&results,
		"SELECT * FROM anime_cache WHERE title LIKE ? OR title_english LIKE ? LIMIT ?",
		"%"+query+"%", "%"+query+"%", limit)
	return results, err
}

func (r *AnimeRepo) DeleteOlderThan(t time.Time) error {
	_, err := r.db.Exec("DELETE FROM anime_cache WHERE cached_at < ?", t)
	return err
}

type LibraryRepo struct {
	db *sqlx.DB
}

func NewLibraryRepo(db *sqlx.DB) *LibraryRepo {
	return &LibraryRepo{db: db}
}

func (r *LibraryRepo) Upsert(entry *models.LibraryEntry) error {
	query := `INSERT INTO library (user_id, anime_id, status, score, episodes_watched, notes, priority, is_favorite, started_at, completed_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT(user_id, anime_id) DO UPDATE SET
			status = excluded.status, score = excluded.score,
			episodes_watched = excluded.episodes_watched, notes = excluded.notes,
			priority = excluded.priority, is_favorite = excluded.is_favorite,
			started_at = excluded.started_at, completed_at = excluded.completed_at,
			updated_at = datetime('now')`
	_, err := r.db.Exec(query, entry.UserID, entry.AnimeID, entry.Status,
		entry.Score, entry.EpisodesWatched, entry.Notes, entry.Priority,
		boolToInt(entry.IsFavorite), entry.StartedAt, entry.CompletedAt)
	return err
}

func (r *LibraryRepo) GetByUserID(userID int64) ([]models.LibraryEntry, error) {
	var entries []models.LibraryEntry
	err := r.db.Select(&entries, "SELECT * FROM library WHERE user_id = ? ORDER BY updated_at DESC", userID)
	return entries, err
}

func (r *LibraryRepo) GetByStatus(userID int64, status string) ([]models.LibraryEntry, error) {
	var entries []models.LibraryEntry
	err := r.db.Select(&entries, "SELECT * FROM library WHERE user_id = ? AND status = ? ORDER BY updated_at DESC", userID, status)
	return entries, err
}

func (r *LibraryRepo) GetByUserAndAnime(userID int64, animeID int64) (*models.LibraryEntry, error) {
	var entry models.LibraryEntry
	err := r.db.Get(&entry, "SELECT * FROM library WHERE user_id = ? AND anime_id = ?", userID, animeID)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *LibraryRepo) Delete(userID, animeID int64) error {
	_, err := r.db.Exec("DELETE FROM library WHERE user_id = ? AND anime_id = ?", userID, animeID)
	return err
}

type SyncQueueRepo struct {
	db *sqlx.DB
}

func NewSyncQueueRepo(db *sqlx.DB) *SyncQueueRepo {
	return &SyncQueueRepo{db: db}
}

func (r *SyncQueueRepo) Enqueue(item *models.SyncQueueItem) error {
	_, err := r.db.Exec(
		"INSERT INTO sync_queue (user_id, action, payload, status) VALUES (?, ?, ?, 'pending')",
		item.UserID, item.Action, item.Payload)
	return err
}

func (r *SyncQueueRepo) GetPending() ([]models.SyncQueueItem, error) {
	var items []models.SyncQueueItem
	err := r.db.Select(&items, "SELECT * FROM sync_queue WHERE status = 'pending' ORDER BY created_at ASC LIMIT 50")
	return items, err
}

func (r *SyncQueueRepo) MarkDone(id int64) error {
	_, err := r.db.Exec("UPDATE sync_queue SET status = 'done', updated_at = datetime('now') WHERE id = ?", id)
	return err
}

func (r *SyncQueueRepo) MarkFailed(id int64, errMsg string) error {
	_, err := r.db.Exec(
		"UPDATE sync_queue SET status = 'failed', error = ?, retries = retries + 1, updated_at = datetime('now') WHERE id = ?",
		errMsg, id)
	return err
}

func (r *SyncQueueRepo) Cleanup() error {
	_, err := r.db.Exec("DELETE FROM sync_queue WHERE status IN ('done', 'failed') AND created_at < datetime('now', '-7 days')")
	return err
}

type HistoryRepo struct {
	db *sqlx.DB
}

func NewHistoryRepo(db *sqlx.DB) *HistoryRepo {
	return &HistoryRepo{db: db}
}

func (r *HistoryRepo) Add(entry *models.HistoryEntry) error {
	_, err := r.db.Exec("INSERT INTO history (user_id, anime_id, action) VALUES (?, ?, ?)",
		entry.UserID, entry.AnimeID, entry.Action)
	return err
}

func (r *HistoryRepo) GetByUserID(userID int64, limit int) ([]models.HistoryEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	var entries []models.HistoryEntry
	err := r.db.Select(&entries, "SELECT * FROM history WHERE user_id = ? ORDER BY created_at DESC LIMIT ?", userID, limit)
	return entries, err
}

type SettingsRepo struct {
	db *sqlx.DB
}

func NewSettingsRepo(db *sqlx.DB) *SettingsRepo {
	return &SettingsRepo{db: db}
}

func (r *SettingsRepo) Get(key string) (string, error) {
	var s models.Setting
	err := r.db.Get(&s, "SELECT * FROM settings WHERE key = ?", key)
	if err != nil {
		return "", err
	}
	return s.Value, nil
}

func (r *SettingsRepo) Set(key, value string) error {
	_, err := r.db.Exec("INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value", key, value)
	return err
}

func (r *SettingsRepo) GetAll() (map[string]string, error) {
	var settings []models.Setting
	err := r.db.Select(&settings, "SELECT * FROM settings")
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(settings))
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

type ImageRepo struct {
	db *sqlx.DB
}

func NewImageRepo(db *sqlx.DB) *ImageRepo {
	return &ImageRepo{db: db}
}

func (r *ImageRepo) Save(img *models.CachedImage) error {
	_, err := r.db.Exec(
		"INSERT INTO images (url, local_path, file_size) VALUES (?, ?, ?) ON CONFLICT(url) DO UPDATE SET local_path = excluded.local_path, file_size = excluded.file_size",
		img.URL, img.LocalPath, img.FileSize)
	return err
}

func (r *ImageRepo) GetByURL(url string) (*models.CachedImage, error) {
	var img models.CachedImage
	err := r.db.Get(&img, "SELECT * FROM images WHERE url = ?", url)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
