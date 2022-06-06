package handlers

import (
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go

type service interface {
	Shorten(url string) string
	Expand(id string) (string, error)
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
