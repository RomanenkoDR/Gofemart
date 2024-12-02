package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"io"
	"log"
	"net/http"
	"strings"
)

//617449
//12345678903
//6011111111111117

// OrdersPost обрабатывает создание нового заказа.
func (h *Handler) OrdersPost(w http.ResponseWriter, r *http.Request) {
	var (
		existingOrder models.Order
		user          = &models.User{}
	)

	// Получаем логин из запросов
	username := r.Header.Get("X-Username")

	// Проверяем Content-Type
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Invalid Content-Type, expected text/plain", http.StatusBadRequest)
		return
	}

	// Читаем номер заказа
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	orderNumber := strings.TrimSpace(string(body))

	// Проверяем номер заказа
	if !models.ValidLun(orderNumber) {
		http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	// Получаем пользователя
	if err := db.GetUserByLogin(h.DB, username, user); err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Проверяем существование заказа
	if err := db.GetOrderByNumber(h.DB, orderNumber, &existingOrder); err == nil {
		if existingOrder.UserID == user.ID {
			http.Error(w, "Order already uploaded by this user", http.StatusOK)
		} else {
			http.Error(w, "Order already uploaded by another user", http.StatusConflict)
		}
		return
	}

	// Создаем новый заказ
	newOrder := &models.Order{
		OrderNumber: orderNumber,
		UserID:      user.ID,
	}
	if err := db.CreateOrder(h.DB, newOrder); err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Асинхронно обновляем заказ в системе начислений
	go func() {
		if err := db.UpdateOrderInfo(h.DB, orderNumber, h.AccrualSystemAddress); err != nil {
			log.Printf("Ошибка обновления информации о заказе %s: %v", orderNumber, err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}

// OrdersGet обрабатывает получение заказов пользователя.
func (h *Handler) OrdersGet(w http.ResponseWriter, r *http.Request) {
	var (
		user   = &models.User{}
		orders []models.Order
	)

	username := r.Header.Get("X-Username")

	// Получаем пользователя по логину
	if err := db.GetUserByLogin(h.DB, username, user); err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Получаем заказы пользователя
	if err := db.GetOrdersByUserID(h.DB, user.ID, &orders); err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	// Обновляем статусы заказов перед отправкой ответа
	for _, order := range orders {
		if err := db.UpdateOrderInfo(h.DB, order.OrderNumber, h.AccrualSystemAddress); err != nil {
			log.Printf("Ошибка обновления информации о заказе %s: %v", order.OrderNumber, err)
		}
	}

	// Если заказы отсутствуют
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Формируем ответ
	response := make([]map[string]interface{}, len(orders))
	for i, order := range orders {
		response[i] = map[string]interface{}{
			"order_number": order.OrderNumber,
			"status":       order.Status,
			"uploaded_at":  order.UploadedAt,
		}
		if &order.Accrual != nil {
			response[i]["accrual"] = order.Accrual
		}
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
