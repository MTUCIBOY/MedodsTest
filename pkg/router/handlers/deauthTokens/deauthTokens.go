package deauthtokens

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

// @Summary Деавторизация пользователя
// @Description Деактивирует Refresh Token пользователя
// @Tags auth
// @Accept json
// @Produce json
//
// @Param Access-Token header string true "Access Token для деактивации"
// @Param Refresh-Token header string true "Refresh Token для деактивации"
//
// @Success 200 {object} nil "Сессия успешно завершена"
// @Failure 401 {object} errorresponse.ErrorResponse "Невалидный Refresh Token"
// @Failure 403 {object} errorresponse.ErrorResponse "Refresh Token не активен"
// @Failure 500 {object} errorresponse.ErrorResponse "Ошибка сервера"
//
// @Router /deauthTokens [post]
func DATHandler(log *slog.Logger, dau deauthUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.DeauthTokens"
		log := log.With(
			slog.String("fn", fn),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		refreshToken := r.Header.Get("Refresh-Token")

		refreshID, err := refresh.Check(refreshToken)
		if err != nil {
			log.Error("failed to check refresh token", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusUnauthorized, "Something wrong")

			return
		}

		err = dau.DeactivateRefreshToken(r.Context(), refreshID)
		if err != nil {
			log.Error("failed to check refresh token", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		userEmail, ok := r.Context().Value(middlewares.UserEmailKey).(string)
		if !ok {
			log.Error("user email not found")
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		log.Info(
			"User deauth",
			slog.String("userEmail", userEmail),
			slog.String("refreshToken", refreshToken),
		)

		w.WriteHeader(http.StatusOK)
	}
}
