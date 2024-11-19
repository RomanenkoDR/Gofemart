package handler

import (
	"encoding/json"
	"errors"
	"github.com/RomanenkoDR/Gofemart/iternal/models"
	"github.com/RomanenkoDR/Gofemart/iternal/services/auth"
	"github.com/RomanenkoDR/Gofemart/iternal/services/db"
	"github.com/RomanenkoDR/Gofemart/iternal/services/orders"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strings"
)

//4000000000000002
//12345678903
//6011111111111117

func OrdersPost(w http.ResponseWriter, r *http.Request) {

	// Получаем результат проверки авторизации пользователя
	username, statusCode, err := auth.СheckAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Проверяем Content-Type
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Invalid Content-Type, expected text/plain", http.StatusBadRequest)
		return
	}

	// Считываем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	orderNumber := strings.TrimSpace(string(body))

	// Проверка номера заказа с помощью алгоритма Луна
	if !models.ValidLuhn(orderNumber) {
		http.Error(w, "Неверный формат номера заказа", http.StatusUnprocessableEntity)
		return
	}

	if err := db.Database.Where("login = ?", username).First(&auth.User).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Проверяем наличие заказа в базе данных

	result := db.Database.Where("order_number = ?", orderNumber).First(&orders.ExistingOrder)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Обработка ситуации, если заказ уже существует
	if result.RowsAffected > 0 {
		if orders.ExistingOrder.UserID == auth.User.ID {
			http.Error(w, "Номер заказа уже был загружен этим пользователем", http.StatusOK)
		} else {
			http.Error(w, "Номер заказа уже был загружен другим пользователем", http.StatusConflict)
		}
		return
	}

	// Создаём новый заказ
	newOrder := models.Order{
		OrderNumber: orderNumber,
		UserID:      auth.User.ID,
	}

	err = db.Database.Create(&newOrder).Error
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func OrdersGet(w http.ResponseWriter, r *http.Request) {

	username, statusCode, err := auth.СheckAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	// Проверяем Content-Type
	if r.Header.Get("Content-Length") != "0" {
		http.Error(w, "некорректный Content-Length, ожидается 0", http.StatusBadRequest)
		return
	}

	// Получаем пользователя из базы
	if err := db.Database.Where("login = ?", username).First(&auth.User).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = db.Database.Where("user_id = ?", auth.User.ID).
		Order("uploaded_at DESC").
		Find(&orders.Orders).Error
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Проверка на наличие заказов
	if len(orders.Orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Формируем список заказов в ответе
	response := make([]map[string]interface{}, 0, len(orders.Orders))
	for _, order := range orders.Orders {
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
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}
