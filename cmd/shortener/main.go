package main

import (
	"github.com/ChristinaFomenko/URLShortener/internal/app/generator"
	repositoryURL "github.com/ChristinaFomenko/URLShortener/internal/app/repository/urls"
	serviceURL "github.com/ChristinaFomenko/URLShortener/internal/app/service/urls"
	"github.com/ChristinaFomenko/URLShortener/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
)

//func init() {
//	// loads values from .env into the system
//	if err := godotenv.Load(); err != nil {
//		log.Print("No .env file found")
//	}
//}
//
//type Config struct {
//	ServerAddress string `env:"SERVER_ADDRESS"`
//	BaseURL       string `env:"BASE_URL"`
//}

func main() {

	serverAddress := os.Getenv("SERVER_ADDRESS")
	baseURL := os.Getenv("BASE_URL")
	//var cfg Config
	//err := env.Parse(&cfg)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Repositories
	repository := repositoryURL.NewRepo()

	// Services
	helper := generator.NewGenerator()
	service := serviceURL.NewService(repository, helper, baseURL)

	// Route
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//router.Route("/", func(r chi.Router) {
	router.Post("/", handlers.New(service).Shorten)
	router.Get("/{id}", handlers.New(service).Expand)
	router.Post("/api/shorten", handlers.New(service).APIJSONShortener)
	//})
	//port := configs.HTTPPort()

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
