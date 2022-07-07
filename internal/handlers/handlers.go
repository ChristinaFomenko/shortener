package handlers

import (
	"context"
	"encoding/json"
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
	"time"
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

	id, err := createURL(h, ctx, string(bytes))

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
	for _, shortID := range userURL {
		shortURL := fmt.Sprintf("%s/%s", h.Config.BaseURL, shortID)
		url, _ := h.Storage.GetURL(shortID.ShortURL)
		userURL = append(userURL, models.UserURL{ShortURL: shortURL, OriginalURL: url})
	}

	ctx.JSON(http.StatusOK, userURL)
}

func (h Handler) Ping(ctx *gin.Context) {
	timoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := h.Storage.GetDBConn().Ping(timoutCtx)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, "")
}

type batchShortenRequest []batchShortenRequest

func (h Handler) BatchShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read batch request body", http.StatusInternalServerError)
		return
	}

	var req batchShortenRequest

	if err := json.Unmarshal(b, &req); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
	}

	resp := make([]models.BatchShortenResponse, len(req))

	serializedResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Can't serialize response", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(serializedResp)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}

}

func createURL(h Handler, ctx *gin.Context, URL string) (shortURLID string, error error) {
	userEncryptID, err := ctx.Cookie("session")
	shortURLID = uuid.New().String()

	if userEncryptID != "" && err == nil {
		userDecryptID, err := utils.Decrypt(userEncryptID, h.Config.AuthKey)
		if err != nil {
			return "", err
		}

		if err := h.Storage.AddURL(shortURLID, URL, userDecryptID); err != nil {
			ctx.String(http.StatusBadRequest, "")
			return
		}
	} else {
		if err := h.Storage.AddURL(shortURLID, URL, ""); err != nil {
			ctx.String(http.StatusBadRequest, "")
			return
		}
	}

	return shortURLID, nil
}
