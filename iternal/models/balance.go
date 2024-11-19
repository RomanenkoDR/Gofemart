package models

import "time"

type Balance struct {
	ID       uint64    `gorm:"primary_key"`
	UserID   uint64    `gorm:"not null"`
	User     User      `gorm:"foreigner:UserID"`
	Current  *float64  `gorm:"default:null"`
	Withdraw *float64  `gorm:"default:null"`
	UpdateAt time.Time `gorm:"autoCreateTime"`
}
