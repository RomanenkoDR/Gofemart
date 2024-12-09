package models

import "time"

// Withdrawal представляет информацию о выводах средств пользователя.
type WithdrawalsJSON struct {
	ID          uint64    `gorm:"primary_key" json:"-"`
	UserID      uint64    `gorm:"not null" json:"-"`
	OrderNumber string    `gorm:"unique;not null" json:"order"`
	Status      string    `gorm:"default:'NEW'" json:"-"`
	Accrual     float32   `gorm:"default:null" json:"-"`
	Sum         float32   `gorm:"default:null" json:"sum"`
	UploadedAt  time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"processed_at"`
}
