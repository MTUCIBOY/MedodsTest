package newuser

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errorrespondes "github.com/MTUCIBOY/MedodsTest/pkg/router/errorRespondes"
	"github.com/MTUCIBOY/MedodsTest/pkg/storage"
	"github.com/go-chi/chi/v5/middleware"
)

type registerUser interface {
	AddUser(ctx context.Context, email, password string) error
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NUHandler(log *slog.Logger, ru registerUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.newUser"
		log := log.With(
			slog.String("fn", fn),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		var req userRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			errorrespondes.JSONResponde(w, http.StatusBadRequest, "Invalid request body")

			return
		}

		if err := ru.AddUser(r.Context(), req.Email, req.Password); err != nil {
			if errors.Is(err, storage.ErrEmailExist) {
				log.Error(err.Error())
				errorrespondes.JSONResponde(w, http.StatusBadRequest, "Email alreary exists")

				return
			}

			log.Error("error in AddUser DB", slog.String("err", err.Error()))
			errorrespondes.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		log.Info("User added in DB")
		w.WriteHeader(http.StatusCreated)
	}
}
