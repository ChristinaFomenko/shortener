package urls

import (
	"fmt"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=urls.go -destination=mocks/mocks.go

type urlRepository interface {
	Add(id, url string) error
	Get(id string) (string, error)
	GetList() ([]models.UserURL, error)
	Ping() error
}

type generator interface {
	GenerateID() string
}

type service struct {
	repository urlRepository
	generator  generator
	host       string
}

func NewService(
	repository urlRepository,
	generator generator,
	host string) *service {
	return &service{
		repository: repository,
		generator:  generator,
		host:       host,
	}
}

func (s *service) Shorten(url string) (string, error) {
	id := s.generator.GenerateID()
	err := s.repository.Add(id, url)
	if err != nil {
		log.WithError(err).WithField("id", id).WithField("url", url).Error("add url error")
		return "", err
	}

	return fmt.Sprintf("%s/%s", s.host, id), nil
}

// Return by id

func (s *service) Expand(id string) (string, error) {
	url, err := s.repository.Get(id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("get url error")
		return "", err
	}

	return url, nil
}

func (s *service) GetList() ([]models.UserURL, error) {
	urls, err := s.repository.GetList()
	if err != nil {
		log.WithError(err).Error("get url list error")
		return nil, err
	}

	return urls, nil
}

func (s *service) Ping() error {
	return fmt.Errorf("method doesn't implemented")
}
