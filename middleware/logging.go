package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CustomLoggerConfig настройки кастомного логгера
func CustomLoggerConfig() middleware.LoggerConfig {
	return middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}","status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}","bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: time.RFC3339,
	}
}

// RequestLogger middleware для логирования запросов
func RequestLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(CustomLoggerConfig())
}
