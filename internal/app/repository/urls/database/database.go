package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	errs "github.com/ChristinaFomenko/shortener/pkg/errors"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"time"
)

const timeout = time.Second * 3

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

func (r *pgRepo) DeleteUserURLs(ctx context.Context, userID string, urls []string) error {
	//TODO implement me
	panic("implement me")
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

	_, err := r.db.ExecContext(ctx, `insert into urls(id,url,user_id) values ($1,$2,$3)`,
		urlID,
		url,
		&userID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
			err = r.db.QueryRowContext(ctx, "select id from urls where url=$1", url).Scan(&urlID)
			if err != nil {
				return err
			}
			err = errs.NewNotUniqueURLErr(urlID, url, err)
			return err
		}

		return err
	}

	return nil

}

func (r *pgRepo) Get(ctx context.Context, urlID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var url sql.NullString
	_ = r.db.QueryRowContext(ctx, `select url from urls where id=$1 and deleted_at is null`, urlID).Scan(&url)
	if url.Valid {
		return url.String, nil
	}

	return "", errs.ErrURLNotFound
}

func (r *pgRepo) FetchURLs(ctx context.Context, userID string) ([]models.UserURL, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res := make([]models.UserURL, 0)
	rows, err := r.db.QueryContext(ctx, `select id, url from urls where user_id=$1 and deleted_at is null;`, userID)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
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

func (r *pgRepo) AddBatch(ctx context.Context, urls []models.UserURL, userID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	stmt, err := tx.PrepareContext(ctx, `insert into urls(id,url,user_id) values ($1,$2,$3);`)
	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	for idx := range urls {
		if _, err = stmt.ExecContext(ctx, urls[idx].ShortURL, urls[idx].OriginalURL, userID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
