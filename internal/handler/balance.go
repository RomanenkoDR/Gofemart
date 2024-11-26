package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/services"
	"net/http"
)

// Balance обрабатывает получение текущего баланса пользователя.
func (h *Handler) Balance(w http.ResponseWriter, r *http.Request) {
	// Проверяем авторизацию
	username, statusCode, err := services.СheckAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Получаем баланс пользователя
	balance, err := db.GetUserBalance(h.DB, username)
	if err != nil {
		http.Error(w, "Failed to retrieve user balance", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	response := map[string]interface{}{
		"current":   balance.Current,
		"withdrawn": balance.Withdraw,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
