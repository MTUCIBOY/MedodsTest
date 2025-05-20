package tokenvalidator

import (
	"context"
	"log/slog"
	"net/http"

	errorresponse "github.com/MTUCIBOY/MedodsTest/pkg/router/errorResponse"
	"github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/access"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/refresh"
	"github.com/go-chi/chi/v5/middleware"
)

func TVMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const fn = "middlewares.TokenValidator"
			log := log.With(
				slog.String("fn", fn),
				slog.String("requestID", middleware.GetReqID(r.Context())),
			)

			accessToken := r.Header.Get("Access-Token")
			refreshToken := r.Header.Get("Refresh-Token")

			if accessToken == "" || refreshToken == "" {
				log.Error("missing tokens")
				errorresponse.JSONResponde(w, http.StatusUnauthorized, "Missing tokens")

				return
			}

			_, err := refresh.Check(refreshToken)
			if err != nil {
				log.Error("failed to check refresh token", slog.String("err", err.Error()))
				errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

				return
			}

			accessClaims, err := access.Check(accessToken, refreshToken)
			if err != nil {
				log.Error("failed to check access token", slog.String("err", err.Error()))
				errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

				return
			}

			ctx := context.WithValue(r.Context(), middlewares.UserEmailKey, accessClaims.Subject)
			ctx = context.WithValue(ctx, middlewares.UserAgentKey, accessClaims.UserAgent)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
