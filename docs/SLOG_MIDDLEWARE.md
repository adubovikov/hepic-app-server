# Echo-Slog Middleware Documentation

## Overview

This middleware provides structured logging for Echo web framework using Go's standard `log/slog` package. It offers comprehensive request/response logging, error tracking, and panic recovery with structured JSON output.

## Features

- ✅ **Structured JSON Logging** - All logs in JSON format for easy parsing
- ✅ **Request/Response Logging** - Automatic logging of HTTP requests and responses
- ✅ **Error Tracking** - Detailed error logging with context
- ✅ **Panic Recovery** - Automatic panic recovery with logging
- ✅ **Customizable Fields** - Add custom fields to logs
- ✅ **Performance Metrics** - Request duration and size tracking
- ✅ **Request ID Support** - Automatic request ID inclusion

## Usage

### Basic Setup

```go
import (
    "log/slog"
    "os"
    "hepic-app-server/v2/middleware"
    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    
    // Setup slog logger
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)
    
    // Add slog middleware
    e.Use(middleware.Slog())
    e.Use(middleware.SlogError())
    e.Use(middleware.SlogRecover())
    
    // Your routes...
}
```

### Advanced Configuration

```go
// Custom slog configuration
config := middleware.SlogConfig{
    Skipper: func(c echo.Context) bool {
        // Skip logging for health checks
        return c.Request().URL.Path == "/health"
    },
    Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
        AddSource: true,
    })),
    IncludeRequestID:   true,
    IncludeUserAgent:   true,
    IncludeRemoteAddr:  true,
    CustomFields: func(c echo.Context) []slog.Attr {
        return []slog.Attr{
            slog.String("user_id", getUserID(c)),
            slog.String("session_id", getSessionID(c)),
        }
    },
}

e.Use(middleware.SlogWithConfig(config))
```

## Middleware Components

### 1. Slog Middleware

Logs HTTP requests and responses with timing information.

**Features:**
- Request method, path, and remote address
- Response status code and duration
- Request size tracking
- Optional fields (User-Agent, Request ID, Remote Address)

**Example Output:**
```json
{
  "time": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "msg": "Request started",
  "method": "GET",
  "path": "/api/v1/analytics/stats",
  "remote_addr": "192.168.1.100:12345",
  "request_id": "req-123456",
  "user_agent": "Mozilla/5.0..."
}
```

### 2. SlogError Middleware

Logs errors with detailed context.

**Features:**
- Error message and stack trace
- Request context (method, path, status)
- Request ID correlation
- Error severity levels

**Example Output:**
```json
{
  "time": "2024-01-15T10:30:00Z",
  "level": "ERROR",
  "msg": "Request error",
  "method": "GET",
  "path": "/api/v1/analytics/stats",
  "status": 500,
  "error": "database connection failed",
  "request_id": "req-123456"
}
```

### 3. SlogRecover Middleware

Recovers from panics and logs them.

**Features:**
- Panic recovery with logging
- Stack trace capture
- Automatic 500 error response
- Request context preservation

**Example Output:**
```json
{
  "time": "2024-01-15T10:30:00Z",
  "level": "ERROR",
  "msg": "Panic recovered",
  "method": "GET",
  "path": "/api/v1/analytics/stats",
  "panic": "runtime error: index out of range",
  "request_id": "req-123456"
}
```

## Configuration Options

### SlogConfig

```go
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
```

### Skipper Function

```go
type Skipper func(c echo.Context) bool

// Example: Skip logging for health checks
func skipHealthChecks(c echo.Context) bool {
    return c.Request().URL.Path == "/health"
}
```

## Log Levels

The middleware uses different log levels based on response status:

- **INFO** - Status 200-299 (Success)
- **WARN** - Status 300-399 (Redirects)
- **ERROR** - Status 400+ (Client/Server Errors)

## Custom Fields

Add custom fields to logs using the `CustomFields` function:

```go
config := middleware.SlogConfig{
    CustomFields: func(c echo.Context) []slog.Attr {
        return []slog.Attr{
            slog.String("user_id", getUserID(c)),
            slog.String("tenant_id", getTenantID(c)),
            slog.String("api_version", c.Request().Header.Get("API-Version")),
        }
    },
}
```

## Performance Considerations

- **Minimal Overhead** - Structured logging with minimal performance impact
- **Async Logging** - Consider using async handlers for high-traffic applications
- **Log Rotation** - Implement log rotation for production environments
- **Sampling** - Use sampling for high-volume endpoints

## Production Setup

### 1. Log Rotation

```go
import "github.com/natefinch/lumberjack"

// Setup log rotation
logFile := &lumberjack.Logger{
    Filename:   "/var/log/hepic-app-server/app.log",
    MaxSize:    100, // MB
    MaxBackups: 3,
    MaxAge:     28, // days
    Compress:   true,
}

logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))
```

### 2. Environment-based Configuration

```go
func setupLogger() *slog.Logger {
    var handler slog.Handler
    
    if os.Getenv("ENV") == "production" {
        // JSON logging for production
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelInfo,
        })
    } else {
        // Text logging for development
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelDebug,
        })
    }
    
    return slog.New(handler)
}
```

### 3. Log Aggregation

For production environments, consider using log aggregation tools:

- **ELK Stack** (Elasticsearch, Logstash, Kibana)
- **Fluentd** with ClickHouse
- **Prometheus** with Grafana
- **Cloud Logging** (AWS CloudWatch, Google Cloud Logging)

## Examples

### Complete Setup

```go
package main

import (
    "log/slog"
    "os"
    "hepic-app-server/v2/middleware"
    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    
    // Setup structured logger
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
        AddSource: false,
    }))
    slog.SetDefault(logger)
    
    // Add middleware
    e.Use(middleware.Slog())
    e.Use(middleware.SlogError())
    e.Use(middleware.SlogRecover())
    
    // Routes
    e.GET("/api/v1/analytics/stats", getAnalyticsStats)
    
    e.Logger.Fatal(e.Start(":8080"))
}
```

### Custom Logging in Handlers

```go
func getAnalyticsStats(c echo.Context) error {
    slog.Info("Analytics stats request",
        "method", c.Request().Method,
        "path", c.Request().URL.Path,
        "user_id", getUserID(c),
    )
    
    // Your handler logic...
    
    slog.Info("Analytics stats retrieved",
        "total_records", stats["total_records"],
        "duration", time.Since(start),
    )
    
    return c.JSON(200, stats)
}
```

## Troubleshooting

### Common Issues

1. **Missing Request ID**: Ensure `middleware.RequestID()` is added before slog middleware
2. **High Log Volume**: Use sampling or filtering for high-traffic endpoints
3. **Performance Impact**: Consider async logging for production
4. **Log Parsing**: Use structured JSON for easier parsing and analysis

### Debug Mode

Enable debug logging for development:

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
    AddSource: true,
}))
```

## Integration with Monitoring

### Prometheus Metrics

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
)
```

### Health Checks

```go
e.GET("/health", func(c echo.Context) error {
    slog.Info("Health check requested",
        "remote_addr", c.Request().RemoteAddr,
    )
    
    return c.JSON(200, map[string]string{
        "status": "ok",
        "timestamp": time.Now().Format(time.RFC3339),
    })
})
```

This middleware provides comprehensive logging capabilities for Echo applications with structured, searchable, and analyzable log output.
