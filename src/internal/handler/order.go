package handler

import (
	"encoding/json"
	"net/http"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/pkg/ctxkey"
	"github.com/K1ver/e-commerce-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// Checkout godoc
// @Summary      Checkout cart to order
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Success      201  {object}  domain.Order
// @Router       /orders/checkout [post]
func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	order, err := h.orderService.Checkout(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, order)
}

// ListMine godoc
// @Summary      List my orders
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  domain.Order
// @Router       /orders [get]
func (h *OrderHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	orders, err := h.orderService.GetAllByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, orders)
}

// Get godoc
// @Summary      Get order by id
// @Tags         orders
// @Security     BearerAuth
// @Param        id  path  string  true  "Order ID"
// @Produce      json
// @Success      200  {object}  domain.Order
// @Router       /orders/{id} [get]
func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	role, _ := ctxkey.RoleFrom(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}
	order, err := h.orderService.GetByID(r.Context(), userID, role, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, order)
}

type updateOrderStatusRequest struct {
	Status domain.OrderStatus `json:"status" validate:"required"`
}

// UpdateStatus godoc
// @Summary      Update order status (admin)
// @Tags         orders
// @Security     BearerAuth
// @Accept       json
// @Param        id    path  string  true  "Order ID"
// @Param        body  body  updateOrderStatusRequest  true  "payload"
// @Success      204
// @Router       /orders/{id}/status [patch]
func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}
	var req updateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	if err := h.orderService.UpdateStatusByID(r.Context(), req.Status, id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
