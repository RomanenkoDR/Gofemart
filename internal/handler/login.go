package handler

import (
	"encoding/json"
	"errors"
	"github.com/RomanenkoDR/Gofemart/internal/config"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"github.com/RomanenkoDR/Gofemart/internal/services/db"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var (
		reqUser models.User
		user    models.User
		jwtKey  string
	)

	// Инициализация декодера
	err := json.NewDecoder(r.Body).Decode(&reqUser)

	// Валидация логина и пароля
	if err != nil || reqUser.Login == "" || reqUser.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result := db.Database.Where("login = ?", reqUser.Login).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Invalid login or password", http.StatusUnauthorized)
			return
		} else if result.Error != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Проверяем пароль
	if !user.CheckPassword(reqUser.Password) {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	// Генерация JWT токена
	jwtKey, err = config.GenerateJWT(user.Login)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: time.Now().Add(24 * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
}
