package db

import (
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var Database *gorm.DB

// InitDB инициализирует глобальную переменную Database.
func InitDB() {
	var err error
	Database, err = ConnectDB(models.DatabaseURI)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	log.Println("Подключение к базе данных успешно установлено.")
}

// ConnectDB инициализирует подключение к базе данных.
func ConnectDB(databaseURI string) (*gorm.DB, error) {
	// Подключение к базе данных
	database, err := gorm.Open(postgres.Open(databaseURI), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автоматическая миграция таблиц
	if err := database.AutoMigrate(&models.User{}, &models.Order{}, &models.Balance{}); err != nil {
		return nil, err
	}

	return database, nil
}

// CloseDB закрывает соединение с базой данных.
func CloseDB() {
	sqlDB, err := Database.DB()
	if err != nil {
		log.Printf("Ошибка получения необработанного подключения к базе данных: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Ошибка при закрытии базы данных: %v", err)
	}
}
