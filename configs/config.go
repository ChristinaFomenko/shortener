package configs

//func ServerAddress() string {
//	address := os.Getenv("SERVER_ADDRESS")
//	if address == "" {
//		address = "8080"
//	}
//
//	return ":" + address
//}
//
//func BaseURL() string {
//	url := os.Getenv("BASE_URL")
//	if url == "" {
//		url = "http://localhost:8080"
//	}
//
//	return url
//}

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"http://localhost:8080"`
	BaseUrl       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}
