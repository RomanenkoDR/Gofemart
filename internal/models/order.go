package models

import "time"

type Order struct {
	ID          uint      `gorm:"primary_key"`
	UserID      uint      `gorm:"not null"`
	OrderNumber string    `gorm:"unique;not null"`
	Status      string    `gorm:"default:'NEW'"`
	Accrual     float64   `gorm:"default:null"`
	Sum         float64   `gorm:"default:null"`
	UploadedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type OrdersUserJSON struct {
	ID          uint      `gorm:"primary_key" json:"-"`
	UserID      uint      `gorm:"not null" json:"-"`
	OrderNumber string    `gorm:"unique;not null" json:"number"`
	Status      string    `gorm:"default:'NEW'" json:"status"`
	Accrual     float64   `gorm:"default:null" json:"accrual,omitempty"`
	UploadedAt  time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`
}

// AccrualInfo представляет информацию о расчёте начислений для заказа.
type AccrualInfo struct {
	OrderNumber string `json:"order"`
	Status      string
	Accrual     float64 `json:"accrual"`
}
