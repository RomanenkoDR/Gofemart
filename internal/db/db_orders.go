package db

import (
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
func GetOrdersByUserID(db *gorm.DB, userID uint64, orders *[]models.Order) error {
	return db.Where("user_id = ?", userID).Order("uploaded_at DESC").Find(orders).Error
}

// GetAccrualInfoByOrderNumber получает информацию о расчёте начисления для заказа из таблицы orders.
func GetAccrualInfoByOrderNumber(db *gorm.DB, orderNumber string, accrualInfo *models.Order) error {

	// Запрос к базе данных для получения информации о начислении для заказа из таблицы orders
	if err := db.Where("order_number = ?", orderNumber).First(accrualInfo).Error; err != nil {
		log.Printf("Ошибка при получении информации о расчете для заказа %s: %v", orderNumber, err)
		return err
	}
	return nil
}

// UpdateOrderStatus обновляет статус заказа в базе данных
func UpdateOrderStatus(db *gorm.DB, order *models.Order) error {
	log.Printf("Ордер ID %v", order.ID)
	if err := db.Save(order).Error; err != nil {
		log.Printf("Ошибка при обновлении статуса заказа %s: %v", order.OrderNumber, err)
		return err
	}
	log.Printf("Статус заказа %s успешно обновлен на %s", order.OrderNumber, order.Status)
	return nil
}
