package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"github.com/RomanenkoDR/Gofemart/internal/services"
	"net/http"
	"time"
)

// Login обрабатывает аутентификацию пользователя.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.User

	// Декодируем данные из запроса
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Login == "" || input.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Получаем пользователя из базы
	storedUser := &models.User{}
	if err := db.GetUserByLogin(h.DB, input.Login, storedUser); err != nil {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	// Проверяем пароль
	if !services.CheckPassword(input.Password, storedUser.Password) {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	// Генерация JWT токена
	jwtToken, err := services.GenerateJWT(storedUser.Login)
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

	w.Header().Set("Authorization", "Bearer "+jwtToken)
	w.WriteHeader(http.StatusOK)
}
