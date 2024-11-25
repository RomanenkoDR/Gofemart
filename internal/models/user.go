package models

import (
	"github.com/golang-jwt/jwt/v5"
)

var (
	User         userStr
	Claims       jwt.MapClaims
	JwtKey       string
	ExistingUser userStr
)

// User Структура таблицы пользователей
type userStr struct {
	ID       uint64 `gorm:"primary_key"`
	Login    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}
