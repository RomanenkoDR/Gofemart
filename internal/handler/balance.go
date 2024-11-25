package handler

import (
	"encoding/json"
	"errors"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func Balance(w http.ResponseWriter, r *http.Request) {

	// Получаем результат проверки авторизации пользователя
	username, statusCode, err := models.СheckAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Проверяем Content-Type
	if r.Header.Get("Content-Length") != "0" {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Получаем пользователя из базы
	if err := models.Database.Where("login = ?", username).First(&models.User).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = models.Database.Where("user_id = ?", models.User.ID).Find(&models.Balance).Error
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	balanceResponse := map[string]interface{}{
		"current":   models.Balance.Current,
		"withdrawn": models.Balance.Withdraw,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(balanceResponse); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}

}
