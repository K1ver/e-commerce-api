package handler

import (
	"encoding/json"
	"net/http"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/pkg/ctxkey"
	"github.com/K1ver/e-commerce-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService service.ProductService
	validate       *validator.Validate
}

func NewProductHandler(productService service.ProductService, validate *validator.Validate) *ProductHandler {
	return &ProductHandler{productService: productService, validate: validate}
}

type productRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"required,min=2,max=2000"`
	Price       int64  `json:"price" validate:"required,gt=0"`
	Stock       int    `json:"stock" validate:"gte=0"`
}

// List godoc
// @Summary      List products
// @Tags         products
// @Produce      json
// @Success      200  {array}   productResponse
// @Router       /products [get]
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.FindAll(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProductsResponse(products))
}

// Get godoc
// @Summary      Get product
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  productResponse
// @Router       /products/{id} [get]
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}
	product, err := h.productService.FindByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProductResponse(product))
}

// Create godoc
// @Summary      Create product (seller/admin)
// @Tags         products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  productRequest  true  "payload"
// @Success      201   {object}  productResponse
// @Router       /products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	var req productRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	product := &domain.Product{
		SellerID:    &userID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	if err := h.productService.Create(r.Context(), product); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toProductResponse(*product))
}

// Update godoc
// @Summary      Update product (seller/admin)
// @Tags         products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string  true  "Product ID"
// @Param        body  body  productRequest  true  "payload"
// @Success      200   {object}  productResponse
// @Router       /products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	role, _ := ctxkey.RoleFrom(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}
	var req productRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	product := domain.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	if err := h.productService.Update(r.Context(), userID, role, product); err != nil {
		writeError(w, err)
		return
	}
	updated, err := h.productService.FindByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProductResponse(updated))
}

// Delete godoc
// @Summary      Delete product (seller/admin)
// @Tags         products
// @Security     BearerAuth
// @Param        id  path  string  true  "Product ID"
// @Success      204
// @Router       /products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	role, _ := ctxkey.RoleFrom(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}
	if err := h.productService.Delete(r.Context(), userID, role, id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListMine godoc
// @Summary      List seller products
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  productResponse
// @Router       /seller/products [get]
func (h *ProductHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, _ := ctxkey.UserIDFrom(r.Context())
	products, err := h.productService.FindBySeller(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProductsResponse(products))
}
