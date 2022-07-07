package main

import (
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

var conf configs.AppConfig

func main() {
	// Config
	cfg, err := configs.NewConfig()
	if err = env.Parse(cfg); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}

	s := storage.ConstructStorage(conf)

	r := router.Router(conf, s)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		if err := storage.DestructStorage(cfg.FileStoragePath, s); err != nil {
			fmt.Printf("ERROR: %s", err)
		}
		os.Exit(0)
	}()

	address := cfg.ServerAddress
	log.WithField("address", address).Info("server starts")
	log.Fatal(http.ListenAndServe(address, r), nil)
}
