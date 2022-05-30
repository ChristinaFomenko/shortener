package main

import (
	"fmt"
	"github.com/ChristinaFomenko/URLShortener/configs"
	"github.com/ChristinaFomenko/URLShortener/internal/app/helpers"
	repositoryURL "github.com/ChristinaFomenko/URLShortener/internal/app/repository/urls"
	serviceURL "github.com/ChristinaFomenko/URLShortener/internal/app/service/urls"
	"github.com/ChristinaFomenko/URLShortener/internal/handlers"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.New(service).Shorten)
	mux.HandleFunc("/{id}", handlers.New(service).Expand)

	port := configs.HTTPPort()

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(port, mux))
}
