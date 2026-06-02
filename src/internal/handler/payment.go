package handler

import (
	"net/http"

	"github.com/K1ver/e-commerce-api/internal/pkg/ctxkey"
	"github.com/K1ver/e-commerce-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// Create godoc
// @Summary      Create payment for order
// @Tags         payments
// @Security     BearerAuth
// @Param        orderId  path  string  true  "Order ID"
// @Produce      json
// @Success      201  {object}  map[string]interface{}
// @Router       /orders/{orderId}/payments [post]
func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	role, _ := ctxkey.RoleFrom(r.Context())
	orderID, err := uuid.Parse(chi.URLParam(r, "orderId"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid order id"})
		return
	}
	payment, url, err := h.paymentService.CreateForOrder(r.Context(), userID, role, orderID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"payment":         payment,
		"confirmationUrl": url,
	})
}

// GetByOrder godoc
// @Summary      Get payment by order
// @Tags         payments
// @Security     BearerAuth
// @Param        orderId  path  string  true  "Order ID"
// @Produce      json
// @Success      200  {object}  domain.Payment
// @Router       /orders/{orderId}/payments [get]
func (h *PaymentHandler) GetByOrder(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	role, _ := ctxkey.RoleFrom(r.Context())
	orderID, err := uuid.Parse(chi.URLParam(r, "orderId"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid order id"})
		return
	}
	payment, err := h.paymentService.GetByOrderID(r.Context(), userID, role, orderID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, payment)
}

// Sync godoc
// @Summary      Sync payment status
// @Tags         payments
// @Security     BearerAuth
// @Param        id  path  string  true  "Payment ID"
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /payments/{id}/sync [post]
func (h *PaymentHandler) Sync(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	role, _ := ctxkey.RoleFrom(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}
	status, err := h.paymentService.Sync(r.Context(), userID, role, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": string(status)})
}
