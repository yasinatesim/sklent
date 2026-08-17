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
	"github.com/yasinatesim/vela-commerce/api/internal/promotion"
	"github.com/yasinatesim/vela-commerce/api/internal/reservation"
	"github.com/yasinatesim/vela-commerce/api/internal/returnreq"
	"github.com/yasinatesim/vela-commerce/api/internal/review"
	"github.com/yasinatesim/vela-commerce/api/internal/server/middleware"
	"github.com/yasinatesim/vela-commerce/api/internal/stocktracking"
)

type mailerAdapter struct{ svc *email.Service }

func (m mailerAdapter) SendOrderConfirmationAsync(orderID, to string, totalCents int64) {
	m.svc.SendOrderConfirmationAsync(email.OrderSummary{ID: orderID, Email: to, TotalCents: totalCents})
}

func (m mailerAdapter) SendPasswordResetAsync(to, resetURL string) {
	m.svc.SendPasswordResetAsync(to, resetURL)
}

func (m mailerAdapter) SendLowStockAlertAsync(adminEmail, productTitle string, currentStock int) {
	m.svc.SendLowStockAlertAsync(adminEmail, productTitle, currentStock)
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

	mailSvc := email.NewService(email.LogSender{Log: log}, log)
	mailer := mailerAdapter{svc: mailSvc}

	productRepo := product.NewRepo(db)
	authRepo := auth.NewRepo(db)
	authHandler := auth.NewHandler(signer, authRepo, cfg.CookieDomain, secure, mailer, cfg.FrontendBaseURL)

	productHandler := product.NewHandler(productRepo)
	categoryHandler := category.NewHandler(category.NewRepo(db))

	reservations := reservation.NewService(db, productRepo, mailer, cfg.AdminEmail)
	orderRepo := order.NewRepo(db)
	orderHandler := order.NewHandler(orderRepo, reservations, mailer)

	iyzicoHandler := iyzico.NewHandler(orderRepo, reservations, sandboxVerifier{}, cfg.FrontendBaseURL)

	promotionHandler := promotion.NewHandler(promotion.NewRepo(db))
	stockTrackingHandler := stocktracking.NewHandler(stocktracking.NewRepo(db))
	reviewHandler := review.NewHandler(review.NewRepo(db), productRepo)

	returnRepo := returnreq.NewRepository(db)
	returnHandler := returnreq.NewHandler(returnreq.NewService(returnRepo, orderRepo), returnRepo)

	r.GET("/healthz", health.Handler)

	a := r.Group("/auth")
	{
		a.POST("/register", auth.RateLimit(auth.Rule{Burst: 5, Interval: 3 * time.Minute}), authHandler.Register)
		a.POST("/login", auth.RateLimit(auth.Rule{Burst: 10, Interval: 6 * time.Second}), authHandler.Login)
		a.POST("/logout", authHandler.Logout)
		a.POST("/refresh", authHandler.Refresh)
		a.GET("/me", authHandler.Me)
		a.POST("/password/forgot", auth.RateLimit(auth.Rule{Burst: 5, Interval: 3 * time.Minute}), authHandler.ForgotPassword)
		a.POST("/password/reset", auth.RateLimit(auth.Rule{Burst: 5, Interval: 3 * time.Minute}), authHandler.ResetPassword)
	}

	r.GET("/products", productHandler.List)
	r.GET("/products/facets", productHandler.Facets)
	r.GET("/products/:slug", productHandler.GetBySlug)
	r.GET("/products/:slug/reviews", reviewHandler.List)
	r.POST("/products/:slug/reviews", auth.RequireCSRF(), reviewHandler.Create)
	r.GET("/categories", categoryHandler.List)
	r.POST("/returns", auth.RequireCSRF(), returnHandler.Create)

	registerOrderRoutes(r, orderHandler, signer)

	pay := r.Group("/payments/iyzico")
	{
		pay.POST("/callback", iyzicoHandler.Callback)
	}

	admin := r.Group("/admin", signer.RequireAdmin(), auth.RequireCSRF())
	{
		admin.POST("/products", productHandler.AdminCreate)

		admin.GET("/returns", returnHandler.AdminList)
		admin.PATCH("/returns/:id", returnHandler.AdminUpdateStatus)
		admin.GET("/promotions", promotionHandler.AdminListPromotions)
		admin.POST("/promotions", promotionHandler.AdminCreatePromotion)
		admin.PATCH("/promotions/:id", promotionHandler.AdminUpdatePromotion)
		admin.DELETE("/promotions/:id", promotionHandler.AdminDeletePromotion)

		admin.GET("/coupons", promotionHandler.AdminListCoupons)
		admin.POST("/coupons", promotionHandler.AdminCreateCoupon)
		admin.PATCH("/coupons/:id", promotionHandler.AdminUpdateCoupon)
		admin.DELETE("/coupons/:id", promotionHandler.AdminDeleteCoupon)

		admin.GET("/orders", orderHandler.AdminList)
		admin.PATCH("/orders/:id/status", orderHandler.AdminUpdateStatus)

		admin.GET("/stock-tracking", stockTrackingHandler.AdminList)
		admin.POST("/stock-tracking", stockTrackingHandler.AdminCreate)
		admin.PATCH("/stock-tracking/:id", stockTrackingHandler.AdminUpdate)
		admin.DELETE("/stock-tracking/:id", stockTrackingHandler.AdminDelete)

		admin.GET("/reviews", reviewHandler.AdminList)
		admin.PATCH("/reviews/:id/status", reviewHandler.AdminUpdateStatus)
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
