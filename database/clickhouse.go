package database

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"hepic-app-server/v2/config"
	"hepic-app-server/v2/models"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseDB struct {
	conn clickhouse.Conn
}

func NewClickHouseConnection(cfg *config.Config) (*ClickHouseDB, error) {
	slog.Info("Connecting to ClickHouse",
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"database", cfg.Database.Database,
		"user", cfg.Database.User,
	)

	// Create ClickHouse connection options
	options := &clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Database.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database.Database,
			Username: cfg.Database.User,
			Password: cfg.Database.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout:      time.Duration(10) * time.Second,
		MaxOpenConns:     5,
		MaxIdleConns:     5,
		ConnMaxLifetime:  time.Duration(10) * time.Minute,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		BlockBufferSize:  10,
	}

	// Add compression if enabled
	if cfg.Database.Compress {
		options.Settings["enable_http_compression"] = 1
	}

	// Add SSL if enabled
	if cfg.Database.SSLMode == "require" {
		options.Settings["secure"] = 1
	}

	// Connect to ClickHouse
	conn, err := clickhouse.Open(options)
	if err != nil {
		slog.Error("Failed to connect to ClickHouse",
			"error", err,
			"host", cfg.Database.Host,
			"port", cfg.Database.Port,
		)
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		slog.Error("Failed to ping ClickHouse",
			"error", err,
			"host", cfg.Database.Host,
			"port", cfg.Database.Port,
		)
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	slog.Info("Successfully connected to ClickHouse",
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"database", cfg.Database.Database,
	)
	return &ClickHouseDB{conn: conn}, nil
}

func (ch *ClickHouseDB) Close() error {
	return ch.conn.Close()
}

// InitClickHouseTables creates necessary tables for HEP analytics
func (ch *ClickHouseDB) InitClickHouseTables() error {
	ctx := context.Background()

	// Create database if not exists
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", "hepic_analytics")
	if err := ch.conn.Exec(ctx, createDBQuery); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// Create HEP records table for analytics
	createHepTableQuery := `
	CREATE TABLE IF NOT EXISTS hep_analytics (
		id UInt64,
		call_id String,
		source_ip IPv4,
		destination_ip IPv4,
		protocol String,
		method String,
		status_code UInt16,
		timestamp DateTime64(3),
		raw_data String,
		created_at DateTime64(3) DEFAULT now64(3)
	) ENGINE = MergeTree()
	PARTITION BY toYYYYMM(timestamp)
	ORDER BY (timestamp, call_id)
	SETTINGS index_granularity = 8192
	`

	// Create users table for authentication
	createUsersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id UInt64,
		username String,
		email String,
		password String,
		role String,
		is_active UInt8,
		created_at DateTime,
		updated_at DateTime,
		last_login Nullable(DateTime)
	) ENGINE = MergeTree()
	ORDER BY (id)
	SETTINGS index_granularity = 8192
	`

	if err := ch.conn.Exec(ctx, createHepTableQuery); err != nil {
		return fmt.Errorf("failed to create hep_analytics table: %w", err)
	}

	if err := ch.conn.Exec(ctx, createUsersTableQuery); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create materialized view for real-time statistics
	mvQuery := `
	CREATE MATERIALIZED VIEW IF NOT EXISTS hep_stats_mv
	ENGINE = SummingMergeTree()
	PARTITION BY toYYYYMM(timestamp)
	ORDER BY (timestamp, protocol, method, status_code)
	AS SELECT
		toStartOfMinute(timestamp) as timestamp,
		protocol,
		method,
		status_code,
		count() as count
	FROM hep_analytics
	GROUP BY timestamp, protocol, method, status_code
	`

	if err := ch.conn.Exec(ctx, mvQuery); err != nil {
		return fmt.Errorf("failed to create materialized view: %w", err)
	}

	// Create distributed table for scaling (optional)
	distributedQuery := `
	CREATE TABLE IF NOT EXISTS hep_analytics_distributed AS hep_analytics
	ENGINE = Distributed('cluster', 'hepic_analytics', 'hep_analytics', rand())
	`

	if err := ch.conn.Exec(ctx, distributedQuery); err != nil {
		log.Printf("Warning: Failed to create distributed table (cluster not configured): %v", err)
	}

	log.Println("ClickHouse tables initialized successfully")
	return nil
}

// InsertHEPRecord inserts a HEP record into ClickHouse
func (ch *ClickHouseDB) InsertHEPRecord(ctx context.Context, record HEPRecord) error {
	query := `
	INSERT INTO hep_analytics (
		id, call_id, source_ip, destination_ip, protocol, 
		method, status_code, timestamp, raw_data, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	return ch.conn.Exec(ctx, query,
		record.ID,
		record.CallID,
		record.SourceIP,
		record.DestinationIP,
		record.Protocol,
		record.Method,
		record.StatusCode,
		record.Timestamp,
		record.RawData,
		record.CreatedAt,
	)
}

