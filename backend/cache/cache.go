package cache

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"anidesk/backend/models"
	"anidesk/backend/repository"
)

type Manager struct {
	basePath string
	imageRepo *repository.ImageRepo
	client   *http.Client
}

func New(basePath string, imageRepo *repository.ImageRepo) *Manager {
	os.MkdirAll(basePath, 0755)
	return &Manager{
		basePath:  basePath,
		imageRepo: imageRepo,
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (m *Manager) GetImage(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("empty url")
	}

	cached, err := m.imageRepo.GetByURL(url)
	if err == nil && cached.LocalPath != "" {
		if _, err := os.Stat(cached.LocalPath); err == nil {
			return cached.LocalPath, nil
		}
	}

	localPath, err := m.downloadImage(url)
	if err != nil {
		return "", err
	}

	m.imageRepo.Save(&models.CachedImage{
		URL:       url,
		LocalPath: localPath,
	})

	return localPath, nil
}

func (m *Manager) downloadImage(url string) (string, error) {
	resp, err := m.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(url))
	filename := fmt.Sprintf("%x", hash[:8])
	ext := ".jpg"

	localPath := filepath.Join(m.basePath, filename+ext)
	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return "", err
	}

	return localPath, nil
}

func (m *Manager) Clear() error {
	entries, err := os.ReadDir(m.basePath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		os.Remove(filepath.Join(m.basePath, entry.Name()))
	}
	return nil
}

func (m *Manager) Cleanup(maxAge time.Duration) error {
	entries, err := os.ReadDir(m.basePath)
	if err != nil {
		return err
	}
	cutoff := time.Now().Add(-maxAge)
	for _, entry := range entries {
		info, err := entry.Info()
		if err == nil && info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(m.basePath, entry.Name()))
		}
	}
	return nil
}
