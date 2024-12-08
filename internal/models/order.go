package models

import "time"

type Order struct {
	ID          uint64    `gorm:"primary_key"`
	UserID      uint64    `gorm:"not null"`
	OrderNumber string    `gorm:"unique;not null"`
	Status      string    `gorm:"default:'NEW'"`
	Accrual     float64   `gorm:"default:null"`
	Sum         float64   `gorm:"default:null"`
	UploadedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// AccrualInfo представляет информацию о расчёте начислений для заказа.
type AccrualInfo struct {
	OrderNumber string `json:"order"`
	Status      string `json:"status"`
	Accrual     int    `json:"accrual""`
}
