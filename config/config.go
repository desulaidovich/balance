package config

import (
	"github.com/gookit/ini/v2/dotenv"
)

type ConfigData struct {
	Name             string
	Port             string
	ConnectionString string
}

func Load() *ConfigData {
	if err := dotenv.Load("./config", ".env"); err != nil {
		panic(err)
	}

	return &ConfigData{
		Name:             dotenv.Get("APP_NAME"),
		Port:             dotenv.Get("HTTP_PORT"),
		ConnectionString: dotenv.Get("DATABASE_CONNECTION_STRING"),
	}
}
