package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	"os"
	"sync"
)

var ErrURLNotFound = errors.New("url not found")

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
	file, err := os.OpenFile(filePath, os.O_CREATE, 0777)
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

	err = json.Unmarshal([]byte(scanner.Text()), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Add URL
func (r *fileRepository) Add(urlID, userID, url string) error {
	r.ma.Lock()
	defer r.ma.Unlock()

	userStore, ok := r.store[userID]
	if !ok {
		userStore = map[string]string{}
	}

	userStore[urlID] = url
	r.store[userID] = userStore

	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("open file error: %w", err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := json.Marshal(r.store)
	if err != nil {
		return fmt.Errorf("serialize url error: %w", err)
	}

	_, err = file.WriteString(string(data))
	if err != nil {
		return fmt.Errorf("write url to file error: %w", err)
	}

	return nil
}

// Get URL
func (r *fileRepository) Get(urlID string) (string, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	for _, userStore := range r.store {
		if url, ok := userStore[urlID]; ok {
			return url, nil
		}
	}

	return "", ErrURLNotFound
}

func (r *fileRepository) FetchURls(userID string) ([]models.UserURL, error) {
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

func (r *fileRepository) Ping() error {
	return nil
}
