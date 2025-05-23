package router

import (
	"log/slog"

	_ "github.com/MTUCIBOY/MedodsTest/docs"
	"github.com/MTUCIBOY/MedodsTest/pkg/config"
	authtokens "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/authTokens"
	deauthtokens "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/deauthTokens"
	getguid "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/getGUID"
	newuser "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/newUser"
	updatetokens "github.com/MTUCIBOY/MedodsTest/pkg/router/handlers/updateTokens"
	checkrefreshtoken "github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares/checkRefreshToken"
	checkuseragent "github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares/checkUserAgent"
	checkuserip "github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares/checkUserIP"
	expiretokenvalidator "github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares/expireTokenValidator"
	tokenvalidator "github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares/tokenValidator"
	"github.com/MTUCIBOY/MedodsTest/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8888/swagger/doc.json"),
	))

	router.Group(func(r chi.Router) {
		r.Use(checkrefreshtoken.CRTMiddleware(log, db))

		r.Group(func(r chi.Router) {
			r.Use(
				tokenvalidator.TVMiddleware(log),
			)

			r.Get("/guid", getguid.UUIDHadler(log, db))
			r.Post("/deauthTokens", deauthtokens.DATHandler(log, db))
		})

		r.Group(func(r chi.Router) {
			r.Use(
				expiretokenvalidator.ETVMiddleware(log),
				checkuseragent.CUAMiddleware(log, db),
				checkuserip.CUIPMiddleware(log, cfg.WebhookURL),
			)

			r.Post("/updateTokens", updatetokens.UTHandler(log, cfg.TTLToken))
		})
	})

	return router
}
