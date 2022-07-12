package urls

import (
	"errors"
	"github.com/ChristinaFomenko/shortener/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"

	mocks "github.com/ChristinaFomenko/shortener/internal/app/service/urls/mocks"
)

const (
	host          = "http://localhost:8080"
	defaultUserID = "abcde"
)

func Test_service_Shorten(t *testing.T) {
	tests := []struct {
		name     string
		urlID    string
		url      string
		shortcut string
		err      error
	}{
		{
			name:     "success",
			urlID:    "abcde",
			url:      "yandex.ru",
			shortcut: "http://localhost:8080/abcde",
		},
		{
			name:     "success",
			urlID:    "abcde",
			url:      "yandex.ru",
			shortcut: "",
			err:      errors.New("test err"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		generatorMock := mocks.NewMockgenerator(ctrl)
		generatorMock.EXPECT().Letters(idLength).Return(tt.urlID)

		repositoryMock := mocks.NewMockurlRepository(ctrl)
		repositoryMock.EXPECT().Add(tt.urlID, defaultUserID, tt.url).Return(tt.err)

		s := NewService(repositoryMock, generatorMock, host, nil)
		act, err := s.Shorten(tt.url, defaultUserID)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.shortcut, act)
	}
}

func Test_service_Expand(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		shortcut string
		err      error
	}{
		{
			name:     "success",
			url:      "yandex.ru",
			shortcut: "abcde",
			err:      nil,
		},
		{
			name:     "error",
			url:      "",
			shortcut: "abcde",
			err:      errors.New("test error"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repositoryMock := mocks.NewMockurlRepository(ctrl)
		repositoryMock.EXPECT().Get(tt.shortcut).Return(tt.url, tt.err)

		s := NewService(repositoryMock, nil, host, nil)
		act, err := s.Expand(tt.shortcut)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.url, act)
	}
}

func Test_service_FetchURls(t *testing.T) {
	tests := []struct {
		name string
		urls []models.UserURL
		err  error
	}{
		{
			name: "success",
			urls: []models.UserURL{
				{
					ShortURL:    "http://localhost:8080/abcde",
					OriginalURL: "https://yandex.ru",
				},
				{
					ShortURL:    "http://localhost:8080/qwerty",
					OriginalURL: "https://github.com",
				},
			},
			err: nil,
		},
		{
			name: "repo err",
			urls: nil,
			err:  errors.New("test err"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repositoryMock := mocks.NewMockurlRepository(ctrl)
		repositoryMock.EXPECT().FetchURls(defaultUserID).Return(tt.urls, tt.err)

		s := NewService(repositoryMock, nil, host, nil)
		act, err := s.GetList(defaultUserID)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.urls, act)
	}
}
