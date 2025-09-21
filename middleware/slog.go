package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Skipper defines a function to skip middleware
type Skipper func(c echo.Context) bool

// DefaultSkipper returns false which processes the middleware
var DefaultSkipper Skipper = func(c echo.Context) bool {
	return false
}

// SlogConfig defines the config for Slog middleware
type SlogConfig struct {
	// Skipper defines a function to skip middleware
	Skipper Skipper
	// Logger is the slog.Logger instance to use
	Logger *slog.Logger
	// IncludeRequestID includes request ID in logs
	IncludeRequestID bool
	// IncludeUserAgent includes user agent in logs
	IncludeUserAgent bool
	// IncludeRemoteAddr includes remote address in logs
	IncludeRemoteAddr bool
	// CustomFields allows adding custom fields to logs
	CustomFields func(c echo.Context) []slog.Attr
}

// DefaultSlogConfig is the default Slog middleware config
var DefaultSlogConfig = SlogConfig{
	Skipper:           DefaultSkipper,
	Logger:            slog.Default(),
	IncludeRequestID:  true,
	IncludeUserAgent:  true,
	IncludeRemoteAddr: true,
	CustomFields:      nil,
}

// Slog returns a middleware that logs HTTP requests using slog
func Slog() echo.MiddlewareFunc {
	return SlogWithConfig(DefaultSlogConfig)
}

// SlogWithConfig returns a Slog middleware with config
func SlogWithConfig(config SlogConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultSkipper
	}
	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()

			// Create context with request info
			ctx := context.WithValue(req.Context(), "echo", c)

			// Prepare log attributes
			attrs := []slog.Attr{
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.String("remote_addr", req.RemoteAddr),
			}

			// Add optional fields
			if config.IncludeRequestID {
				if reqID := req.Header.Get(echo.HeaderXRequestID); reqID != "" {
					attrs = append(attrs, slog.String("request_id", reqID))
				}
			}

			if config.IncludeUserAgent {
				if ua := req.Header.Get("User-Agent"); ua != "" {
					attrs = append(attrs, slog.String("user_agent", ua))
				}
			}

			if config.IncludeRemoteAddr {
				attrs = append(attrs, slog.String("remote_addr", req.RemoteAddr))
			}

			// Add custom fields if provided
			if config.CustomFields != nil {
				customAttrs := config.CustomFields(c)
				attrs = append(attrs, customAttrs...)
			}

			// Log request start
			config.Logger.LogAttrs(ctx, slog.LevelInfo, "Request started", attrs...)

			// Process request
			err := next(c)

			// Calculate duration
			duration := time.Since(start)

			// Prepare response log attributes
			responseAttrs := []slog.Attr{
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.Int("status", res.Status),
				slog.Duration("duration", duration),
				slog.Int64("size", res.Size),
			}

			// Add error if present
			if err != nil {
				responseAttrs = append(responseAttrs, slog.String("error", err.Error()))
			}

			// Log response
			level := slog.LevelInfo
			if res.Status >= 400 {
				level = slog.LevelError
			} else if res.Status >= 300 {
				level = slog.LevelWarn
			}

			config.Logger.LogAttrs(ctx, level, "Request completed", responseAttrs...)

			return err
		}
	}
}

// SlogError returns a middleware that logs errors using slog
func SlogError() echo.MiddlewareFunc {
	return SlogErrorWithConfig(DefaultSlogConfig)
}

// SlogErrorWithConfig returns a Slog error middleware with config
func SlogErrorWithConfig(config SlogConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultSkipper
	}
	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				req := c.Request()
				res := c.Response()

				attrs := []slog.Attr{
					slog.String("method", req.Method),
					slog.String("path", req.URL.Path),
					slog.Int("status", res.Status),
					slog.String("error", err.Error()),
				}

				// Add request ID if available
				if reqID := req.Header.Get(echo.HeaderXRequestID); reqID != "" {
					attrs = append(attrs, slog.String("request_id", reqID))
				}

				config.Logger.LogAttrs(req.Context(), slog.LevelError, "Request error", attrs...)
			}
			return err
		}
	}
}

// SlogRecover returns a middleware that recovers from panics and logs them using slog
func SlogRecover() echo.MiddlewareFunc {
	return SlogRecoverWithConfig(DefaultSlogConfig)
}

// SlogRecoverWithConfig returns a Slog recover middleware with config
func SlogRecoverWithConfig(config SlogConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultSkipper
	}
	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					req := c.Request()

					attrs := []slog.Attr{
						slog.String("method", req.Method),
						slog.String("path", req.URL.Path),
						slog.Any("panic", r),
					}

					// Add request ID if available
					if reqID := req.Header.Get(echo.HeaderXRequestID); reqID != "" {
						attrs = append(attrs, slog.String("request_id", reqID))
					}

					config.Logger.LogAttrs(req.Context(), slog.LevelError, "Panic recovered", attrs...)

					// Return internal server error
					c.Error(echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error"))
				}
			}()

			return next(c)
		}
	}
}
