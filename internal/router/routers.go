package router

import (
	"github.com/RomanenkoDR/Gofemart/internal/config"
	"github.com/RomanenkoDR/Gofemart/internal/handler"
	"github.com/RomanenkoDR/Gofemart/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRouter(cfg config.Options, h handler.Handler) (chi.Router, error) {
	// Init rout for server
	router := chi.NewRouter()

	// Use router
	router.Use(middleware.LogHandler)

	// Group api/user
	router.Route("/api/user", func(r chi.Router) {

		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
		r.Post("/orders", handler.OrdersPost)
		r.Post("/balance/withdraw", handler.Withdraw)

		r.Get("/orders", handler.OrdersGet)
		r.Get("/balance", handler.Balance)
		r.Get("/withdrawals", handler.Withdrawals)
	})

	return router, nil
}
