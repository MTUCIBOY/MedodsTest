package authtokens

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	errorresponse "github.com/MTUCIBOY/MedodsTest/pkg/router/errorResponse"
	"github.com/MTUCIBOY/MedodsTest/pkg/storage"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/access"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/refresh"
	"github.com/go-chi/chi/v5/middleware"
)

type checkUser interface {
	Auth(ctx context.Context, email, password string) (bool, error)
	AddRefreshToken(ctx context.Context, email, token, uuidToken string) error
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ATHandler(log *slog.Logger, ttl time.Duration, cu checkUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.authTokens"
		log = log.With(
			slog.String("fn", fn),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		w.Header().Set("Content-Type", "application/json")

		var req userRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusBadRequest, "Invalid request body")

			return
		}

		isAuth, err := cu.Auth(r.Context(), req.Email, req.Password)
		if err != nil && !errors.Is(err, storage.ErrWrongPassword) {
			log.Error("error in DB Auth", slog.String("err", err.Error()))

			if errors.Is(err, storage.ErrEmailNotFound) {
				errorresponse.JSONResponde(w, http.StatusBadRequest, "Email does not exists")

				return
			}

			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		if !isAuth {
			log.Error("auth failed")
			errorresponse.JSONResponde(w, http.StatusUnauthorized, "Wrong email or password")

			return
		}

		refreshToken, err := refresh.New()
		if err != nil {
			log.Error("failed to make refresh token", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		accessToken, err := access.New(
			req.Email, r.UserAgent(), refreshToken, ttl,
		)
		if err != nil {
			log.Error("failed to make access token", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		uuidToken, err := refresh.Check(refreshToken)
		if err != nil {
			log.Error("failed to get refresh ID", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		err = cu.AddRefreshToken(r.Context(), req.Email, refreshToken, uuidToken)
		if err != nil {
			log.Error("failed to write in DB refresh token", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		resp := userResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			log.Error("failed to encode message", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		log.Info(
			"User get tokens",
			slog.String("userEmail", req.Email),
			slog.String("accessToken", accessToken),
			slog.String("refreshToken", refreshToken),
		)
		w.WriteHeader(http.StatusOK)
	}
}
