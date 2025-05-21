package newuser

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errorresponse "github.com/MTUCIBOY/MedodsTest/pkg/router/errorResponse"
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

// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя по email и паролю
// @Tags user
// @Accept json
// @Produce json
//
// @Param request body userRequest true "Данные для регистрации (email и пароль)"
//
// @Success 201 {object} nil "Пользователь успешно зарегистрирован"
// @Failure 400 {object} errorresponse.ErrorResponse "Неверный запрос или email уже существует"
// @Failure 401 {object} errorresponse.ErrorResponse "Невалидные токены"
// @Failure 500 {object} errorresponse.ErrorResponse "Ошибка сервера"
//
// @Router /registrate [post]
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
			errorresponse.JSONResponde(w, http.StatusBadRequest, "Invalid request body")

			return
		}

		if err := ru.AddUser(r.Context(), req.Email, req.Password); err != nil {
			if errors.Is(err, storage.ErrEmailExist) {
				log.Error(err.Error())
				errorresponse.JSONResponde(w, http.StatusBadRequest, "Email alreary exists")

				return
			}

			log.Error("error in AddUser DB", slog.String("err", err.Error()))
			errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

			return
		}

		log.Info("User added in DB", slog.String("email", req.Email))
		w.WriteHeader(http.StatusCreated)
	}
}
