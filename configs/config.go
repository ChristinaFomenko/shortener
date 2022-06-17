package configs

import "os"

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"http://localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func FileStoragePath() string {
	return os.Getenv("FILE_STORAGE_PATH")
}
