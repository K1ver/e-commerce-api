package handler

import (
	"encoding/json"
	"net/http"

	"github.com/K1ver/e-commerce-api/internal/pkg/ctxkey"
	"github.com/K1ver/e-commerce-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CartHandler struct {
	cartService service.CartService
	validate    *validator.Validate
}

func NewCartHandler(cartService service.CartService, validate *validator.Validate) *CartHandler {
	return &CartHandler{cartService: cartService, validate: validate}
}

type cartItemRequest struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

// Get godoc
// @Summary      Get my cart
// @Tags         cart
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  domain.Cart
// @Router       /cart [get]
func (h *CartHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	cart, err := h.cartService.GetMine(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cart)
}

// AddItem godoc
// @Summary      Add item to cart
// @Tags         cart
// @Security     BearerAuth
// @Accept       json
// @Success      204
// @Router       /cart/items [post]
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	var req cartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	if err := h.cartService.AddItem(r.Context(), userID, req.ProductID, req.Quantity); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateItem godoc
// @Summary      Update cart item quantity
// @Tags         cart
// @Security     BearerAuth
// @Accept       json
// @Param        productId  path  string  true  "Product ID"
// @Success      204
// @Router       /cart/items/{productId} [put]
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	productID, err := uuid.Parse(chi.URLParam(r, "productId"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid product id"})
		return
	}
	var req struct {
		Quantity int `json:"quantity" validate:"required,min=1"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	if err := h.cartService.UpdateItem(r.Context(), userID, productID, req.Quantity); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RemoveItem godoc
// @Summary      Remove item from cart
// @Tags         cart
// @Security     BearerAuth
// @Param        productId  path  string  true  "Product ID"
// @Success      204
// @Router       /cart/items/{productId} [delete]
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	productID, err := uuid.Parse(chi.URLParam(r, "productId"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid product id"})
		return
	}
	if err := h.cartService.RemoveItem(r.Context(), userID, productID); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
