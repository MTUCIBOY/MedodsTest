package router

import (
	"log/slog"

	"github.com/MTUCIBOY/MedodsTest/pkg/config"
	authtokens "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/authTokens"
	deauthtokens "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/deauthTokens"
	getguid "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/getGUID"
	newuser "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/newUser"
	updatetokens "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/updateTokens"
	checkrefreshtoken "github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares/checkRefreshToken"
	expiretokenvalidator "github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares/expireTokenValidator"
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

	router.Post("/authTokens", authtokens.ATHandler(log, cfg.TTLToken, db))
	router.Post("/registrate", newuser.NUHandler(log, db))

	router.Group(func(r chi.Router) {
		r.Use(checkrefreshtoken.CRTMiddleware(log, db))

		r.Group(func(r chi.Router) {
			r.Use(tokenvalidator.TVMiddleware(log))

			r.Get("/guid", getguid.UUIDHadler(log, db))
			r.Post("/deauthTokens", deauthtokens.DATHandler(log, db))
		})

		r.Group(func(r chi.Router) {
			r.Use(expiretokenvalidator.ETVMiddleware(log))

			r.Post("/updateTokens", updatetokens.UTHandler(log, cfg.TTLToken))
		})
	})

	return router
}
