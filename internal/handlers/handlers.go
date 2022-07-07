package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/app/storage"
	"github.com/ChristinaFomenko/shortener/internal/app/utils"
	"github.com/ChristinaFomenko/shortener/internal/models"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go

type Handler struct {
	Config  configs.AppConfig
	Storage storage.Repository
}

// Shorten Cut URL
func (h Handler) Shorten(ctx *gin.Context) {
	bytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error(err)
		ctx.String(http.StatusInternalServerError, "failed to validate struct")
		return
	}

	url := string(bytes)
	id, err := createURL(h, ctx, string(bytes))
	if err != nil {
		var urlDupl *storage.URLDuplicateError
		if errors.As(err, &urlDupl) {
			existID, _ := h.Storage.GetShortByOriginal(url)
			existURL := fmt.Sprintf("%s/%s", h.Config.BaseURL, existID)
			ctx.String(http.StatusConflict, existURL)
			return
		}
	}

	shortcut := fmt.Sprintf("%s/%s", h.Config.BaseURL, id)
	if err != nil {
		log.WithError(err).WithField("url", shortcut).Error("shorten url error")
		ctx.String(http.StatusInternalServerError, "url shortcut")
		return
	}

	ctx.Header("content-type", "text/plain; charset=utf-8")
	ctx.String(http.StatusCreated, "%s", shortcut)

}

// Expand Returns full URL by ID of shorted one
func (h Handler) Expand(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.String(http.StatusBadRequest, "id parameter is empty")
		return
	}

	originalURL, err := h.Storage.GetURL(id)
	if err != nil {
		ctx.String(http.StatusNoContent, "url not found")
		return
	}
	ctx.Header("content-type", "text/plain")

	ctx.Redirect(http.StatusTemporaryRedirect, originalURL)
}

func (h Handler) APIJSONShorten(ctx *gin.Context) {
	req := models.ShortenRequest{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.String(http.StatusBadRequest, "request in not valid")
		return
	}

	ok, err := govalidator.ValidateStruct(req)
	if err != nil || !ok {
		ctx.String(http.StatusBadRequest, "request in not valid")
		return
	}

	id, err := createURL(h, ctx, req.URL)
	if err != nil {
		var urlDupl *storage.URLDuplicateError
		if errors.As(err, &urlDupl) {
			existID, _ := h.Storage.GetShortByOriginal(req.URL)
			existURL := fmt.Sprintf("%s/%s", h.Config.BaseURL, existID)
			res := models.ShortenReply{ShortenURLResult: existURL}

			ctx.JSON(http.StatusConflict, res)
			return
		}
		ctx.String(http.StatusInternalServerError, "failed to create url")
		return
	}

	shortcut := fmt.Sprintf("%s/%s", h.Config.BaseURL, id)

	ctx.Header("content-type", "application/json; charset=utf-8")

	resp := models.ShortenReply{ShortenURLResult: shortcut}

	ctx.JSON(http.StatusCreated, resp)
}

func (h Handler) GetList(ctx *gin.Context) {
	userID, err := ctx.Cookie("user_id")
	if err != nil {
		log.WithError(err).Error("failed to set cookie")
		ctx.String(http.StatusInternalServerError, "failed to set cookie")
		return
	}

	userDecryptID, err := utils.Decrypt(userID, h.Config.AuthKey)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "failed decrypt user id")
		return
	}

	urls := h.Storage.GetList(userDecryptID)
	if len(urls) == 0 {
		ctx.JSON(http.StatusNoContent, "{}")
		return
	}

	ctx.Header("content-type", "application/json")

	var userURL []models.UserURL
	for _, shorted := range urls {
		shortURL := fmt.Sprintf("%s/%s", h.Config.BaseURL, shorted.ID)
		userURL = append(userURL, models.UserURL{ShortURL: shortURL, OriginalURL: shorted.OriginalURL})
	}

	ctx.JSON(http.StatusOK, userURL)
}

func (h Handler) Ping(ctx *gin.Context) {
	err := h.Storage.Ping()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, "Database is running")
}

type batchShortenRequest []batchShortenRequest

func (h Handler) BatchShortenHandler(ctx *gin.Context) {
	var batchURLs models.BatchURLs
	userID, _ := ctx.Cookie("session")

	if err := json.NewDecoder(ctx.Request.Body).Decode(&batchURLs); err != nil {
		ctx.String(http.StatusInternalServerError, "batch decoding error")
		return
	}

	shortURLs := make([]storage.UserURL, 0)

	for _, bu := range batchURLs {
		shortURLs = append(shortURLs, storage.UserURL{ID: bu.CorrelationID, OriginalURL: bu.OriginalURL, UserID: userID})
	}

	if err := h.Storage.AddBatchURL(shortURLs); err != nil {
		ctx.String(http.StatusInternalServerError, "batch add error")
		return
	}

	batchShortURL := make([]models.BatchShortURL, 0)

	for _, su := range shortURLs {
		batchShortURL = append(batchShortURL, models.BatchShortURL{
			CorrelationID: su.ID,
			ShortURL:      fmt.Sprintf("%s/%s", h.Config.BaseURL, su.ID),
		})
	}

	ctx.JSON(http.StatusCreated, batchShortURL)
}

func createURL(h Handler, ctx *gin.Context, URL string) (shortURLID string, error error) {
	userEncryptID, err := ctx.Cookie("session")
	shortURLID = uuid.New().String()

	if userEncryptID != "" && err == nil {
		userDecryptID, err := utils.Decrypt(userEncryptID, h.Config.AuthKey)
		if err != nil {
			return "", err
		}

		if err := h.Storage.AddURL(storage.UserURL{ID: shortURLID, OriginalURL: URL, UserID: userDecryptID}); err != nil {
			return "", err
		}
	} else {
		if err := h.Storage.AddURL(storage.UserURL{ID: shortURLID, OriginalURL: URL, UserID: ""}); err != nil {
			ctx.String(http.StatusBadRequest, "")
			return "", err
		}
	}

	return shortURLID, nil
}
