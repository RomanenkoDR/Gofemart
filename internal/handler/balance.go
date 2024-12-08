package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"net/http"
)

// Balance обрабатывает получение текущего баланса пользователя.
func (h *Handler) Balance(w http.ResponseWriter, r *http.Request) {

	// Получаем логин из запросов
	username := r.Header.Get("X-Username")

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
