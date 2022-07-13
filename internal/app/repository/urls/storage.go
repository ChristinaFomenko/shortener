package urls

import (
	"context"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/file"
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/memory"
)

type Repo interface {
	Add(urlID, userID, url string) error
	Get(urlID string) (string, error)
	FetchURls(userID string) ([]models.UserURL, error)
	Ping(ctx context.Context) error
}

func NewStorage(filePath string) (Repo, error) {
	if filePath != "" {
		r, err := file.NewRepo(filePath)
		if err != nil {
			return nil, err
		}

		return r, nil
	}

	return memory.NewRepo(), nil
}
