package router

import (
	"github.com/RomanenkoDR/Gofemart/iternal/config"
	"github.com/RomanenkoDR/Gofemart/iternal/handler"
	"github.com/RomanenkoDR/Gofemart/iternal/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRouter(cfg config.Options, h handler.Handler) (chi.Router, error) {
	// Init rout for server
	router := chi.NewRouter()

	// Use router
	router.Use(middleware.LogHandler)

	router.Route("/api/users", func(r chi.Router) {

		r.Post("/register", handler.Register)
		r.Post("/logger", handler.Login)
		r.Post("/orders", handler.OrdersGet)
		r.Post("/balance/withdraw", handler.Withdraw)

		r.Get("/orders", handler.OrdersGet)
		r.Get("/balance", handler.Balance)
		r.Get("/withdrawals", handler.Withdrawals)
	})

	return router, nil
}
