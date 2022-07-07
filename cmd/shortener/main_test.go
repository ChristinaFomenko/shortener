package main

import (
	"bytes"
	"encoding/json"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/app/storage"
	"github.com/ChristinaFomenko/shortener/internal/models"
	"github.com/ChristinaFomenko/shortener/internal/router"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/caarlos0/env/v6"
)

func TestShortenHandler(t *testing.T) {
	c := configs.AppConfig{}
	if err := env.Parse(&c); err != nil {
		return
	}

	s := storage.ConstructStorage(c)

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "URL added success",
			want: want{
				http.StatusCreated,
				"text/plain; charset=utf-8",
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := router.Router(c, s)
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("https://yandex.ru"))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.want.code, w.Code)
			assert.Equal(t, testCase.want.contentType, w.Header().Get("Content-Type"))
			assert.NoError(t, err)
		})
	}
}

func TestAPIJSONShorten_Success(t *testing.T) {
	c := configs.AppConfig{}
	if err := env.Parse(&c); err != nil {
		return
	}

	s := storage.ConstructStorage(c)

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "URL JSON added",
			want: want{
				http.StatusCreated,
				"application/json; charset=utf-8",
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := router.Router(c, s)
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(models.ShortenRequest{URL: "https://yandex.ru"})
			req, err := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.want.code, w.Code)
			assert.Equal(t, testCase.want.contentType, w.Header().Get("Content-Type"))
			assert.NoError(t, err)
		})
	}
}
