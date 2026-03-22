package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvData struct {
	PORT           string
	DATABASEURL    string
	CONTEXTTIMEOUT string
	SERVERADDRESS  string
	ENVIROMENT     string
}

func LoadEnv() *EnvData {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env file not found")
	}

	config := EnvData{
		DATABASEURL:    os.Getenv("DATABASE_URL"),
		PORT:           os.Getenv("PORT"),
		CONTEXTTIMEOUT: os.Getenv("CONTEXT_TIMEOUT"),
		SERVERADDRESS:  os.Getenv("SERVER_ADDRESS"),
		ENVIROMENT:     os.Getenv("ENVIROMENT"),
	}
	log.Println("Enviroments variables loaded")
	return &config
}
