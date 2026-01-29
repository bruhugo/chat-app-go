package config

type ConfigType struct {
	DbUser     string
	DbPassword string
	DbPort     string
	DbHost     string
	DbDatabase string

	JwtSecret string

	Port string
}
