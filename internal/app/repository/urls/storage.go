package urls

import (
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/file"
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/memory"
)

type Repo interface {
	Add(id, url string) error
	Get(id string) (string, error)
	GetByUserID(userID string) (string, error)
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
