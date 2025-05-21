package checkuserip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	errorresponse "github.com/MTUCIBOY/MedodsTest/pkg/router/errorResponse"
	"github.com/MTUCIBOY/MedodsTest/pkg/router/middlewares"
	"github.com/MTUCIBOY/MedodsTest/pkg/tokens/access"
	"github.com/go-chi/chi/v5/middleware"
)

type webhookResponse struct {
	UserEmail string    `json:"user_email"`
	UserIP    string    `json:"user_ip"`
	JwtUserIP string    `json:"user_jwt_ip"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

func CUIPMiddleware(log *slog.Logger, webhook string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const fn = "middlewares.CheckUserIP"
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

			accessClaims, err := access.CheckWithoutClaims(accessToken, refreshToken)
			if err != nil {
				log.Error("failed to check access token", slog.String("err", err.Error()))
				errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

				return
			}

			userIP := strings.Split(r.RemoteAddr, ":")[0]

			if accessClaims.UserIP != userIP {
				log.Warn(
					"UserIP not match",
					slog.String("jwtUserIP", accessClaims.UserIP),
					slog.String("UserIP", userIP),
				)

				userEmail, ok := r.Context().Value(middlewares.UserEmailKey).(string)
				if !ok {
					log.Error("user email not found")
					errorresponse.JSONResponde(w, http.StatusInternalServerError, "Something wrong")

					return
				}

				resp := webhookResponse{
					UserEmail: userEmail,
					UserIP:    userIP,
					JwtUserIP: accessClaims.UserIP,
					Timestamp: time.Now(),
					Message:   "UserIPs not match",
				}

				err := sendToWebhook(r.Context(), webhook, resp)
				if err != nil {
					log.Warn("failed to send message to webhook", slog.String("err", err.Error()))
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func sendToWebhook(ctx context.Context, webhook string, resp webhookResponse) error {
	body, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	clientResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer clientResp.Body.Close()

	if clientResp.StatusCode < 200 || clientResp.StatusCode >= 300 {
		return fmt.Errorf("bad client status code: %d", clientResp.StatusCode)
	}

	return nil
}
