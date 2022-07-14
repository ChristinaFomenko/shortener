package ping

import (
	"context"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=ping.go -destination=mocks/mocks.go

type urlRepo interface {
	Ping(ctx context.Context) error
}

type service struct {
	urlRepo urlRepo
}

func NewService(urlRepo urlRepo) *service {
	return &service{
		urlRepo: urlRepo,
	}
}

func (s *service) Ping(ctx context.Context) bool {
	err := s.urlRepo.Ping(ctx)
	if err != nil {
		log.WithError(err).Error("ping database error")
		return false
	}

	return true
}
