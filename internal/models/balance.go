package models

import "time"

type Balance struct {
	ID        uint      `gorm:"primary_key"`
	UserID    uint      `gorm:"not null;unique"`
	Current   float64   `gorm:"default:0"`
	Withdraw  float64   `gorm:"default:0"`
	ProcessAt time.Time `gorm:"autoCreateTime"`
}

type BalanceJSON struct {
	ID        uint      `gorm:"primary_key" json:"-"`
	UserID    uint      `gorm:"not null" json:"-"`
	Current   float64   `gorm:"default:null" json:"current"`
	Withdraw  float64   `gorm:"default:null" json:"withdraw"`
	ProcessAt time.Time `gorm:"autoCreateTime" json:"-"`
}
