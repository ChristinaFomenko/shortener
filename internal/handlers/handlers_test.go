package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/app/mocks"
	"github.com/ChristinaFomenko/shortener/internal/app/storage"
	"github.com/ChristinaFomenko/shortener/internal/models"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Shorten(t *testing.T) {
	cfg := configs.AppConfig{}
	if err := env.Parse(&cfg); err != nil {
		return
	}

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name         string
		want         want
		mockBehavior func(*mocks.MockRepository)
	}{
		{
			name: "201 URL created",
			want: want{
				http.StatusCreated,
				"text/plain",
			},
			mockBehavior: func(s *mocks.MockRepository) {
				s.EXPECT().AddURL(gomock.Any()).Return(nil)
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			stor := mocks.NewMockRepository(ctrl)
			testCase.mockBehavior(stor)

			handler := Handler{Config: cfg, Storage: stor}

			r := gin.Default()
			r.POST("/", handler.Shorten)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("https://test1.ru"))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.want.code, w.Code)
			assert.Equal(t, testCase.want.contentType, w.Header().Get("Content-Type"))
			assert.NoError(t, err)
		})
	}
}

func TestHandler_APIJSONShorten(t *testing.T) {
	cfg := configs.AppConfig{}
	if err := env.Parse(&cfg); err != nil {
		return
	}

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name         string
		want         want
		request      string
		mockBehavior func(*mocks.MockRepository)
	}{
		{
			name: "201 URL created",
			want: want{
				http.StatusCreated,
				"application/json; charset=utf-8",
			},
			mockBehavior: func(s *mocks.MockRepository) {
				s.EXPECT().AddURL(gomock.Any()).Return(nil)
			},
		},
		{
			name: "500 URL creation error",
			want: want{
				http.StatusInternalServerError,
				"text/plain; charset=utf-8",
			},
			mockBehavior: func(s *mocks.MockRepository) {
				s.EXPECT().AddURL(gomock.Any()).Return(errors.New(""))
			},
		},
		{
			name: "409 URL creation duplicate error",
			want: want{
				http.StatusConflict,
				"application/json; charset=utf-8",
			},
			mockBehavior: func(s *mocks.MockRepository) {
				s.EXPECT().AddURL(gomock.Any()).Return(&storage.URLDuplicateError{URL: ""})
				s.EXPECT().GetShortByOriginal(gomock.Any()).Return("test", nil)
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			stor := mocks.NewMockRepository(ctrl)
			testCase.mockBehavior(stor)

			handler := Handler{Config: cfg, Storage: stor}

			r := gin.Default()
			r.POST("/api/shorten", handler.APIJSONShorten)

			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(models.ShortenRequest{URL: "https://test3.ru"})
			req, err := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.want.code, w.Code)
			assert.Equal(t, testCase.want.contentType, w.Header().Get("Content-Type"))
			assert.NoError(t, err)
		})
	}
}
