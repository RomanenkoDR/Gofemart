package main

import (
	"github.com/RomanenkoDR/Gofemart/internal/config"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/handler"
	"github.com/RomanenkoDR/Gofemart/internal/router"
	"log"
	"net/http"
)

func main() {
	config.Init()
	// Инициализация базы данных
	database, err := db.ConnectDB(config.DatabaseURI)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.CloseDB(database)

	// Инициализация обработчиков
	h := handler.NewHandler(database)

	// Инициализация маршрутов
	r := router.SetupRouter(h)

	// Запуск HTTP-сервера
	log.Printf("Сервер запущен на %s", config.RunAddress)
	if err := http.ListenAndServe(config.RunAddress, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
