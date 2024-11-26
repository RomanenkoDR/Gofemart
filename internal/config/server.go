package config

import (
	"fmt"
	"os"
)

// ServerConfig содержит конфигурацию сервера.
type ServerConfig struct {
	Host string
	Port int
}

// LoadServerConfig загружает конфигурацию сервера из флагов и переменных окружения.
func LoadServerConfig(runAddress string) ServerConfig {
	// Если флаг для адреса и порта сервера не задан, читаем переменные окружения
	if runAddress == "" {
		runAddress = os.Getenv("RUN_ADDRESS")
	}

	// Разделяем runAddress на хост и порт
	host, port := parseAddress(runAddress)

	return ServerConfig{
		Host: host,
		Port: port,
	}
}

// parseAddress разбивает строку вида "localhost:8080" на хост и порт
func parseAddress(address string) (string, int) {
	var host string
	var port int
	_, err := fmt.Sscanf(address, "%s:%d", &host, &port)
	if err != nil {
		port = 8080 // Значение по умолчанию
		host = "localhost"
	}
	return host, port
}

// Address возвращает полный адрес сервера.
func (sc ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}
