package services

import (
	"context"
	"log/slog"
	"time"

	"hepic-app-server/v2/database"
	"hepic-app-server/v2/models"
)

type AnalyticsService struct {
	clickhouse *database.ClickHouseDB
}

func NewAnalyticsService(clickhouse *database.ClickHouseDB) *AnalyticsService {
	return &AnalyticsService{
		clickhouse: clickhouse,
	}
}

// InsertHEPRecord inserts a HEP record into ClickHouse for analytics
func (s *AnalyticsService) InsertHEPRecord(ctx context.Context, record models.HEPRecord) error {
	// Convert models.HEPRecord to database.HEPRecord
	chRecord := database.HEPRecord{
		ID:            uint64(record.ID),
		CallID:        record.CallID,
		SourceIP:      record.SourceIP,
		DestinationIP: record.DestinationIP,
		Protocol:      record.Protocol,
		Method:        record.Method,
		StatusCode:    uint16(record.StatusCode),
		Timestamp:     record.Timestamp,
		RawData:       record.RawData,
		CreatedAt:     record.CreatedAt,
	}

	return s.clickhouse.InsertHEPRecord(ctx, chRecord)
}

// GetAnalyticsStats returns comprehensive analytics from ClickHouse
func (s *AnalyticsService) GetAnalyticsStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	slog.Info("Getting analytics stats from ClickHouse",
		"start_date", startDate,
		"end_date", endDate,
	)

	stats, err := s.clickhouse.GetHEPStats(ctx, startDate, endDate)
	if err != nil {
		slog.Error("Failed to get analytics stats from ClickHouse",
			"error", err,
			"start_date", startDate,
			"end_date", endDate,
		)
		return nil, err
	}

	slog.Info("Analytics stats retrieved from ClickHouse",
		"start_date", startDate,
		"end_date", endDate,
		"total_records", stats["total_records"],
	)

	return stats, nil
}

// GetRealTimeStats returns real-time statistics using materialized views
func (s *AnalyticsService) GetRealTimeStats(ctx context.Context, minutes int) (map[string]interface{}, error) {
	// This would query the materialized view for real-time stats
	// Implementation depends on specific ClickHouse setup
	return map[string]interface{}{
		"time_range_minutes": minutes,
		"message":            "Real-time stats feature coming soon",
	}, nil
}

// GetTopProtocols returns top protocols by usage
func (s *AnalyticsService) GetTopProtocols(ctx context.Context, limit int, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	stats, err := s.clickhouse.GetHEPStats(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	protocolStats, ok := stats["protocol_stats"].([]map[string]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}

	// Limit results
	if len(protocolStats) > limit {
		protocolStats = protocolStats[:limit]
	}

	return protocolStats, nil
}

// GetTopMethods returns top methods by usage
func (s *AnalyticsService) GetTopMethods(ctx context.Context, limit int, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	stats, err := s.clickhouse.GetHEPStats(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	methodStats, ok := stats["method_stats"].([]map[string]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}

	// Limit results
	if len(methodStats) > limit {
		methodStats = methodStats[:limit]
	}

	return methodStats, nil
}

// GetTrafficByHour returns traffic statistics by hour
func (s *AnalyticsService) GetTrafficByHour(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	// This would implement hourly traffic analysis
	// For now, return a placeholder
	return []map[string]interface{}{
		{
			"hour":  "00:00",
			"count": 0,
		},
	}, nil
}

// GetGeographicStats returns geographic distribution of traffic
func (s *AnalyticsService) GetGeographicStats(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	// This would implement geographic analysis based on IP addresses
	// For now, return a placeholder
	return []map[string]interface{}{
		{
			"country": "Unknown",
			"count":   0,
		},
	}, nil
}

// GetErrorRate returns error rate statistics
func (s *AnalyticsService) GetErrorRate(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	// This would calculate error rates based on status codes
	// For now, return a placeholder
	return map[string]interface{}{
		"total_requests": 0,
		"error_requests": 0,
		"error_rate":     0.0,
	}, nil
}

// GetPerformanceMetrics returns performance metrics
func (s *AnalyticsService) GetPerformanceMetrics(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	// This would calculate performance metrics
	// For now, return a placeholder
	return map[string]interface{}{
		"avg_response_time": 0.0,
		"max_response_time": 0.0,
		"min_response_time": 0.0,
		"throughput":        0.0,
	}, nil
}
