package main

import (
	"flag"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/app/generator"
	repositoryURL "github.com/ChristinaFomenko/shortener/internal/app/repository/urls"
	serviceURL "github.com/ChristinaFomenko/shortener/internal/app/service/urls"
	"github.com/ChristinaFomenko/shortener/internal/handlers"
	"github.com/ChristinaFomenko/shortener/internal/middlewares"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

var Cfg configs.Config

type serviceParams struct {
	ServerAddress   *int
	BaseURL         *string
	FileStoragePath *string
}

func main() {
	serviceParamsObj := &serviceParams{}

	serviceParamsObj.ServerAddress = flag.Int("a", 8080, "port")
	//serviceParamsObj.BaseURL = flag.String("b", "http://localhost:8080", "base url")
	//serviceParamsObj.FileStoragePath = flag.String("f", "http://localhost:8080", "file path")
	flag.Parse()

	if err := env.Parse(&Cfg); err != nil {
		log.Fatal("failed parse configs:", err)
	}
	// Repositories
	repository := repositoryURL.Storage(configs.FileStoragePath())

	// Services
	helper := generator.NewGenerator()
	service := serviceURL.NewService(repository, helper, Cfg.BaseURL)

	// Route
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middlewares.Compressing)
	router.Use(middlewares.Decompressing)

	//router.Route("/", func(r chi.Router) {
	router.Post("/", handlers.New(service).Shorten)
	router.Get("/{id}", handlers.New(service).Expand)
	router.Post("/api/shorten", handlers.New(service).APIJSONShorten)
	//})

	//address := configs.ServerAddress()
	//log.WithField("address", address).Info("server starts")
	u, err := url.Parse(Cfg.BaseURL)
	if err != nil {
		log.Error(err)
	}
	log.WithField("address", *serviceParamsObj.ServerAddress).Info("server starts")
	log.Fatal(http.ListenAndServe(u.Host, router), nil)
}
