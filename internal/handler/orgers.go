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

	//Проверяем Content-Type
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

	log.Printf("В OrderPost Отправляем запрос на создание заказа newOrder")
	if err := db.CreateOrder(h.DB, newOrder); err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	log.Printf("В OrderPost отправляем запрос в систему лояльности")
	go db.UpdateOrderInfo(h.DB, orderNumber, models.Config.AccrualSystemAddress)

	w.WriteHeader(http.StatusAccepted)
}

// OrdersGet обрабатывает получение заказов пользователя.
func (h *Handler) OrdersGet(w http.ResponseWriter, r *http.Request) {
	var (
		user = &models.User{}
		//orders     = []models.Order{}
		ordersJSON = []models.OrdersUserJSON{}
	)

	username := r.Header.Get("X-Username")

	// Получаем пользователя по логину
	if err := db.GetUserByLogin(h.DB, username, user); err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Получаем заказы пользователя
	if err := db.GetOrdersByUserID(h.DB, user.ID, &ordersJSON); err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	// Если заказы отсутствуют
	if len(ordersJSON) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ordersJSON); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
