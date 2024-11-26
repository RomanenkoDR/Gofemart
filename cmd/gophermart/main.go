package main

import (
	"github.com/joho/godotenv" // Импортируем библиотеку godotenv
	"log"
	"net/http"

	"github.com/RomanenkoDR/Gofemart/internal/config"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/handler"
	"github.com/RomanenkoDR/Gofemart/internal/router"
)

func init() {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env, используются системные переменные окружения")
	}
}

func main() {
	// Загрузка конфигурации сервера
	serverConfig := config.LoadServerConfig()

	// Загрузка конфигурации базы данных
	dbConfig := db.LoadDatabaseConfig()

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
