package handlers

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestAPIShorten(t *testing.T) {
	type want struct {
		contentType string
		body        string
		statusCode  int
	}
	testCases := []struct {
		name    string
		body    string
		request string
		want    want
	}{
		{
			name:    "#1 Valid request",
			body:    `{"url" : "http://shetube.com"}`,
			request: "/api/shorten",
			want: want{
				contentType: "application/json",
				body:        `{"result":"http://localhost:8080/"}`,
				statusCode:  http.StatusCreated,
			},
		},
		{
			name:    "#2 Invalid url",
			body:    `{"url" : "hetube.com"}`,
			request: "/api/shorten",
			want: want{
				contentType: "text/plain; charset=utf-8",
				body:        "Wrong URL",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "#3 Invalid json",
			body:    "{url : http://wetube.com}",
			request: "/api/shorten",
			want: want{
				contentType: "text/plain; charset=utf-8",
				body:        "Bad request",
				statusCode:  http.StatusBadRequest,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock.NewMockservice(ctrl)
			serviceMock.EXPECT().APIShortener(tt.name).Return(tt.want)

			httpHandler := New(serviceMock)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			writer := httptest.NewRecorder()
			HandlerFunc := http.HandlerFunc(httpHandler.APIShortener)
			HandlerFunc.ServeHTTP(writer, request)
			result := writer.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			if result.StatusCode == http.StatusCreated {
				assert.JSONEq(t, tt.want.body, string(body))
			} else {
				assert.Equal(t, tt.want.body, strings.TrimSpace(string(body)))
			}
		})
	}
}
