package models

import "time"

// Order Структура для заказов
type Order struct {
	ID          uint64    `gorm:"primary_key"`
	OrderNumber string    `gorm:"unique;not null"`
	Status      string    `gorm:"default:'NEW'"`
	UserID      uint64    `gorm:"not null"`
	User        User      `gorm:"foreigner:UserID"`
	UploadedAt  time.Time `gorm:"autoCreateTime"`
	Accrual     *float64  `gorm:"default:null"`
}

type StatusOrder struct {
	NEW        string
	PROCESSING string
	INVALID    string
	PROCESSED  string
}
