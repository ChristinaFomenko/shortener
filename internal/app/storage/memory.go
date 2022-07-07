package storage

import (
	"encoding/csv"
	"fmt"
	"github.com/ChristinaFomenko/shortener/configs"
	"os"
	"sync"
)

type LocalStorage struct {
	URLsMap  URLsMap
	UserURLs UserURLs
	mutex    sync.RWMutex
}

func (err *URLDuplicateError) Error() string {
	return fmt.Sprintf("URL %s - already exists.", err.URL)
}

func (ls *LocalStorage) AddURL(userShortURL UserURL) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	if ls.URLsMap[userShortURL.ID] != "" {
		return fmt.Errorf(`ID=%s; URL already exists`, userShortURL.ID)
	}

	existID, _ := ls.GetShortByOriginal(userShortURL.OriginalURL)
	if existID != "" {
		return &URLDuplicateError{URL: userShortURL.OriginalURL}
	}

	ls.URLsMap[userShortURL.ID] = userShortURL.OriginalURL
	ls.UserURLs[userShortURL.UserID] = append(ls.UserURLs[userShortURL.UserID], userShortURL.ID)

	return nil
}

func (ls *LocalStorage) GetURL(ID string) (string, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	URL := ls.URLsMap[ID]
	if URL == "" {
		return "", fmt.Errorf("URL not found")
	}

	return URL, nil
}

func (ls *LocalStorage) GetList(userID string) []UserURL {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	userURLs := ls.UserURLs[userID]

	var URLs []UserURL
	for _, ID := range userURLs {
		URL, _ := ls.GetURL(ID)
		URLs = append(URLs, UserURL{ID, URL, userID})
	}

	return URLs
}

func ConstructLocalStorage(cfg configs.AppConfig) Repository {
	ls := &LocalStorage{make(URLsMap), make(UserURLs), sync.RWMutex{}}

	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		fmt.Printf("OpenFile error; %s", err)
		return ls
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Printf("ReadAll error; %s", err)
		return ls
	}

	for _, rec := range records {
		ls.URLsMap[rec[0]] = rec[1]
	}

	return ls
}

func (ls *LocalStorage) DestructStorage(cfg configs.AppConfig) error {
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("OpenFile error; %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	var records [][]string
	for ID, URL := range ls.URLsMap {
		records = append(records, []string{ID, URL})
	}

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("WriteAll error; %s", err)
	}

	writer.Flush()

	return nil
}

func (ls *LocalStorage) Ping() error {
	return nil
}

func (ls *LocalStorage) AddBatchURL(urls []UserURL) error {
	for _, url := range urls {
		err := ls.AddURL(url)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ls *LocalStorage) GetShortByOriginal(originalURL string) (string, error) {
	for ID, URL := range ls.URLsMap {
		if URL == originalURL {
			return ID, nil
		}
	}

	return "", fmt.Errorf("URL not found")
}
