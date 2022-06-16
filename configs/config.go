package configs

import "os"

func ServerAddress() string {
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = "8080"
	}

	return ":" + address
}

func BaseURL() string {
	host := os.Getenv("BASE_URL")
	if host == "" {
		host = "http://localhost:8080"
	}

	return host
}
