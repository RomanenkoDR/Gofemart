package models

import "time"

var (
	Order         orderStr
	Orders        []orderStr
	ExistingOrder orderStr
)

// OrderStr Order Структура для заказов
type orderStr struct {
	ID          uint64    `gorm:"primary_key"`
	User        userStr   `gorm:"foreigner:UserID"`
	UserID      uint64    `gorm:"not null"`
	OrderNumber string    `gorm:"unique;not null"`
	Status      string    `gorm:"default:'NEW'"`
	Accrual     *float64  `gorm:"default:null"`
	Sum         *float64  `gorm:"default:null"`
	UploadedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type StatusOrder struct {
	NEW        string
	PROCESSING string
	INVALID    string
	PROCESSED  string
}
