package routes

import (
	"hepic-app-server/v2/database"
	"hepic-app-server/v2/handlers"
	"hepic-app-server/v2/middleware"
	"hepic-app-server/v2/services"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// SetupRoutes configures all API routes
func SetupRoutes(e *echo.Echo, clickhouse *database.ClickHouseDB, jwtSecret string) {
	// Initialize services
	analyticsService := services.NewAnalyticsService(clickhouse)
	authService := services.NewAuthService(clickhouse, jwtSecret, 24) // 24 hours JWT expiry

	// Initialize handlers
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	authHandler := handlers.NewAuthHandler(authService)

	// Public routes group (no authentication required)
	public := e.Group("/api/v1")
	{
		// Health check
		public.GET("/health", func(c echo.Context) error {
			return c.JSON(200, map[string]interface{}{
				"status":  "ok",
				"message": "HEPIC App Server v2 is running",
			})
		})

		// Swagger documentation
		public.GET("/docs/*", echoSwagger.WrapHandler)
	}

	// Authentication group (public routes)
	auth := e.Group("/api/v1/auth")
	{
		// Registration and login (no authentication required)
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected authentication routes group
	authProtected := e.Group("/api/v1/auth")
	authProtected.Use(middleware.JWT(authService))
	{
		// User profile
		authProtected.GET("/me", authHandler.Me)
		authProtected.PUT("/profile", authHandler.UpdateProfile)
		authProtected.POST("/change-password", authHandler.ChangePassword)
	}

	// Admin routes group
	admin := e.Group("/api/v1/auth")
	admin.Use(middleware.RequireAdmin(authService))
	{
		// User management (admin only)
		admin.GET("/users", authHandler.GetUsers)
		admin.GET("/stats", authHandler.GetUserStats)
	}

	// Analytics routes group
	analytics := e.Group("/api/v1/analytics")
	{
		analytics.GET("/stats", analyticsHandler.GetAnalyticsStats)
		analytics.GET("/protocols", analyticsHandler.GetTopProtocols)
		analytics.GET("/methods", analyticsHandler.GetTopMethods)
		analytics.GET("/traffic", analyticsHandler.GetTrafficByHour)
		analytics.GET("/errors", analyticsHandler.GetErrorRate)
		analytics.GET("/performance", analyticsHandler.GetPerformanceMetrics)
	}
}
