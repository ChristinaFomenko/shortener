package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ChristinaFomenko/shortener/internal/models"
	"os"
	"sync"
)

type fileRepository struct {
	store    map[string]string
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

func readLines(filePath string) (map[string]string, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	res := make(map[string]string)

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
func (r *fileRepository) Add(models models.UserURL) error {
	r.ma.Lock()
	defer r.ma.Unlock()

	r.store[models.UserID] = models.ShortURL

	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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
func (r *fileRepository) Get(id string) (string, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	url, ok := r.store[id]
	if !ok {
		return "", errors.New("url not found")
	}

	return url, nil
}

func (r *fileRepository) GetList() ([]models.UserURL, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	urls := make([]models.UserURL, 0, len(r.store))
	for shortURL, originalURL := range r.store {
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
