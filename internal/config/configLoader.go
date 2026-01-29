package config

import (
	"os"

	"github.com/joho/godotenv"
)

var EnvConfig *ConfigType

func LoadConfig() (*ConfigType, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	EnvConfig = &ConfigType{
		DbUser:     os.Getenv("db_user"),
		DbPassword: os.Getenv("db_password"),
		DbHost:     os.Getenv("db_host"),
		DbPort:     os.Getenv("db_port"),
		DbDatabase: os.Getenv("db_database"),
		JwtSecret:  os.Getenv("jwt_secret"),
		Port:       os.Getenv("port"),
	}

	return EnvConfig, nil
}