// GetHEPStats returns analytics statistics from ClickHouse
func (ch *ClickHouseDB) GetHEPStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total records count
	var totalRecords uint64
	countQuery := `
	SELECT count() as total 
	FROM hep_analytics 
	WHERE timestamp >= ? AND timestamp <= ?
	`

	row := ch.conn.QueryRow(ctx, countQuery, startDate, endDate)
	if err := row.Scan(&totalRecords); err != nil {
		return nil, fmt.Errorf("failed to get total records: %w", err)
	}

	// Protocol statistics
	protocolQuery := `
	SELECT protocol, count() as count
	FROM hep_analytics 
	WHERE timestamp >= ? AND timestamp <= ?
	GROUP BY protocol 
	ORDER BY count DESC
	LIMIT 10
	`

	protocolRows, err := ch.conn.Query(ctx, protocolQuery, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get protocol stats: %w", err)
	}
	defer protocolRows.Close()

	var protocolStats []map[string]interface{}
	for protocolRows.Next() {
		var protocol string
		var count uint64
		if err := protocolRows.Scan(&protocol, &count); err != nil {
			continue
		}
		protocolStats = append(protocolStats, map[string]interface{}{
			"protocol": protocol,
			"count":    count,
		})
	}

	// Method statistics
	methodQuery := `
	SELECT method, count() as count
	FROM hep_analytics 
	WHERE timestamp >= ? AND timestamp <= ? AND method != ''
	GROUP BY method 
	ORDER BY count DESC
	LIMIT 10
	`

	methodRows, err := ch.conn.Query(ctx, methodQuery, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get method stats: %w", err)
	}
	defer methodRows.Close()

	var methodStats []map[string]interface{}
	for methodRows.Next() {
		var method string
		var count uint64
		if err := methodRows.Scan(&method, &count); err != nil {
			continue
		}
		methodStats = append(methodStats, map[string]interface{}{
			"method": method,
			"count":  count,
		})
	}

	stats["total_records"] = totalRecords
	stats["protocol_stats"] = protocolStats
	stats["method_stats"] = methodStats

	return stats, nil
}

// User management methods

// InsertUser inserts a new user into the database
func (ch *ClickHouseDB) InsertUser(ctx context.Context, user *models.User) (int64, error) {
	query := `
	INSERT INTO users (username, email, password, role, is_active, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	// Generate user ID (simple auto-increment simulation)
	userID := time.Now().UnixNano()

	err := ch.conn.Exec(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUserByID retrieves a user by ID
func (ch *ClickHouseDB) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
	SELECT id, username, email, password, role, is_active, created_at, updated_at, last_login
	FROM users
	WHERE id = ?
	LIMIT 1`

	row := ch.conn.QueryRow(ctx, query, userID)

	user := &models.User{}
	var lastLogin *time.Time

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLogin,
	)

	if err != nil {
		return nil, err
	}

	user.LastLogin = lastLogin
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (ch *ClickHouseDB) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
	SELECT id, username, email, password, role, is_active, created_at, updated_at, last_login
	FROM users
	WHERE username = ?
	LIMIT 1`

	row := ch.conn.QueryRow(ctx, query, username)

	user := &models.User{}
	var lastLogin *time.Time

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLogin,
	)

	if err != nil {
		return nil, err
	}

	user.LastLogin = lastLogin
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (ch *ClickHouseDB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
	SELECT id, username, email, password, role, is_active, created_at, updated_at, last_login
	FROM users
	WHERE email = ?
	LIMIT 1`

	row := ch.conn.QueryRow(ctx, query, email)

	user := &models.User{}
	var lastLogin *time.Time

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLogin,
	)

	if err != nil {
		return nil, err
	}

	user.LastLogin = lastLogin
	return user, nil
}

