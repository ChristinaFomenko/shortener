package urls

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=urls.go -destination=mocks/mocks.go

type urlRepository interface {
	Add(id, url string)
	Get(id string) (string, error)
	Post(id string) (string, error)
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

func (s *service) Shorten(url string) string {
	id := s.generator.GenerateID()
	s.repository.Add(id, url)

	return fmt.Sprintf("%s/%s", s.host, id)
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

func (s *service) APIShortener(id string) (string, error) {
	post, err := s.repository.Post(id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("api shortener url error")
		return "", err
	}

	return fmt.Sprintf("%s/%s", post, id), nil
}
