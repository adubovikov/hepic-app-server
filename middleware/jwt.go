package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"hepic-app-server/v2/services"

	"github.com/labstack/echo/v4"
)

// JWTConfig defines the config for JWT middleware
type JWTConfig struct {
	// Skipper defines a function to skip middleware
	Skipper Skipper
	// AuthService is the authentication service
	AuthService *services.AuthService
	// RequiredRole is the required role for access (optional)
	RequiredRole string
}

// DefaultJWTConfig is the default JWT middleware config
var DefaultJWTConfig = JWTConfig{
	Skipper:      DefaultSkipper,
	AuthService:  nil,
	RequiredRole: "",
}

// JWT returns a middleware that validates JWT tokens
func JWT(authService *services.AuthService) echo.MiddlewareFunc {
	config := DefaultJWTConfig
	config.AuthService = authService
	return JWTWithConfig(config)
}

// JWTWithConfig returns a JWT middleware with config
func JWTWithConfig(config JWTConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				slog.Error("Missing Authorization header",
					"method", c.Request().Method,
					"path", c.Request().URL.Path,
					"remote_addr", c.Request().RemoteAddr,
				)
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Missing Authorization header",
				})
			}

			// Check Bearer token format
			if !strings.HasPrefix(authHeader, "Bearer ") {
				slog.Error("Invalid Authorization header format",
					"method", c.Request().Method,
					"path", c.Request().URL.Path,
					"remote_addr", c.Request().RemoteAddr,
				)
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Invalid Authorization header format",
				})
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				slog.Error("Empty token",
					"method", c.Request().Method,
					"path", c.Request().URL.Path,
					"remote_addr", c.Request().RemoteAddr,
				)
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Empty token",
				})
			}

			// Validate JWT token
			payload, err := config.AuthService.ValidateJWT(token)
			if err != nil {
				slog.Error("Invalid JWT token",
					"error", err,
					"method", c.Request().Method,
					"path", c.Request().URL.Path,
					"remote_addr", c.Request().RemoteAddr,
				)
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Invalid token",
				})
			}

			// Check role if required
			if config.RequiredRole != "" && payload.Role != config.RequiredRole {
				slog.Error("Insufficient permissions",
					"required_role", config.RequiredRole,
					"user_role", payload.Role,
					"user_id", payload.UserID,
					"method", c.Request().Method,
					"path", c.Request().URL.Path,
				)
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"success": false,
					"error":   "Insufficient permissions",
				})
			}

			// Set user information in context
			c.Set("user_id", payload.UserID)
			c.Set("username", payload.Username)
			c.Set("user_role", payload.Role)

			slog.Info("JWT token validated successfully",
				"user_id", payload.UserID,
				"username", payload.Username,
				"role", payload.Role,
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
			)

			return next(c)
		}
	}
}

// RequireAdmin returns a middleware that requires admin role
func RequireAdmin(authService *services.AuthService) echo.MiddlewareFunc {
	config := DefaultJWTConfig
	config.AuthService = authService
	config.RequiredRole = "admin"
	return JWTWithConfig(config)
}

// RequireUser returns a middleware that requires user or admin role
func RequireUser(authService *services.AuthService) echo.MiddlewareFunc {
	config := DefaultJWTConfig
	config.AuthService = authService
	// No specific role required, just valid JWT
	return JWTWithConfig(config)
}
