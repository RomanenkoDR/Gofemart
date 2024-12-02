package db

import (
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/Gofemart/internal/config"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
)

// CreateOrder создает новый заказ.
func CreateOrder(db *gorm.DB, order *models.Order) error {
	return db.Create(order).Error
}

// GetOrderByNumber получает заказ по номеру.
func GetOrderByNumber(db *gorm.DB, orderNumber string, order *models.Order) error {
	return db.Where("order_number = ?", orderNumber).First(order).Error
}

// GetOrdersByUserID получает заказы пользователя по ID.
func GetOrdersByUserID(db *gorm.DB, userID uint64, orders *[]models.Order) error {
	return db.Where("user_id = ?", userID).Order("uploaded_at DESC").Find(orders).Error
}

func UpdateOrderInfo(db *gorm.DB, numberOrder string) error {
	var (
		orderFromAccrualSystem models.AccrualInfo
	)

	urlAPI := fmt.Sprintf("%s/api/orders/%s", config.AccrualSystemURI, numberOrder)
	log.Printf("Обращение к API: %s", urlAPI)

	// Выполняем запрос
	resp, err := http.Get(urlAPI)
	if err != nil {
		return fmt.Errorf("ошибка при запросе к системе начислений: %v", err)
	}
	defer resp.Body.Close()

	// Обрабатываем статус ответа
	switch resp.StatusCode {
	case http.StatusOK:
		log.Printf("Получен ответ с кодом HTTP 200 OK")
	case http.StatusNoContent:
		log.Printf("Ответ с кодом HTTP 204 No Content")
		return nil
	case http.StatusInternalServerError:
		log.Printf("Ответ с кодом HTTP 500 Internal Server Error")
		return fmt.Errorf("ошибка сервера начислений")
	case http.StatusTooManyRequests:
		timeSlip := resp.Header.Get("Retry-After")
		intTimeSlip, _ := strconv.Atoi(timeSlip)
		log.Printf("Превышено количество запросов, подождите %d секунд", intTimeSlip)
		return fmt.Errorf("превышено количество запросов")
	default:
		return fmt.Errorf("неожиданный статус ответа: %d", resp.StatusCode)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Десериализуем JSON
	if err := json.Unmarshal(body, &orderFromAccrualSystem); err != nil {
		return fmt.Errorf("ошибка при разборе JSON: %v", err)
	}

	// Обновляем статус заказа
	if err := updateOrderStatus(db, orderFromAccrualSystem); err != nil {
		return fmt.Errorf("ошибка обновления статуса заказа: %v", err)
	}

	log.Printf("Информация о заказе успешно обновлена: %+v", orderFromAccrualSystem)
	return nil
}

// UpdateOrderStatus обновляет статус заказа в базе данных
func updateOrderStatus(db *gorm.DB, orderAccrual models.AccrualInfo) error {
	var order models.Order

	if err := db.Where("order_number = ?", orderAccrual.OrderNumber).
		Updates(map[string]interface{}{
			"status":  orderAccrual.Status,
			"accrual": orderAccrual.Accrual,
		}).First(&order).Error; err != nil {
		return err
	}

	if err := db.Where(models.Balance{}).
		Where("user_id = ?", order.UserID).
		Updates(map[string]interface{}{
			"current": order.Accrual,
		}).Error; err != nil {
		return err
	}

	log.Printf("Статус заказа %s успешно обновлен на %s", order.OrderNumber, order.Status)
	return nil
}
