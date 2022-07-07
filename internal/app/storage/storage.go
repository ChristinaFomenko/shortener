package storage

import (
	"fmt"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/jackc/pgx/v4"
)

type URLsMap = map[string]string
type UserURLs = map[string][]string

type Repository interface {
	AddURL(userShortURL UserURL) error
	GetURL(string) (string, error)
	GetList(string) []UserURL
	Ping() error
	AddBatchURL([]UserURL) error
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
