package handler

import "gorm.io/gorm"

// Handler представляет структуру обработчиков с доступом к базе данных.
type Handler struct {
	DB                   *gorm.DB
	AccrualSystemAddress string
}

// NewHandler Create new handler and previous reports info from file it needed
func NewHandler(db *gorm.DB, accrualSystemAddress string) *Handler {
	return &Handler{
		DB:                   db,
		AccrualSystemAddress: accrualSystemAddress,
	}
}
