package configs

import (
	"errors"
	"flag"
	"os"
)

type appConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	SecretKey       []byte
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:""`
}

func NewConfig() (*appConfig, error) {
	serverAddress := getServerAddress()
	baseURL := getBaseURL()
	fileStoragePath := getFileStoragePath()
	secretKey := getSecretKey()
	databaseDSN := getDatabaseDSN()
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

	if secretKey == nil {
		return nil, errors.New("secret key not specified")
	}
	if databaseDSN == nil {
		return nil, errors.New("database dsn not specified")
	}

	return &appConfig{
		ServerAddress:   *serverAddress,
		BaseURL:         *baseURL,
		FileStoragePath: *fileStoragePath,
		SecretKey:       []byte(*secretKey),
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

func getSecretKey() *string {
	url := os.Getenv("SECRET_KEY")
	if url == "" {
		url = "my-secret-key"
	}

	return flag.String("s", url, "secret key")
}

func getDatabaseDSN() *string {
	databaseDSN := os.Getenv("DATABASE_DSN")

	return flag.String("d", databaseDSN, "database")
}