// UpdateUser updates a user
func (ch *ClickHouseDB) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
	ALTER TABLE users UPDATE
	username = ?, email = ?, role = ?, is_active = ?, updated_at = ?
	WHERE id = ?`

	err := ch.conn.Exec(ctx, query,
		user.Username,
		user.Email,
		user.Role,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// UpdateUserPassword updates a user's password
func (ch *ClickHouseDB) UpdateUserPassword(ctx context.Context, userID int64, hashedPassword string) error {
	query := `
	ALTER TABLE users UPDATE
	password = ?, updated_at = ?
	WHERE id = ?`

	err := ch.conn.Exec(ctx, query,
		hashedPassword,
		time.Now(),
		userID,
	)

	return err
}

// UpdateUserLastLogin updates a user's last login time
func (ch *ClickHouseDB) UpdateUserLastLogin(ctx context.Context, userID int64, lastLogin time.Time) error {
	query := `
	ALTER TABLE users UPDATE
	last_login = ?, updated_at = ?
	WHERE id = ?`

	err := ch.conn.Exec(ctx, query,
		lastLogin,
		time.Now(),
		userID,
	)

	return err
}

// GetUsers retrieves a paginated list of users
func (ch *ClickHouseDB) GetUsers(ctx context.Context, page, perPage int, role string) (*models.UserListResponse, error) {
	offset := (page - 1) * perPage

	// Build query with optional role filter
	whereClause := ""
	args := []interface{}{}
	if role != "" {
		whereClause = "WHERE role = ?"
		args = append(args, role)
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var total int64
	if role != "" {
		err := ch.conn.QueryRow(ctx, countQuery, role).Scan(&total)
		if err != nil {
			return nil, err
		}
	} else {
		err := ch.conn.QueryRow(ctx, countQuery).Scan(&total)
		if err != nil {
			return nil, err
		}
	}

	// Get users
	query := fmt.Sprintf(`
	SELECT id, username, email, role, is_active, created_at, updated_at, last_login
	FROM users %s
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?`, whereClause)

	args = append(args, perPage, offset)

	rows, err := ch.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user := models.User{}
		var lastLogin *time.Time

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&lastLogin,
		)
		if err != nil {
			return nil, err
		}

		user.LastLogin = lastLogin
		users = append(users, user)
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return &models.UserListResponse{
		Users:      users,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetUserStats retrieves user statistics
func (ch *ClickHouseDB) GetUserStats(ctx context.Context) (*models.UserStats, error) {
	query := `
	SELECT 
		COUNT(*) as total_users,
		COUNTIf(is_active = 1) as active_users,
		COUNTIf(role = 'admin') as admin_users,
		COUNTIf(role = 'user') as regular_users,
		COUNTIf(created_at >= today()) as new_users_today
	FROM users`

	row := ch.conn.QueryRow(ctx, query)

	stats := &models.UserStats{}
	err := row.Scan(
		&stats.TotalUsers,
		&stats.ActiveUsers,
		&stats.AdminUsers,
		&stats.RegularUsers,
		&stats.NewUsersToday,
	)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// DeleteUser deletes a user
func (ch *ClickHouseDB) DeleteUser(ctx context.Context, userID int64) error {
	query := "ALTER TABLE users DELETE WHERE id = ?"
	err := ch.conn.Exec(ctx, query, userID)
	return err
}

// HEPRecord represents a HEP record for ClickHouse
type HEPRecord struct {
	ID            uint64    `json:"id"`
	CallID        string    `json:"call_id"`
	SourceIP      string    `json:"source_ip"`
	DestinationIP string    `json:"destination_ip"`
	Protocol      string    `json:"protocol"`
	Method        string    `json:"method"`
	StatusCode    uint16    `json:"status_code"`
	Timestamp     time.Time `json:"timestamp"`
	RawData       string    `json:"raw_data"`
	CreatedAt     time.Time `json:"created_at"`
}
