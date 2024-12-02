package config

import (
	"flag"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/handler"
	"github.com/RomanenkoDR/Gofemart/internal/router"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	SecretKey            string
}

// InitServer Инициализация конфигурации сервера
func InitServer() {
	// Инициализация конфигурации
	cfg := initConfig()

	// Инициализация базы данных
	database, err := db.ConnectDB(cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.CloseDB(database)

	// Инициализация обработчиков
	h := handler.NewHandler(database, cfg.AccrualSystemAddress)

	// Инициализация маршрутов
	r := router.SetupRouter(h)

	// Запуск HTTP-сервера
	log.Printf("Сервер запущен на %s", cfg.RunAddress)
	if err := http.ListenAndServe(cfg.RunAddress, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

// Определение конфигурации сервера (с загрузкой переменных окружения и флагов)
func initConfig() Config {
	envPath := "../../.env"
	// Загружаем переменные из файла .env
	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Не удалось загрузить файл .env: %v", err)
	}

	// Значения по умолчанию
	defaultRunAddress := "localhost:8080"
	defaultDatabaseURI := "postgres://postgres:password@localhost:5432/gofemart"
	defaultAccrualSystemURI := "http://localhost:8081"

	// Считываем переменные окружения
	envRunAddress := os.Getenv("RUN_ADDRESS")
	envDatabaseURI := os.Getenv("DATABASE_URI")
	envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	envSecretKey := os.Getenv("SECRET_KEY")

	// Добавляем флаги
	flagRunAddress := flag.String("a", envOrDefault(envRunAddress, defaultRunAddress), "Адрес сервера")
	flagDatabaseURI := flag.String("d", envOrDefault(envDatabaseURI, defaultDatabaseURI), "URI подключения к базе данных")
	flagAccrualSystemAddress := flag.String("r", envOrDefault(envAccrualSystemAddress, defaultAccrualSystemURI), "Адрес системы начислений")
	flagSecretKey := flag.String("secret_key", envOrDefault(envSecretKey, ""), "Секретный ключ для токенов")

	// Парсинг флагов
	flag.Parse()

	// Проверяем обязательные значения
	if *flagSecretKey == "" {
		log.Fatalf("SECRET_KEY отсутствует в переменных окружения и не задан через флаг")
	}

	// Возвращаем актуальную конфигурацию после проверок
	return Config{
		RunAddress:           *flagRunAddress,
		DatabaseURI:          *flagDatabaseURI,
		AccrualSystemAddress: *flagAccrualSystemAddress,
		SecretKey:            *flagSecretKey,
	}
}

// Функция для обработки переменных окружения
func envOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
