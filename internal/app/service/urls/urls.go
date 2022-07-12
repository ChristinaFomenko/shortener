package urls

import (
	"fmt"
	"github.com/ChristinaFomenko/shortener/internal/app/repository/urls/database"
	"github.com/ChristinaFomenko/shortener/internal/models"
	_ "github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=urls.go -destination=mocks/mocks.go

const idLength int64 = 5

type urlRepository interface {
	Add(urlID, userID, url string) error
	Get(urlID string) (string, error)
	GetList(userID string) ([]models.UserURL, error)
	Ping() error
}

type generator interface {
	Letters(n int64) string
}

type service struct {
	repository urlRepository
	generator  generator
	host       string
	db         *database.Database
}

func NewService(repository urlRepository, generator generator, host string, db *database.Database) *service {
	return &service{
		repository: repository,
		generator:  generator,
		host:       host,
		db:         db,
	}
}

func (s *service) Shorten(url, userID string) (string, error) {
	urlID := s.generator.Letters(idLength)
	err := s.repository.Add(urlID, userID, url)
	if err != nil {
		log.WithError(err).
			WithField("urlID", urlID).
			WithField("userID", userID).
			WithField("url", url).Error("add url error")
		return "", err
	}

	return s.buildShortString(urlID), nil
}

// Return by id

func (s *service) Expand(urlID string) (string, error) {
	url, err := s.repository.Get(urlID)
	if err != nil {
		log.WithError(err).WithField("urlID", urlID).Error("get url error")
		return "", err
	}

	return url, nil
}

func (s *service) GetList(userID string) ([]models.UserURL, error) {
	urls, err := s.repository.GetList(userID)
	if err != nil {
		log.WithError(err).WithField("urlID", userID).Error("get url list error")
		return nil, err
	}

	for idx := range urls {
		urls[idx].ShortURL = s.buildShortString(urls[idx].ShortURL)
	}

	return urls, nil
}

func (s *service) Ping() error {
	return s.db.Ping()
}

func (s *service) buildShortString(id string) string {
	return fmt.Sprintf("%s/%s", s.host, id)
}
