package handler

import (
	"encoding/json"
	"github.com/RomanenkoDR/Gofemart/internal/db"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"github.com/RomanenkoDR/Gofemart/internal/services"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

// OrdersPost обрабатывает создание нового заказа.
func (h *Handler) OrdersPost(w http.ResponseWriter, r *http.Request) {
	var (
		existingOrder *models.Order
		user          = &models.User{}
	)

	// Проверяем авторизацию
	username, statusCode, err := services.СheckAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

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
	if err := db.GetOrderByNumber(h.DB, orderNumber, existingOrder); err == nil {
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

	w.WriteHeader(http.StatusAccepted)
}

// OrdersGet обрабатывает получение заказов пользователя.
func (h *Handler) OrdersGet(w http.ResponseWriter, r *http.Request) {
	var (
		user   = &models.User{}
		orders []models.Order
	)

	// Проверяем авторизацию
	username, statusCode, err := services.СheckAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

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
		if order.Accrual != nil {
			response[i]["accrual"] = *order.Accrual
		}
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetOrderAccrual обрабатывает запрос на получение информации о расчёте начисления баллов для заказа.
func (h *Handler) GetOrderAccrual(w http.ResponseWriter, r *http.Request) {
	// Извлекаем номер заказа из URL
	orderNumber := chi.URLParam(r, "number")

	// Получаем информацию о заказе из базы данных
	var order models.Order
	if err := db.GetAccrualInfoByOrderNumber(h.DB, orderNumber, &order); err != nil {
		// Возвращаем ошибку, если не удалось найти заказ
		http.Error(w, "Заказ не зарегистрирован в системе расчёта", http.StatusNoContent)
		return
	}

	if err := db.UpdateOrderStatus(h.DB, &order); err != nil {
		http.Error(w, "Не удалось обновить статус заказа", http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	response := struct {
		Order   string   `json:"order"`
		Status  string   `json:"status"`
		Accrual *float64 `json:"accrual,omitempty"`
	}{
		Order:  orderNumber,
		Status: order.Status,
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}
