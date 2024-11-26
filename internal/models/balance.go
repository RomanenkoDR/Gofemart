package models

type Balance struct {
	ID       uint64  `gorm:"primary_key"`
	UserID   uint64  `gorm:"not null;unique"`
	Current  float64 `gorm:"default:0"`
	Withdraw float64 `gorm:"default:0"`
}
