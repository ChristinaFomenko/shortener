package urls

import (
	"context"
	"errors"
	"fmt"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	_ "github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=urls.go -destination=mocks/mocks.go

const idLength int64 = 5

var ErrURLNotFound = errors.New("url not found error")

type urlRepository interface {
	Add(ctx context.Context, urlID, userID, url string) (string, error)
	Get(ctx context.Context, urlID string) (string, error)
	FetchURLs(ctx context.Context, userID string) ([]models.UserURL, error)
	AddBatch(ctx context.Context, urls []models.UserURL, userID string) error
}

type generator interface {
	Letters(n int64) string
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
	urlID := s.generator.Letters(idLength)
	url, err := s.repository.Add(ctx, urlID, userID, url)
	if err != nil {
		log.WithError(err).
			WithField("urlID", urlID).
			WithField("userID", userID).
			WithField("url", url).Error("add url error")
		return "", err
	}

	return s.buildShortURL(urlID), nil
}

// Return by id

func (s *service) Expand(ctx context.Context, urlID string) (string, error) {
	url, err := s.repository.Get(ctx, urlID)
	if err != nil {
		if errors.Is(err, ErrURLNotFound) {
			return "", ErrURLNotFound
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
		urls[idx] = models.UserURL{
			CorrelationID: originalURLs[idx].CorrelationID,
			ShortURL:      s.generator.Letters(idLength),
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
