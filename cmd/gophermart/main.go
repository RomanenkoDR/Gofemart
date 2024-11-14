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

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	h := handler.NewHandler()
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	cfg := config.Options{
		Port: port,
		Host: host,
	}

	router, err := router.InitRouter(cfg, h)
	if err != nil {
		log.Fatalf("Error initializing router: %v", err)
	}
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	db, err := db.InitDB() // Инициализация соединения с базой данных
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}

	// Передаем db в вашу глобальную переменную или напрямую в обработчики
	handler.SetDatabase(db)

	log.Printf("Listening on %s", addr)

	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer db.Close()

}
