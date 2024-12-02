package config

import (
	"flag"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/handler"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"github.com/RomanenkoDR/Gofemart/internal/router"
	"log"
	"net/http"
	"os"
)

func InitServer() {
	// Чтение переменных окружения
	envRunAddress := os.Getenv("RUN_ADDRESS")
	envDatabaseURI := os.Getenv("DATABASE_URI")
	envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")

	// Добавляем флаги для конфигурирования
	flag.StringVar(&models.RunAddress, "a", "localhost:8080", "Адрес и порт сервиса (например: localhost:8080)")
	flag.StringVar(&models.DatabaseURI, "d", "postgres://postgres:password@localhost:5432/gofemart", "URI подключения к базе данных")
	flag.StringVar(&models.AccrualSystemURI, "r", "localhost:8081", "Адрес системы расчёта начислений")
	flag.Parse()

	// Если переменные окружения заданы, они переопределяют значения по умолчанию
	if envRunAddress != "" {
		models.RunAddress = envRunAddress
	}
	if envDatabaseURI != "" {
		models.DatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress != "" {
		models.AccrualSystemURI = envAccrualSystemAddress
	}

	// Инициализация базы данных
	db.InitDB()

	// Инициализация обработчиков
	h := handler.NewHandler(db.Database)

	// Инициализация маршрутов
	r := router.SetupRouter(h)

	// Запуск HTTP-сервера
	log.Printf("Сервер запущен на %s", models.RunAddress)
	if err := http.ListenAndServe(models.RunAddress, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
