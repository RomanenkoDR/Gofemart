package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/internal/config"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"net/http"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {

	// Декодируем запрос в глобальную переменную User
	err := json.NewDecoder(r.Body).Decode(&models.User)
	if err != nil || models.User.Login == "" || models.User.Password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	// Проверяем существование пользователя
	exists, err := db.CheckUserExists(models.User.Login)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Хэшируем пароль
	if err := models.HashPassword(); err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Сохраняем пользователя
	if err := db.CreateUser(); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Создаем баланс
	models.Balance.UserID = models.User.ID
	if err := db.CreateBalance(); err != nil {
		http.Error(w, "Failed to create balance record", http.StatusInternalServerError)
		return
	}

	// Генерация JWT токена
	models.JwtKey, err = config.GenerateJWT(models.User.Login)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Устанавливаем токен в cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   models.JwtKey,
		Expires: time.Now().Add(24 * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
}
