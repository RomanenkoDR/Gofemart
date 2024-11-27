package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/handler"
	"github.com/RomanenkoDR/Gofemart/internal/router"
)

var (
	runAddress       string
	databaseURI      string
	accrualSystemURI string
)

func init() {
	// Загружаем переменные окружения из файла .env
	//if err := godotenv.Load(); err != nil {
	//	log.Println("Не удалось загрузить файл .env, используются системные переменные окружения")
	//}

	envRunAddress := os.Getenv("RUN_ADDRESS")
	envDatabaseURI := os.Getenv("DATABASE_URI")
	envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")

	// Добавляем флаги для конфигурирования
	flag.StringVar(&runAddress, "a", "localhost:8080", "Адрес и порт сервиса (например: localhost:8080)")
	flag.StringVar(&databaseURI, "d", "postgres://postgres:password@localhost:5432/gofemart", "URI подключения к базе данных")
	flag.StringVar(&accrualSystemURI, "r", "http://localhost:8081", "Адрес системы расчёта начислений")

	// Парсим флаги
	flag.Parse()

	// Если переменные окружения заданы, они переопределяют значения по умолчанию
	if envRunAddress != "" {
		runAddress = envRunAddress
	}
	if envDatabaseURI != "" {
		databaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress != "" {
		accrualSystemURI = envAccrualSystemAddress
	}
}

func main() {
	// Инициализация базы данных
	database, err := db.ConnectDB(databaseURI)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.CloseDB(database)

	// Инициализация обработчиков
	h := handler.NewHandler(database)

	// Инициализация маршрутов
	r := router.SetupRouter(h)

	// Запуск HTTP-сервера
	log.Printf("Сервер запущен на %s", runAddress)
	if err := http.ListenAndServe(runAddress, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
