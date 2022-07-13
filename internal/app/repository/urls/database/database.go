package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	_ "github.com/jackc/pgx/v4"
	"time"
)

const timeout = time.Second * 3

var ErrURLNotFound = errors.New("url not found")

type database interface {
	PingContext(ctx context.Context) error
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Close() error
	Begin() (*sql.Tx, error)
}

type pgRepo struct {
	db database
}

func NewRepo(dsn string) (*pgRepo, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Minute * 2)

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &pgRepo{
		db: db,
	}, nil
}

func (r *pgRepo) Add(ctx context.Context, urlID, url, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `insert into urls(id,url,user_id) values ($1,$2,$3);`, urlID, url, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *pgRepo) Get(ctx context.Context, urlID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var url sql.NullString
	_ = r.db.QueryRowContext(ctx, `select url from urls where id=$1 where deleted_at is null`, urlID).Scan(&url)
	if url.Valid {
		return url.String, nil
	}

	return "", ErrURLNotFound
}

func (r *pgRepo) FetchURLs(ctx context.Context, userID string) ([]models.UserURL, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res := make([]models.UserURL, 0)
	rows, _ := r.db.QueryContext(ctx, `select id, url from urls where user_id=$1 and deleted_at is null;`, userID)
	err := rows.Err()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var url models.UserURL
		err = rows.Scan(&url.ShortURL, &url.OriginalURL)
		if err != nil {
			return nil, err
		}

		res = append(res, url)
	}

	return res, nil
}

func (r *pgRepo) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return r.db.PingContext(ctx)
}

func (r *pgRepo) Close() error {
	return r.db.Close()
}
