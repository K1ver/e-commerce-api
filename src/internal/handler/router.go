package handler

import (
	"net/http"

	"github.com/K1ver/e-commerce-api/internal/domain"
	"github.com/K1ver/e-commerce-api/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Deps struct {
	AuthMiddleware *middleware.Auth
	Auth           *AuthHandler
	User           *UserHandler
	Product        *ProductHandler
	Cart           *CartHandler
	Order          *OrderHandler
	Payment        *PaymentHandler
	Admin          *AdminHandler
	CORSOrigins    string
}

func NewRouter(deps Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS(deps.CORSOrigins))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", deps.Auth.Register)
			r.Post("/login", deps.Auth.Login)
			r.Post("/refresh", deps.Auth.Refresh)
		})

		r.Route("/products", func(r chi.Router) {
			r.Get("/", deps.Product.List)
			r.Get("/{id}", deps.Product.Get)
		})

		r.Group(func(r chi.Router) {
			r.Use(deps.AuthMiddleware.RequireJWT)

			r.Get("/users/me", deps.User.Me)

			r.Route("/cart", func(r chi.Router) {
				r.With(middleware.RequireRoles(domain.RoleBuyer)).Get("/", deps.Cart.Get)
				r.With(middleware.RequireRoles(domain.RoleBuyer)).Post("/items", deps.Cart.AddItem)
				r.With(middleware.RequireRoles(domain.RoleBuyer)).Put("/items/{productId}", deps.Cart.UpdateItem)
				r.With(middleware.RequireRoles(domain.RoleBuyer)).Delete("/items/{productId}", deps.Cart.RemoveItem)
			})

			r.Route("/orders", func(r chi.Router) {
				r.With(middleware.RequireRoles(domain.RoleBuyer)).Post("/checkout", deps.Order.Checkout)
				r.With(middleware.RequireRoles(domain.RoleBuyer)).Get("/", deps.Order.ListMine)
				r.Get("/{id}", deps.Order.Get)

				r.With(middleware.RequireRoles(domain.RoleBuyer)).Post("/{orderId}/payments", deps.Payment.Create)
				r.Get("/{orderId}/payments", deps.Payment.GetByOrder)

				r.With(middleware.RequireRoles(domain.RoleAdmin)).Patch("/{id}/status", deps.Order.UpdateStatus)
			})

			r.With(middleware.RequireRoles(domain.RoleBuyer)).Post("/payments/{id}/sync", deps.Payment.Sync)

			r.Route("/products", func(r chi.Router) {
				r.With(middleware.RequireRoles(domain.RoleSeller, domain.RoleAdmin)).Post("/", deps.Product.Create)
				r.With(middleware.RequireRoles(domain.RoleSeller, domain.RoleAdmin)).Put("/{id}", deps.Product.Update)
				r.With(middleware.RequireRoles(domain.RoleSeller, domain.RoleAdmin)).Delete("/{id}", deps.Product.Delete)
			})

			r.With(middleware.RequireRoles(domain.RoleSeller, domain.RoleAdmin)).Get("/seller/products", deps.Product.ListMine)

			r.Route("/admin", func(r chi.Router) {
				r.Use(middleware.RequireRoles(domain.RoleAdmin))
				r.Get("/users", deps.Admin.ListUsers)
				r.Patch("/users/{id}/role", deps.Admin.UpdateUserRole)
				r.Delete("/users/{id}", deps.Admin.DeleteUser)
			})
		})
	})

	return r
}
