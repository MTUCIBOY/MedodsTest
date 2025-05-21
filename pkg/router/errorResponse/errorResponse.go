package errorresponse

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func JSONResponde(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	ER := ErrorResponse{
		Code:    status,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(ER); err != nil {
		slog.Error("failed to send message", slog.String("err", err.Error()))
	}
}
