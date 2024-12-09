package models

import "time"

type Order struct {
	ID          uint64    `gorm:"primary_key"`
	UserID      uint64    `gorm:"not null"`
	OrderNumber string    `gorm:"unique;not null"`
	Status      string    `gorm:"default:'NEW'"`
	Accrual     float32   `gorm:"default:null"`
	Sum         float32   `gorm:"default:null"`
	UploadedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type OrdersUserJSON struct {
	ID          uint64    `gorm:"primary_key" json:"-"`
	UserID      uint64    `gorm:"not null" json:"-"`
	OrderNumber string    `gorm:"unique;not null" json:"number"`
	Status      string    `gorm:"default:'NEW'" json:"status"`
	Accrual     float32   `gorm:"default:null" json:"accrual,omitempty"`
	UploadedAt  time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`
}

// AccrualInfo представляет информацию о расчёте начислений для заказа.
type AccrualInfo struct {
	OrderNumber string `json:"order"`
	Status      string
	Accrual     float32 `json:"accrual"`
}
