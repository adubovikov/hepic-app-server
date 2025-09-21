package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORSConfig CORS settings
func CORSConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // In production specify specific domains
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
			echo.PATCH,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXRequestedWith,
		},
		ExposeHeaders: []string{
			echo.HeaderContentLength,
			echo.HeaderContentType,
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}
