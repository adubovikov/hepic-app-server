package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTAuth middleware для проверки JWT токенов
func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Получение токена из заголовка Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Токен авторизации не предоставлен",
				})
			}

			// Проверка формата Bearer token
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Неверный формат токена авторизации",
				})
			}

			tokenString := tokenParts[1]

			// Парсинг и валидация токена
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Проверка метода подписи
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Неверный метод подписи токена")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Неверный токен авторизации",
				})
			}

			// Проверка валидности токена
			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Токен авторизации недействителен",
				})
			}

			// Извлечение claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Неверный формат токена",
				})
			}

			// Сохранение информации о пользователе в контексте
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

// AdminOnly middleware для проверки прав администратора
func AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("user_role")
			if role != "admin" {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"success": false,
					"error":   "Недостаточно прав доступа",
				})
			}
			return next(c)
		}
	}
}
