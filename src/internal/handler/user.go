package handler

import (
	"net/http"

	"github.com/K1ver/e-commerce-api/internal/pkg/ctxkey"
	"github.com/K1ver/e-commerce-api/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := ctxkey.UserIDFrom(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	user, err := h.userService.GetById(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(user))
}
