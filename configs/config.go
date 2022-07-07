package configs

type AppConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./storage.json"`
	AuthKey         string `env:"AUTH_KEY" envDefault:"auth"`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:""`
}

//func (a *AppConfig) NewConfig() error {
//	serverAddress := getServerAddress()
//	baseURL := getBaseURL()
//	fileStoragePath := getFileStoragePath()
//	databaseDSN := getDatabaseDSN()
//	flag.Parse()
//
//	if serverAddress == nil {
//		return errors.New("server address not specified")
//	}
//
//	if baseURL == nil {
//		return errors.New("base url not specified")
//	}
//
//	if fileStoragePath == nil {
//		return errors.New("file storage path not specified")
//	}
//
//	if databaseDSN == nil {
//		return errors.New("database dsn not specified")
//	}
//
//	return nil
//}
//
//func getServerAddress() *string {
//	address := os.Getenv("SERVER_ADDRESS")
//	if address == "" {
//		address = ":8080"
//	}
//
//	return flag.String("a", address, "server address")
//}
//
//func getBaseURL() *string {
//	url := os.Getenv("BASE_URL")
//	if url == "" {
//		url = "http://localhost:8080"
//	}
//
//	return flag.String("b", url, "base url")
//}
//
//func getFileStoragePath() *string {
//	path := os.Getenv("FILE_STORAGE_PATH")
//
//	return flag.String("f", path, "file storage path")
//}
//
//func getDatabaseDSN() *string {
//	databaseDSN := os.Getenv("DATABASE_DSN")
//
//	return flag.String("d", databaseDSN, "database")
//}
