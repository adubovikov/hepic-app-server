package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTAuth middleware for JWT token validation
func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Authorization token not provided",
				})
			}

			// Check Bearer token format
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Invalid authorization token format",
				})
			}

			tokenString := tokenParts[1]

			// Parse and validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Check signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token signing method")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Invalid authorization token",
				})
			}

			// Check token validity
			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Authorization token is invalid",
				})
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Invalid token format",
				})
			}

			// Save user information in context
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", int(userID))
			}

			if role, ok := claims["role"].(string); ok {
				c.Set("user_role", role)
			}

			return next(c)
		}
	}
}

// AdminOnly middleware for admin rights validation
func AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("user_role")
			if role != "admin" {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"success": false,
					"error":   "Insufficient access rights",
				})
			}
			return next(c)
		}
	}
}
