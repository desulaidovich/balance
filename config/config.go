package config

import (
	"github.com/gookit/ini/v2/dotenv"
)

type ConfigData struct {
	Name             string
	Port             string
	ConnectionString string
}

func Load() (*ConfigData, error) {
	if err := dotenv.Load("./config", ".env"); err != nil {
		return nil, err
	}

	return &ConfigData{
		Name:             dotenv.Get("APP_NAME"),
		Port:             dotenv.Get("HTTP_PORT"),
		ConnectionString: dotenv.Get("DATABASE_CONNECTION_STRING"),
	}, nil
}
