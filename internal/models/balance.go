package models

import "time"

type Balance struct {
	ID        uint64    `gorm:"primary_key"`
	UserID    uint64    `gorm:"not null;unique"`
	Current   float32   `gorm:"default:0"`
	Withdraw  float32   `gorm:"default:0"`
	ProcessAt time.Time `gorm:"autoCreateTime"`
}

type BalanceJSON struct {
	ID        uint64    `gorm:"primary_key" json:"-"`
	UserID    uint64    `gorm:"not null" json:"-"`
	Current   float32   `gorm:"default:null" json:"current"`
	Withdraw  float32   `gorm:"default:null" json:"withdraw"`
	ProcessAt time.Time `gorm:"autoCreateTime" json:"-"`
}
