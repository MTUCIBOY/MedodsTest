package checkrefreshtoken

import (
	"context"
	"log/slog"
	"net/http"

	errorresponse "github.com/MTUCIBOY/MedodsTest/pkg/router/errorResponse"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/refresh"
	"github.com/go-chi/chi/v5/middleware"
)

type CheckRefreshToken interface {
	IsActiveRefresh(ctx context.Context, token, uuidToken string) (bool, error)
}

func CRTMiddleware(log *slog.Logger, crt CheckRefreshToken) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const fn = "middlewares.ExpireTokenValidator"
			log := log.With(
				slog.String("fn", fn),
				slog.String("requestID", middleware.GetReqID(r.Context())),
			)

			refreshToken := r.Header.Get("Refresh-Token")

			refreshID, err := refresh.Check(refreshToken)
			if err != nil {
				log.Error("failed to check refresh token", slog.String("err", err.Error()))
				errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

				return
			}

			isActive, err := crt.IsActiveRefresh(r.Context(), refreshToken, refreshID)
			if err != nil {
				log.Error("failed to check is_active refresh token", slog.String("err", err.Error()))
				errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

				return
			}

			if !isActive {
				log.Error("token not active")
				errorresponse.JSONResponde(w, http.StatusForbidden, "Refresh token is not active")

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
