package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/ChristinaFomenko/shortener/internal/app/models"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go

type service interface {
	Shorten(url string) (string, error)
	Expand(id string) (string, error)
	GetByUserID(UserID string) (string, error)
	Ping() error
}

type handler struct {
	service service
	db      *sql.DB
}

func New(service service, db *sql.DB) *handler {
	return &handler{
		service: service,
		db:      db,
	}
}

// Shorten Cut URL
func (h *handler) Shorten(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to validate struct", 400)
		return
	}

	url := string(bytes)
	shortcut, err := h.service.Shorten(url)
	if err != nil {
		log.WithError(err).WithField("url", url).Error("shorten url error")
		http.Error(w, "url shortcut", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortcut))
	if err != nil {
		log.WithError(err).WithField("shortcut", shortcut).Error("write response error")
		return
	}
}

// Expand Returns full URL by ID of shorted one
func (h *handler) Expand(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id parameter is empty", http.StatusBadRequest)
		return
	}

	url, err := h.service.Expand(id)
	if err != nil {
		http.Error(w, "url not found", http.StatusNoContent)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *handler) APIJSONShorten(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	req := models.ShortenRequest{}
	if err = json.Unmarshal(b, &req); err != nil {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(req)
	if err != nil || !ok {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	shortcut, err := h.service.Shorten(req.URL)
	if err != nil {
		log.WithError(err).WithField("url", req.URL).Error("shorten url error")
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := models.ShortenReply{ShortenURLResult: shortcut}
	marshal, err := json.Marshal(&resp)
	if err != nil {
		log.WithError(err).WithField("resp", resp).Error("marshal response error")
		http.Error(w, err.Error(), 400)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		log.WithError(err).WithField("shortcut", shortcut).Error("write response error")
		return
	}
}

func (h *handler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	idCookie, err := r.Cookie("user_id")
	w.Header().Set("Content-Type", "application/json")

	req := models.ResponseEntity{}

	urls, err := h.service.GetByUserID(idCookie.Value)
	if err != nil {
		log.WithError(err).WithField("url", req.OriginalURL).Error("original url error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	resp := models.ResponseEntity{OriginalURL: urls}
	body, err := json.Marshal(&resp)
	if err != nil {
		log.WithError(err).WithField("resp", urls).Error("marshal urls response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = w.Write(body)
	if err != nil {
		log.WithError(err).WithField("shortcut", urls).Error("write response error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
