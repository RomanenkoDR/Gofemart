package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"github.com/RomanenkoDR/Gofemart/internal/services"
)

// Register обрабатывает регистрацию нового пользователя.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Проверяем Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content-Type, expected application/json", http.StatusBadRequest)
		return
	}

	// Декодируем данные из запроса
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil || user.Login == "" || user.Password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	// Проверяем существование пользователя
	exists, err := db.CheckUserExists(h.DB, user.Login)

	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Хэшируем пароль
	hashedPassword, err := services.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// Создаем пользователя с балансом
	if err := db.CreateUserWithBalance(h.DB, &user); err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Генерация JWT токена
	jwtToken, err := services.GenerateJWT(user.Login)
	if err != nil {
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}

	// Устанавливаем токен в cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   jwtToken,
		Expires: time.Now().Add(24 * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
}
