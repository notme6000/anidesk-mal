package models

import "time"

type User struct {
	ID           int64     `db:"id" json:"id"`
	MALID        int64     `db:"mal_id" json:"mal_id"`
	Username     string    `db:"username" json:"username"`
	AccessToken  string    `db:"access_token" json:"-"`
	RefreshToken string    `db:"refresh_token" json:"-"`
	TokenExpiry  time.Time `db:"token_expiry" json:"-"`
	AvatarURL    string    `db:"avatar_url" json:"avatar_url"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type Anime struct {
	ID             int64     `db:"id" json:"id"`
	MALID          int64     `db:"mal_id" json:"mal_id"`
	Title          string    `db:"title" json:"title"`
	TitleEnglish   string    `db:"title_english" json:"title_english"`
	TitleJapanese  string    `db:"title_japanese" json:"title_japanese"`
	Synopsis       string    `db:"synopsis" json:"synopsis"`
	Type           string    `db:"type" json:"type"`
	Episodes       int       `db:"episodes" json:"episodes"`
	Status         string    `db:"status" json:"status"`
	Score          float64   `db:"score" json:"score"`
	ScoredBy       int       `db:"scored_by" json:"scored_by"`
	Rank           int       `db:"rank" json:"rank"`
	Popularity     int       `db:"popularity" json:"popularity"`
	Season         string    `db:"season" json:"season"`
	Year           int       `db:"year" json:"year"`
	Airing         bool      `db:"airing" json:"airing"`
	PosterURL      string    `db:"poster_url" json:"poster_url"`
	BannerURL      string    `db:"banner_url" json:"banner_url"`
	TrailerURL     string    `db:"trailer_url" json:"trailer_url"`
	Genres         string    `db:"genres" json:"genres"`
	Studios        string    `db:"studios" json:"studios"`
	Rating         string    `db:"rating" json:"rating"`
	Source         string    `db:"source" json:"source"`
	CachedAt       time.Time `db:"cached_at" json:"cached_at"`
}

type Manga struct {
	ID            int64     `db:"id" json:"id"`
	MALID         int64     `db:"mal_id" json:"mal_id"`
	Title         string    `db:"title" json:"title"`
	TitleEnglish  string    `db:"title_english" json:"title_english"`
	TitleJapanese string    `db:"title_japanese" json:"title_japanese"`
	Synopsis      string    `db:"synopsis" json:"synopsis"`
	Type          string    `db:"type" json:"type"`
	Chapters      int       `db:"chapters" json:"chapters"`
	Volumes       int       `db:"volumes" json:"volumes"`
	Status        string    `db:"status" json:"status"`
	Score         float64   `db:"score" json:"score"`
	Rank          int       `db:"rank" json:"rank"`
	PosterURL     string    `db:"poster_url" json:"poster_url"`
	Genres        string    `db:"genres" json:"genres"`
	CachedAt      time.Time `db:"cached_at" json:"cached_at"`
}

type LibraryEntry struct {
	ID           int64     `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	AnimeID      int64     `db:"anime_id" json:"anime_id"`
	Status       string    `db:"status" json:"status"`
	Score        int       `db:"score" json:"score"`
	EpisodesWatched int    `db:"episodes_watched" json:"episodes_watched"`
	Notes        string    `db:"notes" json:"notes"`
	Priority     int       `db:"priority" json:"priority"`
	IsFavorite   bool      `db:"is_favorite" json:"is_favorite"`
	StartedAt    *time.Time `db:"started_at" json:"started_at"`
	CompletedAt  *time.Time `db:"completed_at" json:"completed_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	SyncedAt     *time.Time `db:"synced_at" json:"synced_at"`
}

type HistoryEntry struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	AnimeID   int64     `db:"anime_id" json:"anime_id"`
	Action    string    `db:"action" json:"action"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type SyncQueueItem struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Action    string    `db:"action" json:"action"`
	Payload   string    `db:"payload" json:"payload"`
	Status    string    `db:"status" json:"status"`
	Retries   int       `db:"retries" json:"retries"`
	Error     string    `db:"error" json:"error"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CachedImage struct {
	ID        int64     `db:"id" json:"id"`
	URL       string    `db:"url" json:"url"`
	LocalPath string    `db:"local_path" json:"local_path"`
	FileSize  int64     `db:"file_size" json:"file_size"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Setting struct {
	Key   string `db:"key" json:"key"`
	Value string `db:"value" json:"value"`
}

type MALAnime struct {
	ID         int64              `json:"id"`
	Title      string             `json:"title"`
	MainPicture *MALImage         `json:"main_picture"`
	AlternativeTitles *MALAlternateTitles `json:"alternative_titles"`
	Synopsis   string             `json:"synopsis"`
	Mean       float64            `json:"mean"`
	Rank       int                `json:"rank"`
	Popularity int                `json:"popularity"`
	MediaType  string             `json:"media_type"`
	Status     string             `json:"status"`
	NumEpisodes int               `json:"num_episodes"`
	StartSeason *MALSeason        `json:"start_season"`
	Genres     []MALGenre         `json:"genres"`
	Studios    []MALStudio        `json:"studios"`
	Rating     string             `json:"rating"`
	Pictures   []MALImage         `json:"pictures"`
	Recommendations []MALRecommendation `json:"recommendations"`
	RelatedAnime []MALRelated     `json:"related_anime"`
}

type MALImage struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type MALAlternateTitles struct {
	Synonyms []string `json:"synonyms"`
	English  string   `json:"en"`
	Japanese string   `json:"ja"`
}

type MALSeason struct {
	Year   int    `json:"year"`
	Season string `json:"season"`
}

type MALGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MALStudio struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MALRecommendation struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Image  *MALImage `json:"main_picture"`
}

type MALRelated struct {
	Node  *MALAnime `json:"node"`
	RelationType string `json:"relation_type"`
}

type MALUser struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Picture  string    `json:"picture"`
}

type MALAnimeNode struct {
	Node MALAnime `json:"node"`
}

type MALRankedAnimeNode struct {
	Node    MALAnime `json:"node"`
	Ranking *struct {
		Rank int `json:"rank"`
	} `json:"ranking"`
}

type MALListStatusNode struct {
	Node       MALAnime    `json:"node"`
	ListStatus *MALListStatus `json:"list_status"`
}

type MALListStatus struct {
	Status       string `json:"status"`
	Score        int    `json:"score"`
	EpisodesWatched int `json:"num_episodes_watched"`
	IsRewatching bool   `json:"is_rewatching"`
	StartDate    string `json:"start_date"`
	FinishDate   string `json:"finish_date"`
	Priority     int    `json:"priority"`
	Notes        string `json:"notes"`
}

type AnimeSearchResult struct {
	Data       []MALAnimeNode `json:"data"`
	Paging     *struct {
		Next string `json:"next"`
	} `json:"paging"`
}

type AnimeRankingResult struct {
	Data   []MALRankedAnimeNode `json:"data"`
	Paging *struct {
		Next string `json:"next"`
	} `json:"paging"`
}
