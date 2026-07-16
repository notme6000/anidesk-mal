package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func New(path string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", path+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mal_id INTEGER UNIQUE,
		username TEXT NOT NULL,
		access_token TEXT,
		refresh_token TEXT,
		token_expiry TEXT,
		avatar_url TEXT,
		created_at TEXT DEFAULT (datetime('now')),
		updated_at TEXT DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS anime_cache (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mal_id INTEGER UNIQUE,
		title TEXT,
		title_english TEXT,
		title_japanese TEXT,
		synopsis TEXT,
		type TEXT,
		episodes INTEGER DEFAULT 0,
		status TEXT,
		score REAL DEFAULT 0,
		scored_by INTEGER DEFAULT 0,
		rank INTEGER DEFAULT 0,
		popularity INTEGER DEFAULT 0,
		season TEXT,
		year INTEGER DEFAULT 0,
		airing INTEGER DEFAULT 0,
		poster_url TEXT,
		banner_url TEXT,
		trailer_url TEXT,
		genres TEXT,
		studios TEXT,
		rating TEXT,
		source TEXT,
		cached_at TEXT DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS manga_cache (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mal_id INTEGER UNIQUE,
		title TEXT,
		title_english TEXT,
		title_japanese TEXT,
		synopsis TEXT,
		type TEXT,
		chapters INTEGER DEFAULT 0,
		volumes INTEGER DEFAULT 0,
		status TEXT,
		score REAL DEFAULT 0,
		rank INTEGER DEFAULT 0,
		poster_url TEXT,
		genres TEXT,
		cached_at TEXT DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS library (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		anime_id INTEGER NOT NULL,
		status TEXT DEFAULT 'plan_to_watch',
		score INTEGER DEFAULT 0,
		episodes_watched INTEGER DEFAULT 0,
		notes TEXT DEFAULT '',
		priority INTEGER DEFAULT 0,
		is_favorite INTEGER DEFAULT 0,
		started_at TEXT,
		completed_at TEXT,
		updated_at TEXT DEFAULT (datetime('now')),
		synced_at TEXT,
		UNIQUE(user_id, anime_id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (anime_id) REFERENCES anime_cache(mal_id)
	);

	CREATE TABLE IF NOT EXISTS favorites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		anime_id INTEGER NOT NULL,
		created_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (user_id) REFERENCES users(id),
		UNIQUE(user_id, anime_id)
	);

	CREATE TABLE IF NOT EXISTS history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		anime_id INTEGER,
		action TEXT NOT NULL,
		created_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sync_queue (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		action TEXT NOT NULL,
		payload TEXT,
		status TEXT DEFAULT 'pending',
		retries INTEGER DEFAULT 0,
		error TEXT,
		created_at TEXT DEFAULT (datetime('now')),
		updated_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT UNIQUE,
		local_path TEXT,
		file_size INTEGER DEFAULT 0,
		created_at TEXT DEFAULT (datetime('now'))
	);

	CREATE INDEX IF NOT EXISTS idx_library_user ON library(user_id);
	CREATE INDEX IF NOT EXISTS idx_library_status ON library(status);
	CREATE INDEX IF NOT EXISTS idx_sync_queue_status ON sync_queue(status);
	CREATE INDEX IF NOT EXISTS idx_history_user ON history(user_id);
	CREATE INDEX IF NOT EXISTS idx_anime_cache_title ON anime_cache(title);
	`

	_, err := db.Exec(schema)
	return err
}
