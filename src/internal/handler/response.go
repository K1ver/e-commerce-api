package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/K1ver/e-commerce-api/internal/domain"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrForbidden):
		writeJSON(w, http.StatusForbidden, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrProductNotFound),
		errors.Is(err, domain.ErrCartNotFound), errors.Is(err, domain.ErrOrderNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrEmailAlreadyExists), errors.Is(err, domain.ErrUsernameAlreadyExists):
		writeJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrInsufficientStock), errors.Is(err, domain.ErrInvalidRole):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		if err != nil && err.Error() != "" && err.Error() != "internal server error" {
			// expose simple business errors (cart empty, payment not found, etc.)
			if len(err.Error()) < 120 {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
				return
			}
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}
