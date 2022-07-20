package urls

import (
	"context"
	"errors"
	"fmt"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	errs "github.com/ChristinaFomenko/shortener/pkg/errors"
	_ "github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=urls.go -destination=mocks/mocks.go

const idLength int64 = 5

type urlRepository interface {
	Add(ctx context.Context, urlID, url, userID string) error
	Get(ctx context.Context, urlID string) (string, error)
	FetchURLs(ctx context.Context, userID string) ([]models.UserURL, error)
	AddBatch(ctx context.Context, urls []models.UserURL, userID string) error
	DeleteUserURLs(ctx context.Context, userID string, urls []string) error
}

type generator interface {
	Letters(n int64) (string, error)
}

type service struct {
	repository urlRepository
	generator  generator
	host       string
}

func NewService(repository urlRepository, generator generator, host string) *service {
	return &service{
		repository: repository,
		generator:  generator,
		host:       host,
	}
}

func (s *service) Shorten(ctx context.Context, url, userID string) (string, error) {
	urlID, err := s.generator.Letters(idLength)
	if err != nil {
		log.WithError(err).
			WithField("userID", userID).
			WithField("url", url).Error("add url error")
		return "", err
	}

	if err = s.repository.Add(ctx, urlID, url, userID); err != nil {
		var uniqueErr *errs.NotUniqueURLErr
		if errors.As(err, &uniqueErr) {
			return s.buildShortURL(uniqueErr.URLID), errs.ErrNotUniqueURL
		}

		log.WithError(err).
			WithField("userID", userID).
			WithField("urlID", urlID).
			WithField("url", url).
			Error("add url error")
		return "", err

	}

	return s.buildShortURL(urlID), nil
}

// Return by id

func (s *service) Expand(ctx context.Context, urlID string) (string, error) {
	url, err := s.repository.Get(ctx, urlID)
	if err != nil {
		if errors.Is(err, errs.ErrURLNotFound) {
			return "", errs.ErrURLNotFound
		}
		log.WithError(err).WithField("urlID", urlID).Error("get url error")
		return "", err
	}

	return url, nil
}

func (s *service) FetchURLs(ctx context.Context, userID string) ([]models.UserURL, error) {
	urls, err := s.repository.FetchURLs(ctx, userID)
	if err != nil {
		log.WithError(err).WithField("urlID", userID).Error("get url list error")
		return nil, err
	}

	for idx := range urls {
		urls[idx].ShortURL = s.buildShortURL(urls[idx].ShortURL)
	}

	return urls, nil
}

func (s *service) ShortenBatch(ctx context.Context, originalURLs []models.OriginalURL, userID string) ([]models.UserURL, error) {
	urls := make([]models.UserURL, len(originalURLs))
	for idx := range urls {
		urlID, err := s.generator.Letters(idLength)
		if err != nil {
			log.WithError(err).
				WithField("userID", userID).
				WithField("originalURLs", originalURLs).
				Error("generate urlID error")
			return nil, err
		}
		urls[idx] = models.UserURL{
			CorrelationID: originalURLs[idx].CorrelationID,
			ShortURL:      urlID,
			OriginalURL:   originalURLs[idx].URL,
		}
	}

	err := s.repository.AddBatch(ctx, urls, userID)
	if err != nil {
		log.WithError(err).
			WithField("userID", userID).
			WithField("originalURLs", originalURLs).
			WithField("urls", urls).
			Error("add urls batch error")
		return nil, err
	}

	for idx := range urls {
		urls[idx].ShortURL = s.buildShortURL(urls[idx].ShortURL)
	}

	return urls, nil
}

func (s *service) buildShortURL(id string) string {
	return fmt.Sprintf("%s/%s", s.host, id)
}

func (s *service) DeleteUserURLs(ctx context.Context, userID string, urls []string) error {
	usr, err := s.repository.Get(ctx, userID)
	if err != nil {
		log.WithError(err).
			WithField("userID", userID).
			Error("get userID error")
		return err
	}

	err = s.repository.DeleteUserURLs(ctx, usr, urls)
	if err != nil {
		log.WithError(err).
			WithField("userID", usr).
			WithField("urls", urls).
			Error("delete user urls error")
		return err
	}

	return nil
}
