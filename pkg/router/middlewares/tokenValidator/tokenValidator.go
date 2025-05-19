package tokenvalidator

import (
	"errors"
	"log/slog"
	"net/http"

	errorresponse "github.com/MTUCIBOY/MedodsTest/pkg/router/errorResponse"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/access"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/refresh"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
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

			_, err := access.CheckToken(accessToken)
			if err != nil {
				log.Error("failed to check access token", slog.String("err", err.Error()))

				if errors.Is(err, tokens.ErrInvalidToken) || errors.Is(err, jwt.ErrTokenExpired) {
					errorresponse.JSONResponde(w, http.StatusUnauthorized, "Invalid tokens")

					return
				}

				errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

				return
			}

			err = refresh.CheckToken(refreshToken, accessToken)
			if err != nil {
				log.Error("failed to check refresh token", slog.String("err", err.Error()))

				if errors.Is(err, tokens.ErrInvalidToken) {
					errorresponse.JSONResponde(w, http.StatusUnauthorized, "Invalid tokens")
				}

				errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
