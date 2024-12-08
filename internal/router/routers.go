package router

import (
	"github.com/RomanenkoDR/Gofemart/internal/handler"
	"github.com/RomanenkoDR/Gofemart/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// SetupRouter создает и возвращает настроенный маршрутизатор.
func SetupRouter(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Основной middleware
	r.Group(func(r chi.Router) {
		r.Use(middleware.LogHandler)
		r.Use(middleware.GzipHandle)

		// Публичные маршруты
		r.Post("/api/user/register", h.Register) // Регистрация пользователя
		r.Post("/api/user/login", h.Login)       // Аутентификация пользователя

		// Приватные маршруты
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)

			r.Post("/api/user/orders", h.OrdersPost)         // Добавление номера заказа
			r.Post("/api/user/balance/withdraw", h.Withdraw) // Запрос на списание баллов

			r.Get("/api/user/orders", h.OrdersGet)        // Получение списка заказов пользователя
			r.Get("/api/user/balance", h.Balance)         // Получение текущего баланса
			r.Get("/api/user/withdrawals", h.Withdrawals) // Получение информации о выводах средств
		})
	})

	return r
}
