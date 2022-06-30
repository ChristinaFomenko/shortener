package configs

import (
	"errors"
	"flag"
	"os"
)

type appConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./storage.json"`
	AuthKey         string `env:"AUTH_KEY" envDefault:"auth"`
}

func NewConfig() (*appConfig, error) {
	serverAddress := getServerAddress()
	baseURL := getBaseURL()
	fileStoragePath := getFileStoragePath()
	flag.Parse()

	if serverAddress == nil {
		return nil, errors.New("server address not specified")
	}

	if baseURL == nil {
		return nil, errors.New("base url not specified")
	}

	if fileStoragePath == nil {
		return nil, errors.New("file storage path not specified")
	}

	return &appConfig{
		ServerAddress:   *serverAddress,
		BaseURL:         *baseURL,
		FileStoragePath: *fileStoragePath,
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
