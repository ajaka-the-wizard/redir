package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvData struct {
	PORT                 string
	DATABASE_URL         string
	CONTEXT_TIMEOUT      string
	SERVER_ADDRESS       string
	ENVIRONMENT          string
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	GOOGLE_REDIRECT_URL  string
	GITHUB_CLIENT_ID     string
	GITHUB_CLIENT_SECRET string
	GITHUB_REDIRECT_URL  string
	PRODUCTION           bool
	CLIENT_DASHBOARD     string
}

func LoadEnv() *EnvData {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env file not found")
	}

	config := EnvData{
		DATABASE_URL:         os.Getenv("DATABASE_URL"),
		PORT:                 os.Getenv("PORT"),
		CONTEXT_TIMEOUT:      os.Getenv("CONTEXT_TIMEOUT"),
		SERVER_ADDRESS:       os.Getenv("SERVER_ADDRESS"),
		ENVIRONMENT:          os.Getenv("ENVIRONMENT"),
		GOOGLE_CLIENT_ID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GOOGLE_CLIENT_SECRET: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GOOGLE_REDIRECT_URL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		GITHUB_CLIENT_ID:     os.Getenv("GITHUB_CLIENT_ID"),
		GITHUB_CLIENT_SECRET: os.Getenv("GITHUB_CLIENT_SECRET"),
		GITHUB_REDIRECT_URL:  os.Getenv("GITHUB_REDIRECT_URL"),
		CLIENT_DASHBOARD:     os.Getenv("CLIENT_DASHBOARD"),
	}
	config.PRODUCTION = config.ENVIRONMENT == "production"
	log.Println("Enviroments variables loaded")
	return &config
}
