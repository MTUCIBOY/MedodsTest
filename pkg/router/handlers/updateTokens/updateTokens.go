package updatetokens

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	errorresponse "github.com/MTUCIBOY/MedodsTest/pkg/router/errorResponse"
	"github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/access"
	"github.com/go-chi/chi/v5/middleware"
)

type userResponse struct {
	AccessToken string `json:"access_token"`
}

// @Summary Обновление Access Token
// @Description Генерирует новый Access Token на основе Refresh Token
// @Tags auth
// @Accept json
// @Produce json
//
// @Param Access-Token header string true "Access Token"
// @Param Refresh-Token header string true "Refresh Token"
//
// @Success 200 {object} userResponse "Новый Access Token"
// @Failure 401 {object} errorresponse.ErrorResponse "Невалидный Refresh Token или отсутствие email в контексте"
// @Failure 403 {object} errorresponse.ErrorResponse "Refresh Token не активен"
// @Failure 500 {object} errorresponse.ErrorResponse "Ошибка сервера"
//
// @Router /updateTokens [post]
func UTHandler(log *slog.Logger, ttl time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.UpdateTokens"
		log := log.With(
			slog.String("fn", fn),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		w.Header().Set("Content-Type", "application/json")

		refreshToken := r.Header.Get("Refresh-Token")

		userEmail, ok := r.Context().Value(middlewares.UserEmailKey).(string)
		if !ok {
			log.Error("user email not found")
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		accessToken, err := access.New(
			userEmail,
			r.UserAgent(),
			strings.Split(r.RemoteAddr, ":")[0],
			refreshToken, ttl,
		)
		if err != nil {
			log.Error("failed to generate access token", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		resp := userResponse{AccessToken: accessToken}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			log.Error("failed to encode message", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		log.Info(
			"User update token",
			slog.String("userEmail", userEmail),
			slog.String("accessToken", accessToken),
			slog.String("refreshToken", refreshToken),
		)
		w.WriteHeader(http.StatusOK)
	}
}
