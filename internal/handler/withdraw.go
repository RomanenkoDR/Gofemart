package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"log"
	"net/http"
)

// Withdraw обрабатывает запрос на списание средств.
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var (
		user    models.User
		balance *models.Balance
	)

	// Структура для парсинга тела запроса
	var requestBody struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}

	// Получаем логин из запросов
	username := r.Header.Get("X-Username")

	// Получаем пользователя по логину
	if err := db.GetUserByLogin(h.DB, username, &user); err != nil {
		log.Printf("В Withdraw(POST) ошибка при получении id пользователя по логину: %s", err)
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	// Декодируем JSON
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("В Withdraw (POST) ошибка при парсинге json: %s", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверяем номер заказа алгоритмом Луна
	if !models.ValidLun(requestBody.Order) {
		log.Print("В Withdraw (POST) поступил невалидный номер заказа")
		http.Error(w, "Неверный номер заказа", http.StatusUnprocessableEntity)
		return
	}

	// Получаем баланс пользователя
	balance, err := db.GetUserBalance(h.DB, username)
	if err != nil {
		log.Print("В Withdraw (POST) ошибка при получении баланса пользователя.")
		http.Error(w, "Ошибка при получении баланса", http.StatusInternalServerError)
		return
	}

	// Проверяем, есть ли достаточно средств на счету
	if balance.Current <= requestBody.Sum {
		log.Print("В Withdraw (POST) ошибка при получении баланса пользователя. Недостаточно средств")
		http.Error(w, "Недостаточно средств на счету", http.StatusPaymentRequired)
		return
	}

	// Формируем модель заказа
	newOrder := &models.Order{
		OrderNumber: requestBody.Order,
		UserID:      user.ID,
		Sum:         requestBody.Sum,
	}

	// Создаем заказ в базе
	if err := db.CreateOrder(h.DB, newOrder); err != nil {
		log.Print("в Withdraw ошибка при создании заказ пользователя. Заказ уже есть в базе")
		http.Error(w, "Ошибка при создании заказа", http.StatusInternalServerError)
		return
	}

	// Обновляем current и withdranw
	balance.Current -= requestBody.Sum
	balance.Withdrawn += requestBody.Sum

	// Обновляем баланс пользователя в базе данных
	if err := db.UpdateUserBalance(h.DB, newOrder, balance.Current, balance.Withdrawn); err != nil {
		log.Printf("В Withdraw (POST) ошибка при обновлении баланса: %s", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	log.Print(w, "Средства успешно списаны с баланса")
}

// Withdrawals обрабатывает запрос на получение информации о выводах средств.
func (h *Handler) Withdrawals(w http.ResponseWriter, r *http.Request) {

	var user models.User
	// Получаем логин из запросов
	username := r.Header.Get("X-Username")

	// Получаем пользователя по имени из токена
	if err := db.GetUserByLogin(h.DB, username, &user); err != nil {
		log.Printf("В Withdraw (POST) ошибка при получении id пользователя по логину: %s", err)
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	// Извлекаем все выводы средств пользователя
	withdrawals, err := db.GetWithdrawalsByUserID(h.DB, user.ID)
	if err != nil {
		log.Printf("Ошибка при получении баланса в GetWithdrawalsByUserID: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	log.Printf("В ручке Withdrawals получили баланс: %v", withdrawals)

	// Если выводы отсутствуют, возвращаем статус 204
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Отправляем ответ с выводами
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(withdrawals)
	w.WriteHeader(http.StatusOK)
	log.Printf("Информация о средствах успешно отправлена пользователю. Ответ: %v", withdrawals)
}
