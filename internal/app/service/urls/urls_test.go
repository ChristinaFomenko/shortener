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

func Test_service_ShortenBatch(t *testing.T) {
	tests := []struct {
		name         string
		originalURLs []models.OriginalURL
		urls         []models.UserURL
		err          error
		exp          []models.UserURL
	}{
		{
			name: "success",
			originalURLs: []models.OriginalURL{
				{
					CorrelationID: "1",
					URL:           "https://yandex.ru",
				},
				{
					CorrelationID: "2",
					URL:           "https://github.com",
				},
			},
			urls: []models.UserURL{
				{
					CorrelationID: "1",
					ShortURL:      "abcde",
					OriginalURL:   "https://yandex.ru",
				},
				{
					CorrelationID: "2",
					ShortURL:      "qwerty",
					OriginalURL:   "https://github.com",
				},
			},
			exp: []models.UserURL{
				{
					CorrelationID: "1",
					ShortURL:      "http://localhost:8080/abcde",
					OriginalURL:   "https://yandex.ru",
				},
				{
					CorrelationID: "2",
					ShortURL:      "http://localhost:8080/qwerty",
					OriginalURL:   "https://github.com",
				},
			},
			err: nil,
		},
		{
			name: "repo err",
			originalURLs: []models.OriginalURL{
				{
					CorrelationID: "1",
					URL:           "https://yandex.ru",
				},
				{
					CorrelationID: "2",
					URL:           "https://github.com",
				},
			},
			urls: []models.UserURL{
				{
					CorrelationID: "1",
					ShortURL:      "abcde",
					OriginalURL:   "https://yandex.ru",
				},
				{
					CorrelationID: "2",
					ShortURL:      "qwerty",
					OriginalURL:   "https://github.com",
				},
			},
			err: errors.New("test err"),
		},
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repositoryMock := mocks.NewMockurlRepository(ctrl)
		repositoryMock.EXPECT().AddBatch(ctx, tt.urls, defaultUserID).Return(tt.err)

		generatorMock := mocks.NewMockgenerator(ctrl)
		for _, url := range tt.urls {
			generatorMock.EXPECT().Letters(idLength).Return(url.ShortURL)
		}

		s := NewService(repositoryMock, generatorMock, host)
		act, err := s.ShortenBatch(ctx, tt.originalURLs, defaultUserID)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.exp, act)
	}
}
