package models

import "time"

// Withdrawal представляет информацию о выводах средств пользователя.
type WithdrawalsJSON struct {
	ID          uint      `gorm:"primary_key" json:"-"`
	UserID      uint      `gorm:"not null" json:"-"`
	OrderNumber string    `gorm:"unique;not null" json:"order"`
	Status      string    `gorm:"default:'NEW'" json:"-"`
	Accrual     float64   `gorm:"default:null" json:"-"`
	Sum         float64   `gorm:"default:null" json:"sum"`
	UploadedAt  time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"processed_at"`
}
