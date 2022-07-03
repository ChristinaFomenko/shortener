package memory

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

func (r *repository) GetByUserID(userID string) (string, error) {
	r.ma.Lock()
	defer r.ma.Unlock()

	user, ok := r.store[userID]
	if !ok {
		return "", errors.New("user not found")
	}

	return user, nil
}

func (r *repository) Ping() error {
	//TODO implement me
	panic("implement me")
}
