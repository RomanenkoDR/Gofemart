package config

import (
	"flag"
	"os"
)

var (
	RunAddress       string
	DatabaseURI      string
	AccrualSystemURI string
)

func Init() {

	envRunAddress := os.Getenv("RUN_ADDRESS")
	envDatabaseURI := os.Getenv("DATABASE_URI")
	envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")

	// Добавляем флаги для конфигурирования
	flag.StringVar(&RunAddress, "a", "localhost:8080", "Адрес и порт сервиса (например: localhost:8080)")
	flag.StringVar(&DatabaseURI, "d", "postgres://postgres:password@localhost:5432/gofemart", "URI подключения к базе данных")
	flag.StringVar(&AccrualSystemURI, "r", "localhost:8081", "Адрес системы расчёта начислений")

	// Парсим флаги
	flag.Parse()

	// Если переменные окружения заданы, они переопределяют значения по умолчанию
	if envRunAddress != "" {
		RunAddress = envRunAddress
	}
	if envDatabaseURI != "" {
		DatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress != "" {
		AccrualSystemURI = envAccrualSystemAddress
	}
}
