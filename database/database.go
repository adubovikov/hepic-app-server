package database

import (
	"fmt"
	"log"

	"hepic-app-server/v2/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func NewConnection(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Database,
		cfg.Database.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Connection check
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}

	log.Println("Database connection successful")

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// InitTables creates necessary tables
func (db *DB) InitTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(20) DEFAULT 'user',
			active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			token VARCHAR(255) UNIQUE NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS hep_records (
			id SERIAL PRIMARY KEY,
			call_id VARCHAR(100),
			source_ip VARCHAR(45),
			destination_ip VARCHAR(45),
			protocol VARCHAR(20),
			method VARCHAR(20),
			status_code INTEGER,
			timestamp TIMESTAMP,
			raw_data TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_hep_call_id ON hep_records(call_id)`,
		`CREATE INDEX IF NOT EXISTS idx_hep_timestamp ON hep_records(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_hep_source_ip ON hep_records(source_ip)`,
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("query execution error: %w", err)
		}
	}

	log.Println("Database tables initialized")
	return nil
}
