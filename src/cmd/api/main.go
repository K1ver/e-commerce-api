// @title           E-Commerce API
// @version         1.0
// @description     E-Commerce API with JWT, roles (admin/seller/buyer), cart, orders and payments.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT access token.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/K1ver/e-commerce-api/docs"
	"github.com/K1ver/e-commerce-api/internal/config"
	"github.com/K1ver/e-commerce-api/internal/handler"
	"github.com/K1ver/e-commerce-api/internal/middleware"
	jwtmanager "github.com/K1ver/e-commerce-api/internal/pkg/jwt"
	"github.com/K1ver/e-commerce-api/internal/repository/postgres"
	"github.com/K1ver/e-commerce-api/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := postgres.Connect(cfg.Postgres)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer db.Close()

	validate := validator.New()
	jwt := jwtmanager.NewManager(cfg.JWT)

	userRepo := postgres.NewUserRepository(db)
	productRepo := postgres.NewProductRepository(db)
	cartRepo := postgres.NewCartRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)

	userService := service.NewUserService(userRepo, validate)
	authService := service.NewAuthService(userService, jwt)
	productService := service.NewProductService(productRepo, validate)
	cartService := service.NewCartService(cartRepo, productRepo)
	orderService := service.NewOrderService(orderRepo)

	var paymentHandler *yookassa.PaymentHandler
	if cfg.YooKassa.ShopID != "" && cfg.YooKassa.SecretKey != "" {
		yooclient := yookassa.NewClient(cfg.YooKassa.ShopID, cfg.YooKassa.SecretKey)
		paymentHandler = yookassa.NewPaymentHandler(yooclient)
		log.Println("YooKassa payment handler enabled")
	} else {
		log.Println("YooKassa not configured — mock payments enabled")
	}
	paymentService := service.NewPaymentService(paymentRepo, orderRepo, paymentHandler, cfg.YooKassa, validate)

	authMiddleware := middleware.NewAuth(jwt)
	router := handler.NewRouter(handler.Deps{
		AuthMiddleware: authMiddleware,
		Auth:           handler.NewAuthHandler(authService, validate),
		User:           handler.NewUserHandler(userService),
		Product:        handler.NewProductHandler(productService, validate),
		Cart:           handler.NewCartHandler(cartService, validate),
		Order:          handler.NewOrderHandler(orderService),
		Payment:        handler.NewPaymentHandler(paymentService),
		Admin:          handler.NewAdminHandler(userService),
		CORSOrigins:    cfg.Cors.AllowOrigins,
	})

	addr := ":" + cfg.Server.InternalPort
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s (swagger: http://localhost%s/swagger/index.html)", addr, addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
