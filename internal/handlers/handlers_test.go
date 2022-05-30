package handlers

import (
	"bytes"
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
		})
	}
}

func Test_handler_Expand(t *testing.T) {
	type want struct {
		location   string
		statusCode int
		shortcut   string
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
				location:   "location",
				statusCode: 307,
				shortcut:   "http://localhost:8080/abcde",
			},
			request: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock.NewMockservice(ctrl)
			serviceMock.EXPECT().Expand(tt.shortcut).Return(tt.url)

			httpHandler := New(serviceMock)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.url)
			request := httptest.NewRequest(http.MethodGet, tt.request, buffer)

			writer := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(httpHandler.Expand)
			handlerFunc.ServeHTTP(writer, request)
			result := writer.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.shortcut, string(userResult))
		})
	}
}
