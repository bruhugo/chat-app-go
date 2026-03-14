package config

import (
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
)

var EnvConfig *ConfigType

func LoadConfig() (*ConfigType, error) {

	secretType := os.Getenv("SECRETS")
	if secretType == "" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
		secretType = os.Getenv("SECRETS")
	}

	switch secretType {
	case "LOCAL":
		return LoadLocalConfig()
	case "AWS":
		return LoadAWSSecrets()
	default:
		panic("Unknown secret source.")
	}

}

func LoadAWSSecrets() (*ConfigType, error) {
	EnvConfig = &ConfigType{}
	err := json.Unmarshal([]byte(os.Getenv("DB_SECRETS")), EnvConfig)
	if err != nil {
		return nil, err
	}

	EnvConfig.DbHost = os.Getenv("db_host")
	EnvConfig.DbDatabase = os.Getenv("db_database")
	EnvConfig.DbPort = os.Getenv("db_port")
	EnvConfig.Port = os.Getenv("port")
	EnvConfig.FrontendHost = os.Getenv("frontend_host")

	return EnvConfig, nil
}

func LoadLocalConfig() (*ConfigType, error) {
	EnvConfig = &ConfigType{
		DbUser:       os.Getenv("db_user"),
		DbPassword:   os.Getenv("db_password"),
		DbHost:       os.Getenv("db_host"),
		DbPort:       os.Getenv("db_port"),
		DbDatabase:   os.Getenv("db_database"),
		JwtSecret:    os.Getenv("jwt_secret"),
		Port:         os.Getenv("port"),
		FrontendHost: os.Getenv("frontend_host"),
	}

	return EnvConfig, nil
}
