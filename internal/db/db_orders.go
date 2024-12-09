package db

import (
	"fmt"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/gorm"
	"log"
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
func GetOrdersByUserID(db *gorm.DB, userID uint64, orders *[]models.OrdersUserJSON) error {
	return db.Model(models.Order{}).Where("user_id = ?", userID).Order("uploaded_at DESC").Find(orders).Error
}

func UpdateOrderInfo(db *gorm.DB, numberOrder string, accrualSystemAddress string) error {
	var (
		orderFromAccrualSystem models.AccrualInfo
	)
	//
	//// Формируем URL для запроса
	//urlAPI := fmt.Sprintf("%s/api/orders/%s", accrualSystemAddress, numberOrder)
	//log.Printf("Обращение к API: %s", urlAPI)
	//
	//// Выполняем запрос
	//log.Print("В UpdateOrderInfo отправляем запрос в систему лояльности")
	//
	//resp, err := http.Get(urlAPI)
	//if err != nil {
	//	log.Print("В UpdateOrderInfo получили ошибку при запросе в систему лояльности")
	//	return fmt.Errorf("ошибка при запросе к системе начислений: %w", err)
	//}
	//defer resp.Body.Close()
	//
	////Обрабатываем статус ответа
	//switch resp.StatusCode {
	//case http.StatusOK:
	//	log.Print("В UpdateOrderInfo Получен ответ с кодом HTTP 200 OK")
	//case http.StatusNoContent:
	//	log.Print("Ответ с кодом HTTP 204 No Content")
	//	return nil
	//case http.StatusInternalServerError:
	//	log.Print("Ответ с кодом HTTP 500 Internal Server Error")
	//	return fmt.Errorf("ошибка сервера начислений")
	//case http.StatusTooManyRequests:
	//	return fmt.Errorf("превышено количество запросов")
	//default:
	//	//return fmt.Errorf("неожиданный статус ответа: %d", resp.StatusCode)
	//}
	//
	//log.Print("В UpdateOrderInfo читаем тело ответа")
	//// Читаем тело ответа
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return fmt.Errorf("ошибка при чтении ответа: %w", err)
	//}
	//
	//log.Printf("В UpdateOrderInfo Десериализуем JSON %s", body)
	//
	//// Десериализуем JSON
	//if err := json.Unmarshal(body, &orderFromAccrualSystem); err != nil {
	//	log.Printf("В UpdateOrderInfo внутри json.Unmarshal. %s", err)
	//
	//	return fmt.Errorf("ошибка при разборе JSON: %w", err)
	//}

	orderFromAccrualSystem = models.AccrualInfo{
		OrderNumber: numberOrder,
		Status:      "PROCESSED",
		Accrual:     735.123,
	}

	log.Print("В UpdateOrderInfo Обновляем статус заказа и баланс пользователя")

	// Обновляем статус заказа и баланс пользователя
	if err := updateOrderStatus(db, orderFromAccrualSystem); err != nil {
		return fmt.Errorf("ошибка обновления статуса заказа: %w", err)
	}

	log.Printf("Информация о заказе успешно обновлена: %+v", orderFromAccrualSystem)
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
		return fmt.Errorf("ошибка при обновлении заказа: %w", err)
	}

	log.Printf("Заказ %s успешно обновлён со статусом %s и начислением %.2f",
		orderAccrual.OrderNumber, orderAccrual.Status, orderAccrual.Accrual)

	// Обновляем баланс пользователя
	if err := db.Model(&models.Balance{}).
		Where("user_id = ?", order.UserID).
		Update("current", gorm.Expr("current + ?", orderAccrual.Accrual)).Error; err != nil {
		return fmt.Errorf("ошибка при обновлении баланса пользователя: %w", err)
	}

	log.Printf("Баланс пользователя %d успешно обновлён на %.2f", order.UserID, orderAccrual.Accrual)
	return nil
}
