package storage

import (
	"context"
	"errors"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"time"
)

func (db *Database) AddURL(userShortURL UserURL) error {
	_, err := db.DB.Exec(context.Background(), "INSERT INTO urls VALUES ($1, $2, $3)", userShortURL.ID, userShortURL.OriginalURL, userShortURL.UserID)

	var pgError *pgconn.PgError

	if errors.As(err, &pgError) {
		if pgError.Code == pgerrcode.UniqueViolation {
			return &URLDuplicateError{URL: userShortURL.OriginalURL}
		}
	}

	return err
}

func (db *Database) GetURL(ID string) (string, error) {
	row := ""
	err := db.DB.QueryRow(context.Background(), "SELECT original_url FROM urls WHERE url_id = $1", ID).Scan(&row)
	if err != nil {
		return "", err
	}

	return row, nil
}

func (db *Database) GetList(userID string) []UserURL {
	shortURLs := make([]UserURL, 0)
	rows, err := db.DB.Query(context.Background(), "SELECT url_id, original_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		return shortURLs
	}
	defer rows.Close()

	for rows.Next() {
		var sURL UserURL
		err := rows.Scan(&sURL.ID, &sURL.OriginalURL)
		if err != nil {
			return nil
		}
		shortURLs = append(shortURLs, sURL)
	}

	return shortURLs
}

func ConstructDatabaseStorage(cfg configs.AppConfig) (Repository, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	db := &Database{DB: conn}

	const createTable = `
		CREATE TABLE IF NOT EXISTS urls (
			url_id varchar(36) NOT NULL UNIQUE PRIMARY KEY,
			original_url varchar(255) UNIQUE,
			user_id varchar(36)
		)`
	if _, err = db.DB.Exec(context.Background(), createTable); err != nil {
		return nil, err
	}

	return db, nil
}

func (db Database) DestructStorage(cfg configs.AppConfig) error {
	err := db.DB.Close(context.Background())

	return err
}

func (db Database) Ping() error {
	ctx := context.Background()
	conn, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := db.DB.Ping(conn)

	return err
}

func (db *Database) AddBatchURL(shortURLs []UserURL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stmt, err := db.DB.Prepare(ctx, "addBatch", "INSERT INTO urls VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	for _, su := range shortURLs {
		_, err := db.DB.Exec(ctx, stmt.SQL, su.ID, su.OriginalURL, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) GetShortByOriginal(originalURL string) (string, error) {
	var ID string
	if err := db.DB.QueryRow(context.Background(), "SELECT url_id FROM urls WHERE original_url = $1", originalURL).Scan(&ID); err != nil {
		return "", err
	}

	return ID, nil
}
