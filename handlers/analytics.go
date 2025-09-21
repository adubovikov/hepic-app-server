package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"hepic-app-server/v2/models"
	"hepic-app-server/v2/services"

	"github.com/labstack/echo/v4"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetAnalyticsStats godoc
// @Summary Get analytics statistics
// @Description Get comprehensive analytics statistics from ClickHouse
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Param start_date query string false "Start date (RFC3339)"
// @Param end_date query string false "End date (RFC3339)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/analytics/stats [get]
func (h *AnalyticsHandler) GetAnalyticsStats(c echo.Context) error {
	slog.Info("Analytics stats request",
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
		"remote_addr", c.Request().RemoteAddr,
	)
	var startDate, endDate time.Time
	var err error

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			slog.Error("Invalid start date format",
				"start_date", startDateStr,
				"error", err,
			)
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid start date format",
			})
		}
	} else {
		// Default to last 24 hours
		startDate = time.Now().Add(-24 * time.Hour)
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			slog.Error("Invalid end date format",
				"end_date", endDateStr,
				"error", err,
			)
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid end date format",
			})
		}
	} else {
		endDate = time.Now()
	}

	slog.Info("Getting analytics stats",
		"start_date", startDate,
		"end_date", endDate,
	)

	stats, err := h.analyticsService.GetAnalyticsStats(c.Request().Context(), startDate, endDate)
	if err != nil {
		slog.Error("Failed to get analytics stats",
			"error", err,
			"start_date", startDate,
			"end_date", endDate,
		)
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	slog.Info("Analytics stats retrieved successfully",
		"start_date", startDate,
		"end_date", endDate,
		"total_records", stats["total_records"],
	)

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    stats,
	})
}

// GetTopProtocols godoc
// @Summary Get top protocols
// @Description Get top protocols by usage
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit results" default(10)
// @Param start_date query string false "Start date (RFC3339)"
// @Param end_date query string false "End date (RFC3339)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/analytics/protocols [get]
func (h *AnalyticsHandler) GetTopProtocols(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var startDate, endDate time.Time
	var err error

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid start date format",
			})
		}
	} else {
		startDate = time.Now().Add(-24 * time.Hour)
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid end date format",
			})
		}
	} else {
		endDate = time.Now()
	}

	protocols, err := h.analyticsService.GetTopProtocols(c.Request().Context(), limit, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    protocols,
	})
}

// GetTopMethods godoc
// @Summary Get top methods
// @Description Get top methods by usage
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit results" default(10)
// @Param start_date query string false "Start date (RFC3339)"
// @Param end_date query string false "End date (RFC3339)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/analytics/methods [get]
func (h *AnalyticsHandler) GetTopMethods(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var startDate, endDate time.Time
	var err error

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid start date format",
			})
		}
	} else {
		startDate = time.Now().Add(-24 * time.Hour)
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid end date format",
			})
		}
	} else {
		endDate = time.Now()
	}

	methods, err := h.analyticsService.GetTopMethods(c.Request().Context(), limit, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    methods,
	})
}

// GetTrafficByHour godoc
// @Summary Get traffic by hour
// @Description Get traffic statistics by hour
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Param start_date query string false "Start date (RFC3339)"
// @Param end_date query string false "End date (RFC3339)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/analytics/traffic [get]
func (h *AnalyticsHandler) GetTrafficByHour(c echo.Context) error {
	var startDate, endDate time.Time
	var err error

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid start date format",
			})
		}
	} else {
		startDate = time.Now().Add(-24 * time.Hour)
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid end date format",
			})
		}
	} else {
		endDate = time.Now()
	}

	traffic, err := h.analyticsService.GetTrafficByHour(c.Request().Context(), startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    traffic,
	})
}

// GetErrorRate godoc
// @Summary Get error rate
// @Description Get error rate statistics
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Param start_date query string false "Start date (RFC3339)"
// @Param end_date query string false "End date (RFC3339)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/analytics/errors [get]
func (h *AnalyticsHandler) GetErrorRate(c echo.Context) error {
	var startDate, endDate time.Time
	var err error

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid start date format",
			})
		}
	} else {
		startDate = time.Now().Add(-24 * time.Hour)
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid end date format",
			})
		}
	} else {
		endDate = time.Now()
	}

	errorRate, err := h.analyticsService.GetErrorRate(c.Request().Context(), startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    errorRate,
	})
}

// GetPerformanceMetrics godoc
// @Summary Get performance metrics
// @Description Get performance metrics
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Param start_date query string false "Start date (RFC3339)"
// @Param end_date query string false "End date (RFC3339)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/analytics/performance [get]
func (h *AnalyticsHandler) GetPerformanceMetrics(c echo.Context) error {
	var startDate, endDate time.Time
	var err error

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid start date format",
			})
		}
	} else {
		startDate = time.Now().Add(-24 * time.Hour)
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Invalid end date format",
			})
		}
	} else {
		endDate = time.Now()
	}

	metrics, err := h.analyticsService.GetPerformanceMetrics(c.Request().Context(), startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    metrics,
	})
}
