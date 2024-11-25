package db

import (
	"errors"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/gorm"
)

// GetOrderByNumber ищет заказ по номеру.
func GetOrderByNumber(orderNumber string) error {
	result := models.Database.Where("order_number = ?", orderNumber).First(&models.ExistingOrder)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	return nil
}

// CreateOrder создает новый заказ.
func CreateOrder() error {
	return models.Database.Create(&models.Order).Error
}

// GetOrdersByUserID получает заказы пользователя по ID.
func GetOrdersByUserID(userID uint64) error {
	return models.Database.Where("user_id = ?", userID).Order("uploaded_at DESC").Find(&models.Orders).Error
}
