package db

import (
	"errors"
	"github.com/RomanenkoDR/Gofemart/internal/models"
	"gorm.io/gorm"
)

// CreateUserWithBalance CreateUser создает пользователя и связанную запись в таблице баланса.
func CreateUserWithBalance(db *gorm.DB, user *models.User) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Создаем пользователя
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// Создаем баланс пользователя
		balance := models.Balance{UserID: user.ID}
		if err := tx.Create(&balance).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetUserByLogin получает пользователя по логину.
func GetUserByLogin(db *gorm.DB, login string, user *models.User) error {
	return db.Where("login = ?", login).First(user).Error
}

// CheckUserExists проверяет существование пользователя.
func CheckUserExists(db *gorm.DB, login string) (bool, error) {
	var user models.User
	err := db.Where("login = ?", login).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
