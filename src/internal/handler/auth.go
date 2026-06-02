package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/service"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{authService: authService, validate: validate}
}

type registerRequest struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string `json:"lastName" validate:"required,min=2,max=50"`
	Username  string `json:"username" validate:"required,min=3,max=30"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=100"`
}

type loginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// Register godoc
// @Summary      Register buyer
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      registerRequest  true  "payload"
// @Success      201   {object}  map[string]interface{}
// @Failure      400,409,500
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	user := &domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
	}

	tokens, err := h.authService.Register(r.Context(), user)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"user":   toUserResponse(user),
		"tokens": tokens,
	})
}

// Login godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      loginRequest  true  "payload"
// @Success      200   {object}  jwtmanager.TokenPair
// @Failure      400,401,500
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

// Refresh godoc
// @Summary      Refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      refreshRequest  true  "payload"
// @Success      200   {object}  jwtmanager.TokenPair
// @Failure      400,401,500
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid refresh token"})
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}
