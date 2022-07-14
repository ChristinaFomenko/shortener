package urls

import (
	"context"
	"errors"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
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

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		generatorMock := mocks.NewMockgenerator(ctrl)
		generatorMock.EXPECT().Letters(idLength).Return(tt.urlID)

		repositoryMock := mocks.NewMockurlRepository(ctrl)
		repositoryMock.EXPECT().Add(ctx, tt.urlID, tt.url, defaultUserID).Return(tt.url, tt.err)

		s := NewService(repositoryMock, generatorMock, host)
		act, err := s.Shorten(ctx, tt.url, defaultUserID)

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

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repositoryMock := mocks.NewMockurlRepository(ctrl)
		repositoryMock.EXPECT().Get(ctx, tt.shortcut).Return(tt.url, tt.err)

		s := NewService(repositoryMock, nil, host)
		act, err := s.Expand(ctx, tt.shortcut)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.url, act)
	}
}

func Test_service_FetchURLs(t *testing.T) {
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

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repositoryMock := mocks.NewMockurlRepository(ctrl)
		repositoryMock.EXPECT().FetchURLs(ctx, defaultUserID).Return(tt.urls, tt.err)

		s := NewService(repositoryMock, nil, host)
		act, err := s.FetchURLs(ctx, defaultUserID)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.urls, act)
	}
}

func Test_service_Ping(t *testing.T) {
	tests := []struct {
		name string
		err  error
		exp  bool
	}{
		{
			name: "success",
			err:  nil,
			exp:  true,
		},
		{
			name: "repo err",
			err:  errors.New("test err"),
			exp:  false,
		},
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repoMock := mocks.NewMockurlRepository(ctrl)
		repoMock.EXPECT().Ping(ctx).Return(tt.err)

		s := NewService(repoMock, nil, host)
		act := s.Ping(ctx)

		assert.Equal(t, tt.exp, act)
	}
}
