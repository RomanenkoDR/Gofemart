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
	log.Print("В ручке Balance Отправили запрос на получение баланса")
	balance, err := db.GetUserBalance(h.DB, username)
	if err != nil {
		log.Printf("В ручке Balance ошибка при получение баланса пользователя %s", err)
		http.Error(w, "Failed to retrieve user balance", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	log.Print("В ручке Balance формируем ответ")
	response := map[string]interface{}{
		"current":  balance.Current,
		"withdraw": balance.Withdraw,
	}
	log.Print("В ручке Balance устанавливаем ответ в body")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
