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

func (h *handler) APIShortener(w http.ResponseWriter, r *http.Request) {
	urlReq := Request{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&urlReq); err != nil {
		log.Printf("APIShortenURL: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}
	statusCode := http.StatusCreated
	shortURL := h.service.Shorten(urlReq.URL)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	enc := json.NewEncoder(w)
	var err error
	err = enc.Encode(&Result{Result: shortURL})
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Printf("APIShortHandler: %v", err)
	}
}
