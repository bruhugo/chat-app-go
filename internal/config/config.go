package config

type ConfigType struct {
	DbUser       string `json:"db_username"`
	DbPassword   string `json:"db_password"`
	DbPort       string `json:"db_port"`
	DbHost       string `json:"db_host"`
	DbDatabase   string `json:"db_database"`
	JwtSecret    string `json:"jwt_secret"`
	Port         string `json:"port"`
	FrontendHost string `json:"frontend_host"`
}
