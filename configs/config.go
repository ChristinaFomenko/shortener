package configs

import (
	"errors"
	"flag"
	"github.com/caarlos0/env"
	"os"
)

type appConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./storage.json"`
	AuthKey         []byte
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func NewConfig() (*appConfig, error) {
	serverAddress := getServerAddress()
	baseURL := getBaseURL()
	fileStoragePath := getFileStoragePath()
	databaseDSN := getDatabaseDSN()
	flag.Parse()

	instance := &appConfig{}
	if err := env.Parse(instance); err != nil {
		return nil, errors.New("auth key not specified")
	}

	if serverAddress == nil {
		return nil, errors.New("server address not specified")
	}

	if baseURL == nil {
		return nil, errors.New("base url not specified")
	}

	if fileStoragePath == nil {
		return nil, errors.New("file storage path not specified")
	}

	if databaseDSN == nil {
		return nil, errors.New("database dsn not specified")
	}

	instance.AuthKey = make([]byte, 16)

	return &appConfig{
		ServerAddress:   *serverAddress,
		BaseURL:         *baseURL,
		FileStoragePath: *fileStoragePath,
		DatabaseDSN:     *databaseDSN,
	}, nil
}

func getServerAddress() *string {
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = ":8080"
	}

	return flag.String("a", address, "server address")
}

func getBaseURL() *string {
	url := os.Getenv("BASE_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	return flag.String("b", url, "base url")
}

func getFileStoragePath() *string {
	path := os.Getenv("FILE_STORAGE_PATH")

	return flag.String("f", path, "file storage path")
}

func getDatabaseDSN() *string {
	databaseDSN := os.Getenv("DATABASE_DSN")

	return flag.String("d", databaseDSN, "database")
}
