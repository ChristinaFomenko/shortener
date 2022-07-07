package main

import (
	"flag"
	"fmt"
	"github.com/ChristinaFomenko/shortener/configs"
	"github.com/ChristinaFomenko/shortener/internal/app/storage"
	"github.com/ChristinaFomenko/shortener/internal/router"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := configs.AppConfig{}
	if err := env.Parse(&cfg); err != nil {
		return
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "127.0.0.1:8080")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "http://localhost:8080")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "./storage.json")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "postgres://username:password@host:port/database")
	flag.Parse()

	s := storage.ConstructStorage(cfg)

	r := router.Router(cfg, s)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		if err := s.DestructStorage(cfg); err != nil {
			fmt.Printf("ERROR: %s", err)
		}
		os.Exit(0)
	}()

	address := cfg.ServerAddress
	log.WithField("address", address).Info("server starts")
	log.Fatal(http.ListenAndServe(address, r), nil)
}
