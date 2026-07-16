package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	MALClientID     string `json:"mal_client_id"`
	MALClientSecret string `json:"mal_client_secret"`
	RedirectPort    int    `json:"redirect_port"`
	DBPath          string `json:"db_path"`
	CachePath       string `json:"cache_path"`
	LogLevel        string `json:"log_level"`
}

var defaultConfig = Config{
	MALClientID:     "",
	MALClientSecret: "",
	RedirectPort:    43829,
	DBPath:          "",
	CachePath:       "",
	LogLevel:        "info",
}

func Load(path string) (*Config, error) {
	cfg := defaultConfig

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(home, ".anidesk")
	if path == "" {
		path = filepath.Join(appDir, "config.json")
	}

	cfg.DBPath = filepath.Join(appDir, "anidesk.db")
	cfg.CachePath = filepath.Join(appDir, "cache")

	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	}

	for _, dir := range []string{appDir, cfg.CachePath} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func (c *Config) Save(path string) error {
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".anidesk", "config.json")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
