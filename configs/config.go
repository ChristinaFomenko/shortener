package configs

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"os"
)

const AuthKeyLength = 32

type AppConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./storage.json"`
	AuthKey         string `env:"AUTH_KEY" envDefault:"auth"`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:"postgres://christina:123@postgres:5432/praktikum"`
}

func NewConfig() (*AppConfig, error) {
	serverAddress := getServerAddress()
	baseURL := getBaseURL()
	fileStoragePath := getFileStoragePath()
	authKey, _ := configureSecretKey()
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

	if databaseDSN == nil {
		return nil, errors.New("database dsn not specified")
	}

	return &AppConfig{
		ServerAddress:   *serverAddress,
		BaseURL:         *baseURL,
		FileStoragePath: *fileStoragePath,
		AuthKey:         string(authKey),
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

func configureSecretKey() ([]byte, error) {
	authKey := os.Getenv("BASE_URL")
	if authKey != "" {
		confKey, err := hex.DecodeString(authKey)
		if err != nil {
			return nil, err
		}
		return confKey, nil
	}
	return GenerateSecretKey(AuthKeyLength)
}

func GenerateSecretKey(length int) ([]byte, error) {
	randKey := make([]byte, length)
	_, err := rand.Read(randKey)
	if err != nil {
		return nil, err
	}
	return randKey, nil
}
func getDatabaseDSN() *string {
	databaseDSN := os.Getenv("DATABASE_DSN")

	return flag.String("d", databaseDSN, "database")
}
