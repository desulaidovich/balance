package config

import (
	"github.com/gookit/ini/v2/dotenv"
)

type Data struct {
	Name             string
	Port             string
	ConnectionString string
}

func LoadEnvFromFile() (*Data, error) {
	if err := dotenv.Load("./config", ".env"); err != nil {
		return nil, err
	}

	config := new(Data)
	config.Name = dotenv.Get("APP_NAME")
	config.Port = dotenv.Get("HTTP_PORT")
	config.ConnectionString = dotenv.Get("DATABASE_CONNECTION_STRING")

	return config, nil
}
