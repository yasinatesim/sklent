package server

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yasinatesim/vela-commerce/api/internal/auth"
	"github.com/yasinatesim/vela-commerce/api/internal/category"
	"github.com/yasinatesim/vela-commerce/api/internal/config"
	"github.com/yasinatesim/vela-commerce/api/internal/email"
	"github.com/yasinatesim/vela-commerce/api/internal/health"
	"github.com/yasinatesim/vela-commerce/api/internal/order"
	"github.com/yasinatesim/vela-commerce/api/internal/payment/iyzico"
	"github.com/yasinatesim/vela-commerce/api/internal/product"
	"github.com/yasinatesim/vela-commerce/api/internal/reservation"
	"github.com/yasinatesim/vela-commerce/api/internal/server/middleware"
)

type mailerAdapter struct{ svc *email.Service }

func (m mailerAdapter) SendOrderConfirmationAsync(orderID, to string, totalCents int64) {
	m.svc.SendOrderConfirmationAsync(email.OrderSummary{ID: orderID, Email: to, TotalCents: totalCents})
}

func New(cfg config.Config, db *gorm.DB, log *slog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(gin.Logger())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigin))

	signer := auth.NewSigner(cfg.JWTSecret)
	secure := cfg.IsProduction()
	r.Use(auth.IssueCSRF(cfg.CookieDomain, secure))

	authRepo := auth.NewRepo(db)
	authHandler := auth.NewHandler(signer, authRepo, cfg.CookieDomain, secure)

	productHandler := product.NewHandler(product.NewRepo(db))
	categoryHandler := category.NewHandler(category.NewRepo(db))

	reservations := reservation.NewService(db)
	mailSvc := email.NewService(email.LogSender{Log: log}, log)
	orderRepo := order.NewRepo(db)
	orderHandler := order.NewHandler(orderRepo, reservations, mailerAdapter{svc: mailSvc})

	iyzicoHandler := iyzico.NewHandler(orderRepo, reservations, sandboxVerifier{}, cfg.FrontendBaseURL)

	r.GET("/healthz", health.Handler)

	a := r.Group("/auth")
	{
		a.POST("/register", auth.RateLimit(auth.Rule{Burst: 5, Interval: 3 * time.Minute}), authHandler.Register)
		a.POST("/login", auth.RateLimit(auth.Rule{Burst: 10, Interval: 6 * time.Second}), authHandler.Login)
		a.POST("/logout", authHandler.Logout)
		a.POST("/refresh", authHandler.Refresh)
		a.GET("/me", authHandler.Me)
	}

	r.GET("/products", productHandler.List)
	r.GET("/products/:slug", productHandler.GetBySlug)
	r.GET("/categories", categoryHandler.List)

	registerOrderRoutes(r, orderHandler, signer)

	pay := r.Group("/payments/iyzico")
	{
		pay.POST("/callback", iyzicoHandler.Callback)
	}

	admin := r.Group("/admin", signer.RequireAdmin(), auth.RequireCSRF())
	{
		admin.POST("/products", productHandler.AdminCreate)
	}

	return r
}

func registerOrderRoutes(r *gin.Engine, h *order.Handler, signer *auth.Signer) {
	trackRule := auth.Rule{Burst: 20, Interval: time.Minute}
	og := r.Group("/orders")
	og.POST("", auth.RequireCSRF(), signer.OptionalAuth(), h.Place)
	og.GET("/track/:token", auth.RateLimit(trackRule), h.GetByGuestToken)
	og.GET("", signer.RequireAuth(), h.ListForUser)
	og.GET("/:id", signer.RequireAuth(), h.GetByID)
}
