package handlers

import (
	"bytes"
	"context"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mock "github.com/ChristinaFomenko/shortener/internal/handlers/mocks"
)

const defaultUserID = "abcde"

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
		err      error
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
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock.NewMockservice(ctrl)
			serviceMock.EXPECT().Shorten(ctx, tt.url, defaultUserID).Return(tt.shortcut, nil)

			authMock := mock.NewMockauth(ctrl)
			authMock.EXPECT().UserID(gomock.Any()).Return(defaultUserID)

			httpHandler := New(serviceMock, authMock, nil)

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

func TestAPIJSONShorten_Success(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name     string
		request  string
		url      string
		body     string
		shortcut string
		want     want
		err      error
	}{
		{
			name:     "success",
			url:      "https://yandex.ru",
			body:     "{\"url\":\"https://yandex.ru\"}",
			shortcut: "http://localhost:8080/abcde",
			want: want{
				contentType: "application/json",
				statusCode:  201,
				response:    "{\"result\":\"http://localhost:8080/abcde\"}",
			},
			request: "/api/shorten",
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock.NewMockservice(ctrl)
			serviceMock.EXPECT().Shorten(ctx, tt.url, defaultUserID).Return(tt.shortcut, tt.err)

			authMock := mock.NewMockauth(ctrl)
			authMock.EXPECT().UserID(gomock.Any()).Return(defaultUserID)

			httpHandler := New(serviceMock, authMock, nil)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			writer := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(httpHandler.APIJSONShorten)
			handlerFunc.ServeHTTP(writer, request)
			result := writer.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			bodyResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(bodyResult))
		})
	}
}

func TestAPIJSONShorten_BadRequest(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name     string
		request  string
		url      string
		body     string
		shortcut string
		want     want
	}{
		{
			name:     "bad-request",
			url:      "https://yandex.ru",
			body:     "{\"url\":\"\"}",
			shortcut: "http://localhost:8080/abcde",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "request in not valid\n",
			},
			request: "/api/shorten",
		},
		{
			name:     "bad-request",
			url:      "https://yandex.ru",
			body:     "{\"url\":\"qwerty\"}",
			shortcut: "http://localhost:8080/abcde",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "request in not valid\n",
			},
			request: "/api/shorten",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			httpHandler := New(nil, nil, nil)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			writer := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(httpHandler.APIJSONShorten)
			handlerFunc.ServeHTTP(writer, request)
			result := writer.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			bodyResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(bodyResult))
		})
	}
}

func TestExpandHandler_Success(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
		location    string
	}
	tests := []struct {
		name     string
		request  string
		url      string
		urlID    string
		shortcut string
		err      error
		want     want
	}{
		{
			name:     "success",
			url:      "https://yandex.ru",
			urlID:    "abc",
			shortcut: "http://localhost:8080/abc",
			err:      nil,
			want: want{
				contentType: "",
				statusCode:  307,
				response:    "",
			},
			request: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlsSrvMock := mock.NewMockservice(ctrl)
			urlsSrvMock.EXPECT().Expand(gomock.Any(), tt.urlID).Return(tt.url, tt.err)

			httpHandler := New(urlsSrvMock, nil, nil)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlID)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.Expand)

			h.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(userResult))
		})
	}
}

func Test_handler_FetchURLs_Success(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name    string
		request string
		urls    []models.UserURL
		err     error
		want    want
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
			want: want{
				contentType: "application/json",
				statusCode:  200,
				response:    "[{\"short_url\":\"http://localhost:8080/abcde\",\"original_url\":\"https://yandex.ru\"},{\"short_url\":\"http://localhost:8080/qwerty\",\"original_url\":\"https://github.com\"}]",
			},
			request: "/api/user/urls",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock.NewMockservice(ctrl)
			serviceMock.EXPECT().FetchURLs(ctx, defaultUserID).Return(tt.urls, tt.err)

			authMock := mock.NewMockauth(ctrl)
			authMock.EXPECT().UserID(gomock.Any()).Return(defaultUserID)

			httpHandler := New(serviceMock, authMock, nil)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.FetchURLs)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(userResult))
		})
	}
}

func TestGetUrlsHandler_Ping(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		request string
		success bool
		want    want
	}{
		{
			name:    "success",
			success: true,
			want: want{
				statusCode: 200,
			},
			request: "/ping",
		},
		{
			name:    "fail",
			success: false,
			want: want{
				statusCode: 500,
			},
			request: "/ping",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			infraMock := mock.NewMockpingService(ctrl)
			infraMock.EXPECT().Ping(ctx).Return(tt.success)

			httpHandler := New(nil, nil, infraMock)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.Ping)
			h.ServeHTTP(w, request)

			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
