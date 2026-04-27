package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvData struct {
	PORT                              string
	DATABASE_URL                      string
	CONTEXT_TIMEOUT                   string
	SERVER_ADDRESS                    string
	ENVIRONMENT                       string
	GOOGLE_CLIENT_ID                  string
	GOOGLE_CLIENT_SECRET              string
	GOOGLE_REDIRECT_URL               string
	GITHUB_CLIENT_ID                  string
	GITHUB_CLIENT_SECRET              string
	GITHUB_REDIRECT_URL               string
	PRODUCTION                        bool
	CLIENT_DASHBOARD                  string
	STORAGE_SERVICE_ACCESS_KEY_ID     string
	STORAGE_SERVICE_SECRET_ACCESS_KEY string
	BUCKET_NAME                       string
	STORAGE_SERVICE_ENDPOINT          string
	BUCKET_ROOT                       string
	DOMAIN                            string
	DATA_GET_PATH                     string
	REDIS_ADDR                        string
	REDIS_PASSWORD                    string
}

func LoadEnv() *EnvData {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env file not found")
	}

	config := EnvData{
		DATABASE_URL:                      os.Getenv("DATABASE_URL"),
		PORT:                              os.Getenv("PORT"),
		CONTEXT_TIMEOUT:                   os.Getenv("CONTEXT_TIMEOUT"),
		SERVER_ADDRESS:                    os.Getenv("SERVER_ADDRESS"),
		ENVIRONMENT:                       os.Getenv("ENVIRONMENT"),
		GOOGLE_CLIENT_ID:                  os.Getenv("GOOGLE_CLIENT_ID"),
		GOOGLE_CLIENT_SECRET:              os.Getenv("GOOGLE_CLIENT_SECRET"),
		GOOGLE_REDIRECT_URL:               os.Getenv("GOOGLE_REDIRECT_URL"),
		GITHUB_CLIENT_ID:                  os.Getenv("GITHUB_CLIENT_ID"),
		GITHUB_CLIENT_SECRET:              os.Getenv("GITHUB_CLIENT_SECRET"),
		GITHUB_REDIRECT_URL:               os.Getenv("GITHUB_REDIRECT_URL"),
		CLIENT_DASHBOARD:                  os.Getenv("CLIENT_DASHBOARD"),
		STORAGE_SERVICE_ACCESS_KEY_ID:     os.Getenv("STORAGE_SERVICE_ACCESS_KEY_ID"),
		STORAGE_SERVICE_SECRET_ACCESS_KEY: os.Getenv("STORAGE_SERVICE_SECRET_ACCESS_KEY"),
		BUCKET_NAME:                       os.Getenv("BUCKET_NAME"),
		STORAGE_SERVICE_ENDPOINT:          os.Getenv("STORAGE_SERVICE_ENDPOINT"),
		BUCKET_ROOT:                       os.Getenv("BUCKET_ROOT"),
		DOMAIN:                            os.Getenv("DOMAIN"),
		DATA_GET_PATH:                     os.Getenv("DATA_GET_PATH"),
		REDIS_ADDR:                        os.Getenv("REDIS_ADDR"),
		REDIS_PASSWORD:                    os.Getenv("REDIS_PASSWORD"),
	}
	config.PRODUCTION = config.ENVIRONMENT == "production"
	log.Println("Enviroments variables loaded")
	return &config
}
