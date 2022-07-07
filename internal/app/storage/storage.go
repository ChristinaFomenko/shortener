package storage

import (
	"fmt"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/jackc/pgx/v4"
)

//go:generate mockgen -source=storage.go -destination=mocks/mocks.go -package=mocks

type URLsMap = map[string]string
type UserURLs = map[string][]string

type Repository interface {
	AddURL(userShortURL UserURL) error
	GetURL(string) (string, error)
	GetList(string) []UserURL
	Ping() error
	AddBatchURL([]UserURL) error
	Destruct(cfg configs.AppConfig) error
	GetShortByOriginal(string) (string, error)
}

type Database struct {
	DB *pgx.Conn
}

func New(cfg configs.AppConfig) Repository {
	if cfg.DatabaseDSN != "" {
		dbStore, err := constructDatabaseStorage(cfg)
		if err != nil {
			fmt.Println("constructLocalStorage ERROR: ", err)
			return constructLocalStorage(cfg)
		}

		return dbStore
	}

	return constructLocalStorage(cfg)
}
