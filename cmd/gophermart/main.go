package main

import (
	"flag"
	"github.com/joho/godotenv" // Импортируем библиотеку godotenv
	"log"
	"net/http"
	"os"

	"github.com/RomanenkoDR/Gofemart/internal/config"
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
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env, используются системные переменные окружения")
	}
	// Добавляем флаги для конфигурирования
	flag.StringVar(&runAddress, "a", os.Getenv("RUN_ADDRESS"), "Адрес и порт сервиса (например: localhost:8080)")
	flag.StringVar(&databaseURI, "d", os.Getenv("DATABASE_URI"), "URI подключения к базе данных")
	flag.StringVar(&accrualSystemURI, "r", os.Getenv("ACCRUAL_SYSTEM_ADDRESS"), "Адрес системы расчёта начислений")

	// Парсим флаги
	flag.Parse()

}

func main() {
	// Загрузка конфигурации сервера
	serverConfig := config.LoadServerConfig(runAddress)

	// Загрузка конфигурации базы данных
	dbConfig := db.LoadDatabaseConfig(databaseURI)

	// Инициализация базы данных
	database, err := db.ConnectDB(dbConfig)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.CloseDB(database)

	// Инициализация обработчиков
	h := handler.NewHandler(database)

	// Инициализация маршрутов
	r := router.SetupRouter(h)

	// Запуск HTTP-сервера
	address := serverConfig.Address()
	log.Printf("Сервер запущен на %s", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
