package urls

import (
	"context"
	"fmt"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/database"
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/file"
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/memory"
)

type Repo interface {
	Add(ctx context.Context, urlID, url, userID string) error
	Get(ctx context.Context, urlID string) (string, error)
	FetchURLs(ctx context.Context, userID string) ([]models.UserURL, error)
	Ping(ctx context.Context) error
	AddBatch(ctx context.Context, urls []models.UserURL, userID string) error
	Close() error
	DeleteUserURLs(ctx context.Context, userID string, urls []string) error
}

func NewStorage(filePath string, databaseDSN string) (Repo, error) {
	switch {
	case databaseDSN != "":
		r, err := database.NewRepo(databaseDSN)
		if err != nil {
			return nil, fmt.Errorf("initialize database repo error: %w", err)
		}
		return r, nil

	case filePath != "":

		r, err := file.NewRepo(filePath)
		if err != nil {
			return nil, fmt.Errorf("initialize file repo error: %w", err)
		}
		return r, nil
	}
	return memory.NewRepo(), nil
}
