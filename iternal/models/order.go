package models

type Order struct {
	ID          uint   `gorm:"primaryKey"`
	OrderNumber string `gorm:"uniqueIndex;not null"`
	UserID      uint   `gorm:"not null"`
	Status      string `gorm:"default:processing;not null"`
}
