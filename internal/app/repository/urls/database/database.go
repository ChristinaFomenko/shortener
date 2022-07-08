package database

import (
	"database/sql"
	"github.com/ChristinaFomenko/shortener/internal/models"
	_ "github.com/jackc/pgx/v4"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{DB: db}
}

func (db *Database) Add(id, url string) error {
	panic("doesn't implemented")
}

func (db *Database) Get(id string) (string, error) {
	panic("doesn't implemented")
}

func (db *Database) GetList() ([]models.UserURL, error) {
	panic("doesn't implemented")
}

func (db Database) Ping() error {
	err := db.DB.Ping()
	return err
}
