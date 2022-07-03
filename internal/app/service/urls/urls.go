package urls

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=urls.go -destination=mocks/mocks.go

type urlRepository interface {
	Add(id, url string) error
	Get(id string) (string, error)
	GetByUserID(userID string) (string, error)
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

func (s *service) Ping() error {
	//TODO implement me
	panic("implement me")
}

func (s *service) GetByUserID(UserID string) (string, error) {
	user, err := s.repository.GetByUserID(UserID)
	if err != nil {
		log.WithError(err).WithField("userID", UserID).Error("get user ID error")
		return "", err
	}

	return user, nil
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
