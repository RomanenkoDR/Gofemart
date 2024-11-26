package db

import (
	"errors"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/gorm"
)

// GetUserBalance возвращает баланс пользователя по логину.
func GetUserBalance(db *gorm.DB, username string) (*models.Balance, error) {
	var user models.User
	var balance models.Balance

	// Ищем пользователя по логину
	if err := db.Where("login = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// Ищем баланс пользователя по user.ID
	if err := db.Where("user_id = ?", user.ID).First(&balance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &balance, nil
}

// UpdateUserBalance обновляет баланс пользователя в базе данных.
func UpdateUserBalance(db *gorm.DB, user *models.User) error {
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// GetWithdrawalsByUserID получает все выводы средств пользователя по его ID.
func GetWithdrawalsByUserID(db *gorm.DB, userID uint64, withdrawals *[]models.Withdrawal) error {
	// Предположим, что есть таблица для выводов с соответствующей моделью
	if err := db.Where("user_id = ?", userID).Order("processed_at DESC").Find(withdrawals).Error; err != nil {
		return err
	}
	return nil
}
