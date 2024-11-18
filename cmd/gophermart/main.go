package main

import (
	"fmt"
	"github.com/RomanenkoDR/Gofemart/iternal/config"
	"github.com/RomanenkoDR/Gofemart/iternal/db"
	"github.com/RomanenkoDR/Gofemart/iternal/handler"
	"github.com/RomanenkoDR/Gofemart/iternal/router"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {

	// Загрузка параметров системы референции
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Инициализация DB
	database, err := db.InitDB() // Инициализация соединения с базой данных
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}
	defer func() {
		sqlDB, err := database.DB() // Извлечение базового *sql.DB из Gorm
		if err != nil {
			log.Printf("Error retrieving raw DB instance: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	handler.SetDatabase(database)

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	cfg := config.Options{
		Port: port,
		Host: host,
	}

	// Инициализация hanlder
	h := handler.NewHandler()

	// Инициализация маршрутов
	router, err := router.InitRouter(cfg, h)
	if err != nil {
		log.Fatalf("Error initializing router: %v", err)
	}

	// Выводим адрес запуска сервера
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	// Запуск HTTP-сервера
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	//defer database.Close()
}
