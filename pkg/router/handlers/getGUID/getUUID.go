package getguid

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	errorresponse "github.com/MTUCIBOY/MedodsTest/pkg/router/errorResponse"
	"github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares"
	"github.com/go-chi/chi/v5/middleware"
)

type UUID interface {
	UserUUID(ctx context.Context, email string) (string, error)
}

type userResponse struct {
	GUID string `json:"guid"`
}

func UUIDHadler(log *slog.Logger, u UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.GetUUID"
		log := log.With(
			slog.String("fn", fn),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		w.Header().Set("Content-Type", "application/json")

		userEmail, ok := r.Context().Value(middlewares.UserEmailKey).(string)
		if !ok {
			log.Error("user email not found")
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		uuid, err := u.UserUUID(r.Context(), userEmail)
		if err != nil {
			log.Error("error in UserUUID", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		resp := userResponse{GUID: uuid}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			log.Error("failed to encode message", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		log.Info("User get GUID", slog.String("userEmail", userEmail))
		w.WriteHeader(http.StatusOK)
	}
}
