package config

import (
	"fmt"
	"os"
	"strconv"
)

// ServerConfig содержит конфигурацию сервера.
type ServerConfig struct {
	Host string
	Port int
}

// LoadServerConfig загружает конфигурацию сервера из переменных окружения.
func LoadServerConfig() ServerConfig {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil || port <= 0 {
		port = 8080 // Значение по умолчанию
	}

	return ServerConfig{
		Host: os.Getenv("SERVER_HOST"),
		Port: port,
	}
}

// Address возвращает полный адрес сервера.
func (sc ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}
