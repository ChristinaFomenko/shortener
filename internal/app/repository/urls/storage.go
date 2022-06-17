package urls

import (
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/file"
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/memory"
)

type Repo interface {
	Add(id, url string) error
	Get(id string) (string, error)
}

func Storage(filePath string) Repo {
	if filePath != "" {
		r, err := file.NewRepo(filePath)
		if err != nil {
			panic(err)
		}

		return r
	}

	return memory.NewRepo()
}
