package main

import (
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/app/generator"
	repositoryURL "github.com/ChristinaFomenko/shortener/internal/app/repository/urls"
	serviceURL "github.com/ChristinaFomenko/shortener/internal/app/service/urls"
	"github.com/ChristinaFomenko/shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// Repositories
	repository := repositoryURL.NewRepo()

	// Services
	helper := generator.NewGenerator()
	service := serviceURL.NewService(repository, helper, configs.BaseURL())

	// Route
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//router.Route("/", func(r chi.Router) {
	router.Post("/", handlers.New(service).Shorten)
	router.Get("/{id}", handlers.New(service).Expand)
	router.Post("/api/shorten", handlers.New(service).APIJSONShorten)
	//})
	address := configs.ServerAddress()
	log.WithField("address", address).Info("server starts")
	log.Fatal(http.ListenAndServe(address, router))
}
