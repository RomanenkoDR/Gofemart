package models

var (
	Balance balanceStr
)

type balanceStr struct {
	ID       uint64   `gorm:"primary_key"`
	UserID   uint64   `gorm:"not null"`
	User     userStr  `gorm:"foreigner:UserID"`
	Current  *float64 `gorm:"default:null"`
	Withdraw *float64 `gorm:"default:null"`
}
