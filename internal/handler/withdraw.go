package handler

import (
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"github.com/RomanenkoDR/Gofemart/internal/services"
	"net/http"
	"time"
)

// Withdraw обрабатывает запрос на списание средств.
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {

	// Получаем логин из запросов
	username := r.Header.Get("X-Username")

	// Парсим тело запроса
	var requestBody struct {
		Order string `json:"order"`
		Sum   int    `json:"sum"`
	}

	// Декодируем JSON
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверяем номер заказа (например, с помощью алгоритма Луна)
	if !models.ValidLun(requestBody.Order) {
		http.Error(w, "Неверный номер заказа", http.StatusUnprocessableEntity)
		return
	}

	// Получаем пользователя по токену
	var (
		user    models.User
		balance models.Balance
	)
	if err := db.GetUserByLogin(h.DB, username, &user); err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	// Проверяем, есть ли достаточно средств на счету
	if balance.Current < float64(requestBody.Sum) {
		http.Error(w, "Недостаточно средств на счету", http.StatusPaymentRequired)
		return
	}

	// Списываем средства
	balance.Current -= float64(requestBody.Sum)

	// Обновляем баланс пользователя в базе данных
	if err := db.UpdateUserBalance(h.DB, &user); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	fmt.Println(w, "Средства успешно списаны с баланса")
}

// Withdrawals обрабатывает запрос на получение информации о выводах средств.
func (h *Handler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	// Проверяем авторизацию
	username, statusCode, err := services.СheckAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Получаем пользователя по имени из токена
	var user models.User
	if err := db.GetUserByLogin(h.DB, username, &user); err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	// Извлекаем все выводы средств пользователя
	var withdrawals []models.Withdrawal
	if err := db.GetWithdrawalsByUserID(h.DB, user.ID, &withdrawals); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Если выводы отсутствуют, возвращаем статус 204
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Сортируем выводы по времени от новых к старым
	// (Это будет уже сделано в SQL запросе через Order)
	// Можно добавлять сортировку вручную, если нужно

	// Преобразуем поле времени в строку формата RFC3339
	for i := range withdrawals {
		withdrawals[i].ProcessedAtStr = withdrawals[i].ProcessedAt.Format(time.RFC3339)
	}

	// Отправляем ответ с выводами
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(withdrawals); err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}
