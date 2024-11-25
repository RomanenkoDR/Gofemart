package db

import (
	"errors"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/gorm"
)

// GetUserByLogin Получаем id пользователя по его логину
func GetUserByLogin(login string) error {
	result := models.Database.Where("login = ?", login).First(&models.User)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil // Пользователь не найден
	}
	return result.Error
}

// CheckUserExists проверяет существование пользователя.
func CheckUserExists(login string) (bool, error) {
	result := models.Database.Where("login = ?", login).First(&models.ExistingUser)
	if result.RowsAffected > 0 {
		return true, nil
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, result.Error
}

// CreateUser создает пользователя.
func CreateUser() error {
	return models.Database.Create(&models.User).Error
}

// CreateBalance создает баланс пользователя.
func CreateBalance() error {
	return models.Database.Create(&models.Balance).Error
}
