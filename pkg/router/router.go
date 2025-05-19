package router

import (
	"log/slog"

	"github.com/MTUCIBOY/MedodsTest/pkg/config"
	authtokens "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/authTokens"
	newuser "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/newUser"
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

	router.Post("/authTokens", authtokens.ATHandler(log, db))
	router.Post("/registrate", newuser.NUHandler(log, db))

	return router
}
