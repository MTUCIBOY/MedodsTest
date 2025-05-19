package router

import (
	"log/slog"

	"github.com/MTUCIBOY/MedodsTest/pkg/config"
	authtokens "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/authTokens"
	newuser "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/newUser"
	tokenvalidator "github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares/tokenValidator"
	"github.com/MTUCIBOY/MedodsTest/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(cfg *config.Config, log *slog.Logger, db storage.DB) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.RequestID,
	)

	router.Group(func(r chi.Router) {
		r.Use(
			tokenvalidator.TVMiddleware(log),
		)
	})

	router.Post("/authTokens", authtokens.ATHandler(log, cfg.TTLToken, db))
	router.Post("/registrate", newuser.NUHandler(log, db))

	return router
}
