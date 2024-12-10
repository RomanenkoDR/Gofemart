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
		log.Printf("Некорректное тело запроса \n%v", err)
		http.Error(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	orderNumber := strings.TrimSpace(string(body))

	// Проверяем номер заказа
	if !models.ValidLun(orderNumber) {
		log.Print("В OrderPost поступил невалидный номер заказа")
		http.Error(w, "Неверный номер заказа", http.StatusUnprocessableEntity)
		return
	}

	// Получаем пользователя
	if err := db.GetUserByLogin(h.DB, username, user); err != nil {
		log.Printf("В OrderPost ошибка при получении id пользователя по логину: %s", err)
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	// Проверяем существование заказа
	if err := db.GetOrderByNumber(h.DB, orderNumber, &existingOrder); err == nil {
		if existingOrder.UserID == user.ID {
			log.Printf("Такой заказ уже загружен текущим пользователем: \n%v", err)
			http.Error(w, "Такой заказ уже загружен текущим пользователем", http.StatusOK)
		} else {
			log.Printf("Такой заказ уже загружен другим пользователем: \n%v", err)
			http.Error(w, "Такой заказ уже загружен другим пользователем", http.StatusConflict)
		}
		return
	}

	// Создаем новый заказ
	newOrder := &models.Order{
		OrderNumber: orderNumber,
		UserID:      user.ID,
	}

	if err := db.CreateOrder(h.DB, newOrder); err != nil {
		log.Printf("Ошибка при создании заказ: \n%v", err)
		http.Error(w, "ошибка при создании заказа", http.StatusInternalServerError)
		return
	}

	go db.UpdateOrderInfo(h.DB, orderNumber, models.Config.AccrualSystemAddress)

	w.WriteHeader(http.StatusAccepted)
}

// OrdersGet обрабатывает получение заказов пользователя.
func (h *Handler) OrdersGet(w http.ResponseWriter, r *http.Request) {
	var (
		user       = &models.User{}
		ordersJSON []models.OrdersUserJSON
	)

	username := r.Header.Get("X-Username")

	// Получаем пользователя по логину
	if err := db.GetUserByLogin(h.DB, username, user); err != nil {
		log.Printf("В OrderGet ошибка при получении id пользователя по логину: %s", err)
		http.Error(w, "пользователь не найден", http.StatusUnauthorized)
		return
	}

	// Получаем заказы пользователя
	log.Printf("В ручке OrdersGet отправляем запрос на получение заказов пользователя по id %v", user.ID)
	if err := db.GetOrdersByUserID(h.DB, user.ID, &ordersJSON); err != nil {
		log.Printf("Ошибка при поиске заказов: %s", err)
		http.Error(w, "ошибка при поиске заказов", http.StatusInternalServerError)
		return
	}

	// Если заказы отсутствуют
	if len(ordersJSON) == 0 {
		log.Print("В OrdersGet заказы отсутствуют")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ordersJSON); err != nil {
		log.Printf("В ручке OrdersGet ошибка при Кодировании json: %s", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
