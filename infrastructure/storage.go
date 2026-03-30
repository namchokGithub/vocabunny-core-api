package infrastructure

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/namchokGithub/vocabunny-core-api/configs"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
)

type LocalFileStorage struct {
	basePath string
	baseURL  string
}

func NewFileStorage(cfg configs.StorageConfig) (port.FileStorage, error) {
	if err := os.MkdirAll(cfg.BasePath, 0o755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}

	return &LocalFileStorage{
		basePath: cfg.BasePath,
		baseURL:  cfg.BaseURL,
	}, nil
}

func (s *LocalFileStorage) Save(_ context.Context, path string, payload []byte) (string, error) {
	fullPath := filepath.Join(s.basePath, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", fmt.Errorf("create storage path: %w", err)
	}

	if err := os.WriteFile(fullPath, payload, 0o644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return filepath.ToSlash(filepath.Join(s.baseURL, path)), nil
}

func (s *LocalFileStorage) Delete(_ context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}
