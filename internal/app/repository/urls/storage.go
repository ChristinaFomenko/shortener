package urls

import (
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/file"
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/memory"
	"github.com/ChristinaFomenko/shortener/internal/models"
)

type Repo interface {
	Add(urlID, userID, url string) error
	Get(urlID string) (string, error)
	GetList(userID string) ([]models.UserURL, error)
	Ping() error
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
