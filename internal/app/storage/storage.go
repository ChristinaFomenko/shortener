package storage

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/jackc/pgx/v4"
	"os"
	"sync"
)

type URLsMap = map[string]string
type UserURLs = map[string][]string

type Repository interface {
	AddURL(ID string, URL string, userID string) error
	GetURL(ID string) (string, error)
	GetList(userID string) []string
	GetDBConn() *pgx.Conn
}

type Storage struct {
	URLsMap  URLsMap
	UserURLs UserURLs
	DB       *pgx.Conn
	mutex    sync.RWMutex
}

func (s *Storage) AddURL(ID string, URL string, userID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.URLsMap[ID] != "" {
		return fmt.Errorf(`ID=%s; URL already exists`, ID)
	}

	s.URLsMap[ID] = URL
	s.UserURLs[userID] = append(s.UserURLs[userID], ID)

	return nil
}

func (s *Storage) GetURL(ID string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	URL := s.URLsMap[ID]
	if URL == "" {
		return "", fmt.Errorf("URL not found")
	}

	return URL, nil
}

func (s *Storage) GetList(userID string) []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.UserURLs[userID]
}

func (s *Storage) GetDBConn() *pgx.Conn {
	return s.DB
}

func ConstructStorage(cfg configs.AppConfig) *Storage {
	s := &Storage{make(URLsMap), make(UserURLs), nil, sync.RWMutex{}}

	if cfg.DatabaseDSN != "" {
		conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		} else {
			s.DB = conn
		}
	}

	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		fmt.Printf("OpenFile error; %s", err)
		return s
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Printf("ReadAll error; %s", err)
		return s
	}

	for _, rec := range records {
		s.URLsMap[rec[0]] = rec[1]
	}

	return s
}

func DestructStorage(fileStoragePath string, s *Storage) error {
	file, err := os.OpenFile(fileStoragePath, os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("OpenFile error; %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	var records [][]string
	for user, shortID := range s.UserURLs {
		fmt.Println("user", user)
		fmt.Println("shortID", shortID)
	}
	for ID, URL := range s.URLsMap {
		records = append(records, []string{ID, URL})
	}

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("WriteAll error; %s", err)
	}

	writer.Flush()

	if s.DB != nil {
		s.DB.Close(context.Background())
	}

	return nil
}
