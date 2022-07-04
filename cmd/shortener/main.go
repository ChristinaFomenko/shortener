package main

import (
	"database/sql"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/app/generator"
	repositoryURL "github.com/ChristinaFomenko/shortener/internal/app/repository/urls"
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

	// Database
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connnect db %v", err)
	}
	defer db.Close()

	// Repositories
	repository, err := repositoryURL.NewStorage(cfg.FileStoragePath)
	if err != nil {
		log.Fatalf("failed to create a storage %v", err)
	}
	// Services
	helper := generator.NewGenerator()
	service := serviceURL.NewService(repository, helper, cfg.BaseURL)

	// Route
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middlewares.GZIPMiddleware)
	router.Use(middlewares.AuthCookie)

	//router.Route("/", func(r chi.Router) {
	router.Post("/", handlers.New(service).Shorten)
	router.Get("/{id}", handlers.New(service).Expand)
	router.Post("/api/shorten", handlers.New(service).APIJSONShorten)
	router.Get("/api/user/urls", handlers.New(service).GetList)
	router.Get("/ping", handlers.New(service).Ping)
	//})

	address := cfg.ServerAddress
	log.WithField("address", address).Info("server starts")
	log.Fatal(http.ListenAndServe(address, router), nil)
}
