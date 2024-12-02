package db

import (
	"fmt"
	"log"

	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB инициализирует подключение к базе данных.
func ConnectDB(databaseURI string) (*gorm.DB, error) {
	// Проверяем наличие URI
	if databaseURI == "" {
		return nil, fmt.Errorf("строка подключения к базе данных пуста")
	}

	// Подключение к базе данных
	database, err := gorm.Open(postgres.Open(databaseURI), &gorm.Config{})
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
