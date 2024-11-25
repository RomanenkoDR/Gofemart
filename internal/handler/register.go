package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/iternal/config"
	"github.com/RomanenkoDR/Gofemart/iternal/models"
	"github.com/RomanenkoDR/Gofemart/iternal/services/auth"
	"github.com/RomanenkoDR/Gofemart/iternal/services/db"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func SetDatabase(database *gorm.DB) {
	db.Database = database
}

func Register(w http.ResponseWriter, r *http.Request) {

	err := json.NewDecoder(r.Body).Decode(&auth.User)
	if err != nil || auth.User.Login == "" || auth.User.Password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	// Проверка существования пользователя
	result := db.Database.Where("login = ?", auth.User.Login).First(&auth.ExistingUser)
	if result.RowsAffected > 0 {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Хэшируем пароль
	err = auth.User.HashPassword()
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Сохранение пользователя
	err = db.Database.Create(&auth.User).Error
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Создание связанной записи в таблице для нового юзера
	newBalance := models.Balance{
		UserID: auth.User.ID,
	}

	// Сохранение таблицы
	err = db.Database.Create(&newBalance).Error
	if err != nil {
		http.Error(w, "Failed to create balance record", http.StatusInternalServerError)
		return
	}
	// Генерация токена
	auth.JwtKey, err = config.GenerateJWT(auth.User.Login)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   auth.JwtKey,
		Expires: time.Now().Add(24 * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
}
