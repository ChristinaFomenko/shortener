package main

import (
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/app/generator"
	"github.com/ChristinaFomenko/shortener/internal/app/hasher"
	repositoryURL "github.com/ChristinaFomenko/shortener/internal/app/repository/urls"
	authService "github.com/ChristinaFomenko/shortener/internal/app/service/auth"
	serviceURL "github.com/ChristinaFomenko/shortener/internal/app/service/urls"
	"github.com/ChristinaFomenko/shortener/internal/handlers"
	"github.com/ChristinaFomenko/shortener/internal/middlewares"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// Config
	cfg, err := configs.NewConfig()
	if err != nil {
		err = env.Parse(cfg)
		log.Fatalf("failed to retrieve env variables, %v", err)
	}

	// Repositories
	repository, err := repositoryURL.NewStorage(cfg.FileStoragePath, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to create a storage %v", err)
	}

	// Services
	helper := generator.NewGenerator()
	hash := hasher.NewHasher(cfg.SecretKey)
	service := serviceURL.NewService(repository, helper, cfg.BaseURL)
	authSrvc := authService.NewService(helper, hash)
	//pingSrvc := pingService.NewService(repository)

	// Route
	router := chi.NewRouter()

	compress, err := middlewares.NewCompressor()
	if err != nil {
		log.Fatalf("compressor failed %v", err)
	}

	auth := middlewares.NewAuthenticator(authSrvc)

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middlewares.Decompressing)
	router.Use(compress.Compressing)
	router.Use(auth.Auth)

	//router.Route("/", func(r chi.Router) {
	router.Post("/", handlers.New(service, auth).Shorten)
	router.Get("/{id}", handlers.New(service, auth).Expand)
	router.Post("/api/shorten", handlers.New(service, auth).APIJSONShorten)
	router.Get("/api/user/urls", handlers.New(service, auth).FetchURLs)
	router.Get("/ping", handlers.New(service, auth).Ping)
	router.Post("/api/shorten/batch", handlers.New(service, auth).ShortenBatch)
	//})

	address := cfg.ServerAddress
	log.WithField("address", address).Info("server starts")
	log.Fatal(http.ListenAndServe(address, router), nil)
}
