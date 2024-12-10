package db

import (
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
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
func GetOrdersByUserID(db *gorm.DB, userID uint, orders *[]models.OrdersUserJSON) error {
	return db.Model(models.Order{}).Where("user_id = ?", userID).Order("uploaded_at DESC").Find(orders).Error
}

func UpdateOrderInfo(db *gorm.DB, numberOrder string, accrualSystemAddress string) error {
	var orderFromAccrualSystem models.AccrualInfo

	// Формируем URL для запроса
	urlAPI := fmt.Sprintf("%s/api/orders/%s", accrualSystemAddress, numberOrder)
	log.Printf("Обращение к API: %s", urlAPI)

	// Выполняем запрос в систему лояльности
	resp, err := http.Get(urlAPI)
	if err != nil {
		return fmt.Errorf("ошибка при запросе к системе лояльности: %w", err)
	}
	defer resp.Body.Close()

	//Обрабатываем статус ответа
	switch resp.StatusCode {
	case http.StatusOK:
		log.Printf("В системе лояльности получен код: %d", http.StatusOK)
	case http.StatusNoContent:
		log.Printf("В системе лояльности получен код: %d", http.StatusNoContent)
		return nil
	case http.StatusInternalServerError:
		log.Printf("В системе лояльности получен код: %d", http.StatusInternalServerError)
		return fmt.Errorf("ошибка сервера начислений")
	case http.StatusTooManyRequests:
		log.Printf("В системе лояльности получен код: %d", http.StatusTooManyRequests)
		return fmt.Errorf("превышено количество запросов")
	default:
		log.Printf("В системе лояльности получен код: %d", resp.StatusCode)
		return fmt.Errorf("неожиданный статус ответа: %d", resp.StatusCode)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении тела ответа из системы лояльности. \nERR: %v", err)
		return fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Десериализуем JSON
	if err := json.Unmarshal(body, &orderFromAccrualSystem); err != nil {
		log.Printf("Ошибка при Unmarshal JSON. \nERR: %v", err)
		return fmt.Errorf("ошибка при разборе JSON: %v", err)
	}

	// Обновляем статус заказа и баланс пользователя
	if err := updateOrderStatus(db, orderFromAccrualSystem); err != nil {
		log.Printf("Ошибка при обновлении статуса заказа. \nERR: %v", err)
		return fmt.Errorf("ошибка обновления статуса заказа: %w", err)
	}

	log.Printf("Информация о заказе успешно обновлена: \n%v", orderFromAccrualSystem)
	return nil
}

// UpdateOrderStatus обновляет статус заказа в базе данных
func updateOrderStatus(db *gorm.DB, orderAccrual models.AccrualInfo) error {

	var order models.Order

	// Обновляем заказ в таблице
	if err := db.Model(&models.Order{}).
		Where("order_number = ?", orderAccrual.OrderNumber).
		Updates(map[string]interface{}{
			"status":  orderAccrual.Status,
			"accrual": orderAccrual.Accrual,
		}).First(&order).Error; err != nil {
		return fmt.Errorf("ошибка при обновлении заказа: \n%w", err)
	}

	// Обновляем баланс пользователя
	if err := db.Model(&models.Balance{}).
		Where("user_id = ?", order.UserID).
		Update("current", gorm.Expr("current + ?", orderAccrual.Accrual)).Error; err != nil {
		return fmt.Errorf("ошибка при обновлении баланса пользователя: %w", err)
	}

	log.Printf("Баланс пользователя %d успешно обновлён на: %f", order.UserID, orderAccrual.Accrual)
	return nil
}
