package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

func New() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		App:      LoadAppConfig(),
		Database: LoadDatabaseConfig(),
	}
}
