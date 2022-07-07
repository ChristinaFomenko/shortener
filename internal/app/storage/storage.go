package storage

import (
	"fmt"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/jackc/pgx/v4"
)

type URLsMap = map[string]string
type UserURLs = map[string][]string

type Repository interface {
	AddURL(ID string, URL string, userID string) error
	GetURL(ID string) (string, error)
	GetList(userID string) []UserURL
	Ping() error
	DestructStorage(cfg configs.AppConfig) error
}

type Database struct {
	DB *pgx.Conn
}

func ConstructStorage(cfg configs.AppConfig) Repository {
	if cfg.DatabaseDSN != "" {
		dbStore, err := ConstructDatabaseStorage(cfg)
		if err != nil {
			fmt.Println("ConstructDatabaseStorage ERROR: ", err)
			return ConstructLocalStorage(cfg)
		}

		return dbStore
	}

	return ConstructLocalStorage(cfg)
}
