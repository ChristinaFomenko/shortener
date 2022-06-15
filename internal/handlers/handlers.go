package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go

type Request struct {
	URL string `json:"url"`
}
type Result struct {
	Result string `json:"result"`
}

type service interface {
	Shorten(url string) string
	Expand(id string) (string, error)
	APIShortener(url string) (string, error)
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
	}

	shortcut := h.service.Shorten(string(bytes))

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

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func (h *handler) APIJSONShortener(w http.ResponseWriter, r *http.Request) {
	var sr ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.APIShortener(sr.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.json(w, http.StatusCreated, &ShortenResponse{Result: shortURL})
}

func (h *handler) json(w http.ResponseWriter, statusCode int, data interface{}) {
	resp, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(resp)
}
