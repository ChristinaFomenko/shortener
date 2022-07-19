package memory

import (
	"context"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	errs "github.com/ChristinaFomenko/shortener/pkg/errors"
	"sync"
)

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
func (r *repository) Add(_ context.Context, urlID, url, userID string) error {
	r.ma.Lock()
	defer r.ma.Unlock()

	if doubleURLID, exists := r.urlExist(url); exists {
		return errs.NewNotUniqueURLErr(doubleURLID, url, nil)
	}

	userStore, ok := r.store[userID]
	if !ok {
		userStore = map[string]string{}
	}

	userStore[urlID] = url
	r.store[userID] = userStore

	return nil
}

// Get URL
func (r *repository) Get(_ context.Context, urlID string) (string, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	for _, userStore := range r.store {
		if url, ok := userStore[urlID]; ok {
			return url, nil
		}
	}

	return "", errs.ErrURLNotFound
}

func (r *repository) FetchURLs(_ context.Context, userID string) ([]models.UserURL, error) {
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

func (r *repository) Close() error {
	return nil
}

func (r *repository) AddBatch(_ context.Context, urls []models.UserURL, userID string) error {
	r.ma.Lock()
	defer r.ma.Unlock()

	userStore, ok := r.store[userID]
	if !ok {
		userStore = map[string]string{}
	}

	for idx := range urls {
		userStore[urls[idx].ShortURL] = urls[idx].OriginalURL
	}

	r.store[userID] = userStore

	return nil
}

func (r *repository) urlExist(url string) (string, bool) {
	for _, userStore := range r.store {
		for urlID, originalURL := range userStore {
			if url == originalURL {
				return urlID, true
			}
		}
	}

	return "", false
}
