package handler

import (
	"encoding/json"
	"errors"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strings"
)

// OrdersPost обрабатывает создание нового заказа.
func OrdersPost(w http.ResponseWriter, r *http.Request) {

	// Проверяем авторизацию
	username, statusCode, err := models.СheckAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Проверяем Content-Type
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Invalid Content-Type, expected text/plain", http.StatusBadRequest)
		return
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	orderNumber := strings.TrimSpace(string(body))

	// Проверяем номер заказа (алгоритм Луна)
	if !models.ValidLuhn(orderNumber) {
		http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	// Получаем пользователя
	if err := db.GetUserByLogin(username); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Проверяем существование заказа
	if err := db.GetOrderByNumber(orderNumber); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Если заказ уже существует
	if models.ExistingOrder.ID != 0 {
		if models.ExistingOrder.UserID == models.User.ID {
			http.Error(w, "Order number already uploaded by this user", http.StatusOK)
		} else {
			http.Error(w, "Order number already uploaded by another user", http.StatusConflict)
		}
		return
	}

	// Создаем новый заказ
	models.Order.OrderNumber = orderNumber
	models.Order.UserID = models.User.ID

	if err := db.CreateOrder(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// OrdersGet обрабатывает получение заказов пользователя.
func OrdersGet(w http.ResponseWriter, r *http.Request) {

	// Проверяем авторизацию
	username, statusCode, err := models.СheckAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Проверяем Content-Length
	if r.Header.Get("Content-Length") != "0" {
		http.Error(w, "Invalid Content-Length, expected 0", http.StatusBadRequest)
		return
	}

	// Получаем пользователя
	if err := db.GetUserByLogin(username); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Получаем заказы пользователя
	if err := db.GetOrdersByUserID; err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Если заказы отсутствуют
	if len(models.Orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Формируем ответ
	response := make([]map[string]interface{}, 0, len(models.Orders))
	for _, order := range models.Orders {
		orderResponse := map[string]interface{}{
			"order_number": order.OrderNumber,
			"status":       order.Status,
			"uploaded_at":  order.UploadedAt,
		}
		if order.Accrual != nil {
			orderResponse["accrual"] = *order.Accrual
		}
		response = append(response, orderResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
	}
}
