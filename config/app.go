package config

import "os"

type AppConfig struct {
	Name    string
	Version string
	Debug   string
	Host    string
	Port    string
	Key     string
	Locale  string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Name:    os.Getenv("APP_NAME"),
		Version: os.Getenv("APP_VERSION"),
		Debug:   os.Getenv("APP_DEBUG"),
		Host:    os.Getenv("APP_HOST"),
		Port:    os.Getenv("APP_PORT"),
		Key:     os.Getenv("APP_KEY"),
		Locale:  os.Getenv("APP_LOCALE"),
	}
}
