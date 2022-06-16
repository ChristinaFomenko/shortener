package main

import (
	"fmt"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/app/generator"
	repositoryURL "github.com/ChristinaFomenko/shortener/internal/app/repository/urls"
	serviceURL "github.com/ChristinaFomenko/shortener/internal/app/service/urls"
	"github.com/ChristinaFomenko/shortener/internal/handlers"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

var Cfg configs.Config

func main() {
	if err := env.Parse(&Cfg); err != nil {
		fmt.Println("failed:", err)
	}
	// Repositories
	repository := repositoryURL.NewRepo()

	// Services
	helper := generator.NewGenerator()
	service := serviceURL.NewService(repository, helper, Cfg.BaseUrl)

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

	//address := configs.ServerAddress()
	//log.WithField("address", address).Info("server starts")
	u, err := url.Parse(Cfg.BaseUrl)
	if err != nil {
		log.Error(err)
	}
	log.WithField("address", Cfg.ServerAddress).Info("server starts")
	log.Fatal(http.ListenAndServe(u.Host, router), nil)
}
