package configs

import (
	"flag"
	"os"
)

func init() {
	cfg.serverAddress = serverAddress()
	cfg.baseURL = baseURL()
	cfg.fileStoragePath = fileStoragePath()

	flag.Parse()
}

var cfg struct {
	serverAddress   *string
	baseURL         *string
	fileStoragePath *string
}

func ServerAddress() string {
	if cfg.serverAddress == nil {
		panic("server address not specified")
	}

	return *cfg.serverAddress
}

func BaseURL() string {
	if cfg.baseURL == nil {
		panic("base url not specified")
	}

	return *cfg.baseURL
}

func FileStoragePath() string {
	if cfg.fileStoragePath == nil {
		panic("file storage path not specified")
	}

	return *cfg.fileStoragePath
}

func serverAddress() *string {
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = ":8080"
	}

	return flag.String("a", address, "server address")
}

func baseURL() *string {
	url := os.Getenv("BASE_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	return flag.String("b", url, "base url")
}

func fileStoragePath() *string {
	path := os.Getenv("FILE_STORAGE_PATH")

	return flag.String("f", path, "file storage path")
}
