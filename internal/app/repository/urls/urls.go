package urls

import (
	"errors"
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
func (r *repository) Add(id, url string) {
	r.ma.Lock()
	defer r.ma.Unlock()

	r.store[id] = url
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
