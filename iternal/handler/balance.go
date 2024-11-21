package handler

import (
	"encoding/json"
	"errors"
	"github.com/RomanenkoDR/Gofemart/iternal/services/auth"
	"github.com/RomanenkoDR/Gofemart/iternal/services/balance"
	"github.com/RomanenkoDR/Gofemart/iternal/services/db"
	"gorm.io/gorm"
	"net/http"
)

func Balance(w http.ResponseWriter, r *http.Request) {

	// Получаем результат проверки авторизации пользователя
	username, statusCode, err := auth.СheckAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Проверяем Content-Type
	if r.Header.Get("Content-Length") != "0" {
		// TODO: хз почему надо только внутренние ошибки сервера и статус не авторизован
		// TODO: пока поменял все выводы на внутренние ошибки сервера
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Получаем пользователя из базы
	if err := db.Database.Where("login = ?", username).First(&auth.User).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = db.Database.Where("user_id = ?", auth.User.ID).Find(&balance.Balance).Error
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	balanceResponce := map[string]interface{}{
		"current":   balance.Balance.Current,
		"withdrawn": balance.Balance.Withdraw,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(balanceResponce); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}

}
