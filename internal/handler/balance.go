package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"log"
	"net/http"
)

// Balance обрабатывает получение текущего баланса пользователя.
func (h *Handler) Balance(w http.ResponseWriter, r *http.Request) {

	// Получаем логин из запросов
	username := r.Header.Get("X-Username")

	// Получаем баланс пользователя
	balance, err := db.GetUserBalance(h.DB, username)
	if err != nil {
		log.Printf("В Balance/GetUserBalance ошибка при получении баланса пользователя %v", err)
		http.Error(w, "ошибка при получении баланса пользователя", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	response := map[string]interface{}{
		"current":   balance.Current,
		"withdrawn": balance.Withdrawn,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	w.WriteHeader(http.StatusOK)
}
