package memory

import (
	"errors"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	"sync"
)

type repository struct {
	store map[string]string
	ma    sync.RWMutex
}

func NewRepo() *repository {
	return &repository{
		store: map[string]string{},
	}
}

// Add URL
func (r *repository) Add(id, url string) error {
	r.ma.Lock()
	defer r.ma.Unlock()

	r.store[id] = url
	return nil
}

// Get URL
func (r *repository) Get(id string) (string, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	url, ok := r.store[id]
	if !ok {
		return "", errors.New("url not found")
	}

	return url, nil
}

func (r *repository) GetList() ([]models.UserURL, error) {
	r.ma.Lock()
	defer r.ma.Unlock()

	urls := make([]models.UserURL, 0, len(r.store))
	for shortURL, originalURL := range r.store {
		urls = append(urls, models.UserURL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		})
	}

	return urls, nil
}

func (r *repository) Ping() error {
	return nil
}
