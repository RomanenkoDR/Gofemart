package handler

import "gorm.io/gorm"

// Handler представляет структуру обработчиков с доступом к базе данных.
type Handler struct {
	DB *gorm.DB
}

// NewHandler Create new handler and previous reports info from file it needed
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB: db,
	}
}
