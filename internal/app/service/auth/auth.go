package auth

import log "github.com/sirupsen/logrus"

const idLength = 8

type generator interface {
	Letters(n int64) string
}

type hasher interface {
	Sign(value string) (string, error)
	Validate(value string, dataLength int64) (string, error)
}

type service struct {
	hasher    hasher
	generator generator
}

func NewService(generator generator, hasher hasher) *service {
	return &service{
		generator: generator,
		hasher:    hasher,
	}
}

func (s *service) SignUp() (string, string, error) {
	userID := s.generator.Letters(idLength)
	signedUserID, err := s.hasher.Sign(userID)
	if err != nil {
		log.WithError(err).WithField("userID", userID).Error("sign userID error")
		return userID, "", err
	}

	return userID, signedUserID, nil
}

func (s *service) SignIn(token string) (string, error) {
	userID, err := s.hasher.Validate(token, idLength)
	if err != nil {
		log.WithError(err).WithField("token", token).Error("validate userID sign error")
		return "", err
	}

	return userID, nil
}
