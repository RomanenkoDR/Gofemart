package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/internal/config"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {

	// Декодируем запрос в глобальную переменную User
	err := json.NewDecoder(r.Body).Decode(&models.User)
	if err != nil || models.User.Login == "" || models.User.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Получаем пользователя из базы
	err = db.GetUserByLogin(models.User.Login)
	if err != nil {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	// Проверяем пароль
	if !models.CheckPassword(models.User.Password) {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
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
