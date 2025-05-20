package checkuseragent

import (
	"context"
	"log/slog"
	"net/http"

	errorresponse "github.com/MTUCIBOY/MedodsTest/pkg/router/errorResponse"
	"github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/refresh"
	"github.com/go-chi/chi/v5/middleware"
)

type deauthUser interface {
	DeactivateRefreshToken(ctx context.Context, uuidToken string) error
}

func CUAMiddleware(log *slog.Logger, dau deauthUser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const fn = "middlewares.CheckUserAgent"
			log := log.With(
				slog.String("fn", fn),
				slog.String("requestID", middleware.GetReqID(r.Context())),
			)

			jwtUserAgent, ok := r.Context().Value(middlewares.UserAgentKey).(string)
			if !ok {
				log.Error("user agent not found")
				errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

				return
			}

			if jwtUserAgent != r.UserAgent() {
				log.Error(
					"user agents not match",
					slog.String("jwtUserAgent", jwtUserAgent),
					slog.String("userAgent", r.UserAgent()),
				)
				errorresponse.JSONResponde(w, http.StatusConflict, "UserAgents not match")

				refreshToken := r.Header.Get("Refresh-Token")

				refreshID, err := refresh.Check(refreshToken)
				if err != nil {
					log.Error("failed to check refresh token", slog.String("err", err.Error()))
					errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

					return
				}

				err = dau.DeactivateRefreshToken(r.Context(), refreshID)
				if err != nil {
					log.Error("failed to deactivate refresh token", slog.String("err", err.Error()))
					errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

					return
				}

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
