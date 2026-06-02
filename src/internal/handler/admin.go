package handler

import (
	"encoding/json"
	"net/http"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AdminHandler struct {
	userService service.UserService
}

func NewAdminHandler(userService service.UserService) *AdminHandler {
	return &AdminHandler{userService: userService}
}

type updateRoleRequest struct {
	Role domain.Role `json:"role" validate:"required"`
}

// ListUsers godoc
// @Summary      List users (admin)
// @Tags         admin
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  userResponse
// @Router       /admin/users [get]
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	resp := make([]userResponse, len(users))
	for i := range users {
		resp[i] = toUserResponse(&users[i])
	}
	writeJSON(w, http.StatusOK, resp)
}

// UpdateUserRole godoc
// @Summary      Update user role (admin)
// @Tags         admin
// @Security     BearerAuth
// @Accept       json
// @Param        id    path  string  true  "User ID"
// @Param        body  body  updateRoleRequest  true  "payload"
// @Success      204
// @Router       /admin/users/{id}/role [patch]
func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}
	var req updateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	if !req.Role.IsValid() {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: domain.ErrInvalidRole.Error()})
		return
	}
	if err := h.userService.UpdateRole(r.Context(), id, req.Role); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser godoc
// @Summary      Delete user (admin)
// @Tags         admin
// @Security     BearerAuth
// @Param        id  path  string  true  "User ID"
// @Success      204
// @Router       /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}
	if err := h.userService.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
