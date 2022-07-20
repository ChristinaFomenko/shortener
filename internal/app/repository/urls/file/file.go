package file

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	errs "github.com/ChristinaFomenko/shortener/pkg/errors"
	"os"
	"sync"
)

type fileRepository struct {
	store    map[string]map[string]string
	ma       sync.RWMutex
	filePath string
}

func NewRepo(filePath string) (*fileRepository, error) {
	store, err := readLines(filePath)
	if err != nil {
		return nil, fmt.Errorf("read urls from file error: %w", err)
	}

	return &fileRepository{
		store:    store,
		filePath: filePath,
	}, nil
}

func readLines(filePath string) (map[string]map[string]string, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	res := make(map[string]map[string]string)

	if ok := scanner.Scan(); !ok {
		return res, nil
	}

	res, err = unmarshal(scanner.Bytes())
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Add URL
func (r *fileRepository) Add(_ context.Context, urlID, url, userID string) error {
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

	return r.save()
}

// Get URL
func (r *fileRepository) Get(_ context.Context, urlID string) (models.UserURL, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	for _, userStore := range r.store {
		if _, ok := userStore[urlID]; ok {
			return models.UserURL{}, nil
		}
	}

	var deleted bool
	if r.store[urlID][string(rune(1))] == "true" {
		deleted = true
	}

	return models.UserURL{
		OriginalURL: r.store[urlID][string(rune(0))],
		IsDeleted:   deleted,
	}, errs.ErrURLNotFound
}

func (r *fileRepository) FetchURLs(_ context.Context, userID string) ([]models.UserURL, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

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

func (r *fileRepository) AddBatch(_ context.Context, urls []models.UserURL, userID string) error {
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

	return r.save()
}

func (r *fileRepository) Ping(_ context.Context) error {
	return nil
}

func (r *fileRepository) Close() error {
	return nil
}

func (r *fileRepository) save() error {
	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open file error: %w", err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := marshal(r.store)
	if err != nil {
		return fmt.Errorf("serialize url error: %w", err)
	}

	_, err = file.WriteString(string(data))
	if err != nil {
		return fmt.Errorf("write url to file error: %w", err)
	}

	return nil
}

func marshal(store map[string]map[string]string) ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(store)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func unmarshal(data []byte) (map[string]map[string]string, error) {
	store := map[string]map[string]string{}

	buff := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buff)

	err := decoder.Decode(&store)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (r *fileRepository) urlExist(url string) (string, bool) {
	for _, userStore := range r.store {
		for urlID, originalURL := range userStore {
			if url == originalURL {
				return urlID, true
			}
		}
	}

	return "", false
}

func (r *fileRepository) DeleteUserURLs(ctx context.Context, userID string, urls []string) error {
	return nil
}
