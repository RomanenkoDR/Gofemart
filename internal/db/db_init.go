package db

import (
	"fmt"
	"log"
	"os"

	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// LoadDatabaseConfig загружает конфигурацию базы данных из флагов и переменных окружения.
func LoadDatabaseConfig(databaseURI string) models.DatabaseConfig {
	// Если флаг для URI базы данных не задан, читаем переменную окружения
	if databaseURI == "" {
		databaseURI = os.Getenv("DATABASE_URI")
	}

	// Формируем строку подключения
	return models.DatabaseConfig{
		Host:     "localhost", // Можно добавить логику для парсинга URI
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		Name:     "gofemart",
		SSLMode:  "disable",
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
