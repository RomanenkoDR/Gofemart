package models

// ServerConfig содержит конфигурацию сервера.
type ServerConfig struct {
	Host string
	Port int
}

type ConfigFlag struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	TokenExpiredIn       string
	TokenMaxAge          int
	SecretKey            string
}
