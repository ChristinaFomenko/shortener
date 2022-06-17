package urls

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mocks "github.com/ChristinaFomenko/shortener/internal/app/service/urls/mocks"
)

const host = "http://localhost:8080"

func TestShorten(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		url      string
		shortcut string
		err      error
	}{
		{
			name:     "success",
			id:       "abcde",
			url:      "yandex.ru",
			shortcut: "http://localhost:8080/abcde",
		},
		{
			name:     "success",
			id:       "abcde",
			url:      "yandex.ru",
			shortcut: "",
			err:      errors.New("test err"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		generatorMock := mocks.NewMockgenerator(ctrl)
		generatorMock.EXPECT().GenerateID().Return(tt.id)

		repositoryMock := mocks.NewMockurlRepository(ctrl)
		repositoryMock.EXPECT().Add(tt.id, tt.url).Return(tt.err)

		s := NewService(repositoryMock, generatorMock, host)
		act, err := s.Shorten(tt.url)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.shortcut, act)
	}
}

func TestExpand(t *testing.T) {
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

		s := NewService(repositoryMock, nil, host)
		act, err := s.Expand(tt.shortcut)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.url, act)
	}
}
