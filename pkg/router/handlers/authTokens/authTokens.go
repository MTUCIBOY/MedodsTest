package authtokens

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errorrespondes "github.com/MTUCIBOY/MedodsTest/pkg/router/errorRespondes"
	"github.com/MTUCIBOY/MedodsTest/pkg/storage"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/access"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/refresh"
	"github.com/go-chi/chi/v5/middleware"
)

type checkUser interface {
	Auth(ctx context.Context, email, password string) (bool, error)
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ATHandler(log *slog.Logger, cu checkUser) http.HandlerFunc {
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
			errorrespondes.JSONResponde(w, http.StatusBadRequest, "Invalid request body")

			return
		}

		isAuth, err := cu.Auth(r.Context(), req.Email, req.Password)
		if err != nil && !errors.Is(err, storage.ErrWrongPassword) {
			log.Error("error in DB Auth", slog.String("err", err.Error()))

			if errors.Is(err, storage.ErrEmailNotFound) {
				errorrespondes.JSONResponde(w, http.StatusBadRequest, "Email does not exists")

				return
			}

			errorrespondes.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		if !isAuth {
			log.Error("auth failed")
			errorrespondes.JSONResponde(w, http.StatusUnauthorized, "Wrong email or password")

			return
		}

		accessToken, err := access.New(req.Email, w.Header().Get("User-Agent"))
		if err != nil {
			log.Error("failed to make access token", slog.String("err", err.Error()))
			errorrespondes.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		refreshToken, err := refresh.New(accessToken)
		if err != nil {
			log.Error("failed to make refresh token", slog.String("err", err.Error()))
			errorrespondes.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		resp := userResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			log.Error("failed to encode message", slog.String("err", err.Error()))
			errorrespondes.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		log.Info("User get tokens")
		w.WriteHeader(http.StatusOK)
	}
}
