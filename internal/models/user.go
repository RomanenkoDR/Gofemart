package models

// User Структура таблицы пользователей
type User struct {
	ID       uint   `gorm:"primary_key"`
	Login    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}
