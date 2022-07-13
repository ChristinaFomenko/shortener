package memory

import (
	"context"
	"errors"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	"sync"
)

var ErrURLNotFound = errors.New("url not found")

type repository struct {
	store map[string]map[string]string
	ma    sync.RWMutex
}

func NewRepo() *repository {
	return &repository{
		store: map[string]map[string]string{},
	}
}

// Add URL
func (r *repository) Add(urlID, userID, url string) error {
	r.ma.Lock()
	defer r.ma.Unlock()

	userStore, ok := r.store[userID]
	if !ok {
		userStore = map[string]string{}
	}

	userStore[urlID] = url
	r.store[userID] = userStore

	return nil
}

// Get URL
func (r *repository) Get(urlID string) (string, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	for _, userStore := range r.store {
		if url, ok := userStore[urlID]; ok {
			return url, nil
		}
	}

	return "", ErrURLNotFound
}

func (r *repository) FetchURls(userID string) ([]models.UserURL, error) {
	r.ma.Lock()
	defer r.ma.Unlock()

	urls := make([]models.UserURL, 0)

	userStore, ok := r.store[userID]
	if !ok {
		return urls, nil
	}

	for shortURL, originalURL := range userStore {
		urls = append(urls, models.UserURL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		})
	}

	return urls, nil
}

func (r *repository) Ping(_ context.Context) error {
	return nil
}
