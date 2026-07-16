package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultConfig(t *testing.T) {
	cfg, err := Load("/tmp/nonexistent/config.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected info, got %s", cfg.LogLevel)
	}
	if cfg.RedirectPort != 43829 {
		t.Errorf("expected 43829, got %d", cfg.RedirectPort)
	}
}

func TestLoadExistingConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{"mal_client_id":"test123","log_level":"debug"}`
	os.WriteFile(path, []byte(content), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.MALClientID != "test123" {
		t.Errorf("expected test123, got %s", cfg.MALClientID)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected debug, got %s", cfg.LogLevel)
	}
}

func TestSaveConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "config.json")
	cfg := &Config{
		MALClientID: "saved_id",
		LogLevel:    "error",
		DBPath:      filepath.Join(dir, "test.db"),
		CachePath:   filepath.Join(dir, "cache"),
	}

	if err := cfg.Save(path); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if loaded.MALClientID != "saved_id" {
		t.Errorf("expected saved_id, got %s", loaded.MALClientID)
	}
}
