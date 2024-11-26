package db

import (
	"fmt"
	"log"
	"os"

	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// LoadDatabaseConfig загружает конфигурацию базы данных из переменных окружения.
func LoadDatabaseConfig() models.DatabaseConfig {
	return models.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  "disable", // Значение по умолчанию для локальной разработки
	}
}

// ConnectDB инициализирует подключение к базе данных.
func ConnectDB(cfg models.DatabaseConfig) (*gorm.DB, error) {
	// Формируем строку подключения
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	// Подключение к базе данных
	database, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	// Автоматическая миграция таблиц
	if err := database.AutoMigrate(&models.User{}, &models.Order{}, &models.Balance{}); err != nil {
		return nil, fmt.Errorf("не удалось выполнить миграции: %w", err)
	}

	log.Println("Подключение к базе данных успешно установлено.")
	return database, nil
}

// CloseDB закрывает соединение с базой данных.
func CloseDB(database *gorm.DB) {
	sqlDB, err := database.DB()
	if err != nil {
		log.Printf("Ошибка получения необработанного подключения к базе данных: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Ошибка при закрытии базы данных: %v", err)
	}
}
