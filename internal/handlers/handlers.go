package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/models"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var cfg = configs.AppConfig{}

var (
	ErrNoTokenFound = errors.New("no token found")
	ErrInvalidToken = errors.New("token is invalid")
)

//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go

type service interface {
	Shorten(url string) (string, error)
	Expand(id string) (string, error)
	GetList() ([]models.UserURL, error)
	Ping() error
}

type handler struct {
	service service
}

func New(service service) *handler {
	return &handler{
		service: service,
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
		log.WithError(err).WithField("id", shortcut).Error("write response error")
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

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	urls, err := h.service.GetList()
	if err != nil {
		log.WithError(err).Error("get urls error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("content-type", "application/json")

	resp := toGetUrlsReply(urls)
	body, err := json.Marshal(&resp)
	if err != nil {
		log.WithError(err).WithField("resp", urls).Error("marshal urls response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = w.Write(body)
	if err != nil {
		log.WithError(err).WithField("resp", urls).Error("write response error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.service.Ping()
	if err != nil {
		log.Infof("DB not avalable %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type batchShortenRequest []batchShortenRequest

func (h *handler) BatchShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read batch request body", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{Name: "user_id", Value: "abcd"}
	if err != nil {
		if errors.Is(err, ErrNoTokenFound) || errors.Is(err, ErrInvalidToken) {
			http.SetCookie(w, &cookie)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
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
