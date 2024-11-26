package models

// DatabaseConfig содержит конфигурацию базы данных.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}
