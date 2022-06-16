package handlers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
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

func (h *handler) APIJSONShorten(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	req := ShortenRequest{}
	if err = json.Unmarshal(b, &req); err != nil {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(req)
	if err != nil || !ok {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	shortcut := h.service.Shorten(req.URL)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := ShortenReply{ShortenUrlResult: shortcut}
	marshal, err := json.Marshal(&resp)
	if err != nil {
		log.WithError(err).WithField("resp", resp).Error("marshal response error")
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		log.WithError(err).WithField("shortcut", shortcut).Error("write response error")
		return
	}
}
