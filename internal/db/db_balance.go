package db

import (
	"errors"
	"fmt"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/gorm"
	"log"
)

// GetUserBalance возвращает баланс пользователя по логину.
func GetUserBalance(db *gorm.DB, username string) (*models.Balance, error) {
	var user models.User
	var balance models.Balance

	// Ищем пользователя по логину
	log.Printf("В ручке Balance в функции GetUserBalance ищем пользователя по логину %s", username)
	if err := db.
		Where("login = ?", username).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("В ручке Balance в функции GetUserBalance ошибка при поиске пользователя по логину: %s", err)
			return nil, gorm.ErrRecordNotFound
		}
		log.Printf("В ручке Balance в функции GetUserBalance ошибка при поиске пользователя по логину: %s", err)
		return nil, err
	}

	// Ищем баланс пользователя по user.ID
	log.Printf("В ручке Balance в функции GetUserBalance ищем баланс по id %v", user.ID)
	if err := db.Where("user_id = ?", user.ID).First(&balance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("В ручке Balance в функции GetUserBalance ошибка при поиске баланса по id: %s", err)
			return nil, gorm.ErrRecordNotFound
		}
		log.Printf("В ручке Balance в функции GetUserBalance ошибка при поиске баланса по id: %s", err)
		return nil, err
	}
	log.Printf("В ручке Balance в функции GetUserBalance нашли баланс пользователя по логину: %v", user.ID)
	return &balance, nil
}

// UpdateUserBalance обновляет баланс пользователя в базе данных.
func UpdateUserBalance(db *gorm.DB, orderAccrual *models.Order) error {
	if err := db.
		Model(&models.Balance{}).
		Where("user_id = ?", orderAccrual.UserID).
		Update("current", gorm.Expr("current - ?", orderAccrual.Sum)).Error; err != nil {
		return fmt.Errorf("ошибка при обновлении баланса пользователя: %w", err)
	}
	return nil
}

// GetWithdrawalsByUserID получает все выводы средств пользователя по его ID.
func GetWithdrawalsByUserID(db *gorm.DB, userID uint64) (withdrawals []models.WithdrawalsJSON, err error) {
	if err = db.
		Model(&models.Order{}).
		Where("user_id = ? AND sum IS NOT NULL AND sum != 0", userID).
		Order("updated_at DESC").
		Find(&withdrawals).Error; err != nil {
		return nil, err
	}
	return withdrawals, nil
}
