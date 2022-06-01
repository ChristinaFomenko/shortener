package main

import (
	"fmt"
	"github.com/ChristinaFomenko/URLShortener/configs"
	"github.com/ChristinaFomenko/URLShortener/internal/app/helpers"
	repositoryURL "github.com/ChristinaFomenko/URLShortener/internal/app/repository/urls"
	serviceURL "github.com/ChristinaFomenko/URLShortener/internal/app/service/urls"
	"github.com/ChristinaFomenko/URLShortener/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	// Repositories
	repository := repositoryURL.NewRepo()

	// Services
	helper := helpers.NewGenerator()
	service := serviceURL.NewService(repository, helper, configs.HTTPHost())

	// Route
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//router.Route("/", func(r chi.Router) {
	router.Post("/", handlers.New(service).Shorten)
	router.Get("/{id}", handlers.New(service).Expand)
	//})
	port := configs.HTTPPort()

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(port, router))
}
