package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/ChristinaFomenko/URLShortener/internal/handlers/mocks"
)

func TestShortenHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		shortcut    string
	}
	tests := []struct {
		name     string
		request  string
		url      string
		shortcut string
		want     want
	}{
		{
			name:     "success",
			url:      "https://yandex.ru",
			shortcut: "http://localhost:8080/abcde",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
				shortcut:    "http://localhost:8080/abcde",
			},
			request: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock.NewMockservice(ctrl)
			serviceMock.EXPECT().Shorten(tt.url).Return(tt.shortcut)

			httpHandler := New(serviceMock)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.url)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			writer := httptest.NewRecorder()
			HandlerFunc := http.HandlerFunc(httpHandler.Shorten)
			HandlerFunc.ServeHTTP(writer, request)
			result := writer.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			bodyResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.shortcut, string(bodyResult))
			assert.JSONEq(t, tt.want.shortcut, string(bodyResult))
		})
	}
}

const (
	shortURLDomain = "http://localhost:8080"
	longURL        = "https://www.yandex.ru/practicum"
)

func TestAPIShortenerHandler_Shorten_ShouldReturnBadRequestWhenShortenRequestIsNotContainsValidJSON(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mock.NewMockservice(controller)
	mockShortenerService.EXPECT().APIShortener(gomock.Any()).Return("", nil).Times(0)

	handler := handler{
		service: mockShortenerService,
	}
	resp := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader([]byte(`invalid json`)))

	handler.APIJSONShortener(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestAPIShortenerHandler_Shorten_ShouldReturnInternalServerErrorWhenShortenerServiceReturnsError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mock.NewMockservice(controller)
	mockShortenerService.EXPECT().APIShortener(gomock.Any()).Return("", errors.New("service error")).Times(1)

	handler := handler{
		service: mockShortenerService,
	}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader([]byte(fmt.Sprintf(`{"url": "%s"}`, longURL)))) //

	handler.APIJSONShortener(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestAPIShortenerHandler_Shorten_ShortenedURL(t *testing.T) {
	expectedShortenedURL := fmt.Sprintf(`{"result":"%s/tTeEsT"}`, shortURLDomain)
	shortenedURL := shortURLDomain + "/tTeEsT"
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockShortenerService := mock.NewMockservice(controller)
	mockShortenerService.EXPECT().APIShortener(gomock.Any()).Return(shortenedURL, nil).Times(1)

	handler := handler{
		service: mockShortenerService,
	}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader([]byte(fmt.Sprintf(`{"url": "%s"}`, longURL))))

	handler.APIJSONShortener(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, expectedShortenedURL, resp.Body.String())
	assert.JSONEq(t, expectedShortenedURL, resp.Body.String())
}
