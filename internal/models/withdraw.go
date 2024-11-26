package models

import "time"

// Withdrawal представляет информацию о выводах средств пользователя.
type Withdrawal struct {
	Order          string    `json:"order"`
	Sum            int       `json:"sum"`
	ProcessedAt    time.Time `json:"processed_at"`
	ProcessedAtStr string    `json:"processed_at_str,omitempty"` // Добавлено для хранения строки
	UserID         uint64    `json:"user_id"`
}
