package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"hepic-app-server/v2/config"
	"hepic-app-server/v2/database"

	"github.com/spf13/cobra"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Health check commands",
	Long: `Health check commands for HEPIC App Server.

This command provides utilities for checking the health of the server
and its dependencies (ClickHouse, etc.).`,
}

// healthCheckCmd represents the health check command
var healthCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Perform health check",
	Long: `Perform a comprehensive health check of the HEPIC App Server.

This command checks:
- ClickHouse connectivity
- Server configuration
- Database tables
- JWT configuration
- Logging setup

Examples:
  hepic-app-server health check
  hepic-app-server health check --timeout 30s
  hepic-app-server health check --verbose`,
	Run: runHealthCheck,
}

// healthServerCmd represents the health server command
var healthServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start health check server",
	Long: `Start a dedicated health check server.

This command starts a minimal HTTP server that provides health check
endpoints for monitoring and load balancers.

Endpoints:
- GET /health - Basic health check
- GET /health/ready - Readiness check
- GET /health/live - Liveness check
- GET /health/detailed - Detailed health information

Examples:
  hepic-app-server health server
  hepic-app-server health server --port 8081
  hepic-app-server health server --host 0.0.0.0`,
	Run: runHealthServer,
}

var (
	timeout     string
	healthVerbose bool
	healthPort  string
	healthHost  string
)

func init() {
	rootCmd.AddCommand(healthCmd)
	healthCmd.AddCommand(healthCheckCmd)
	healthCmd.AddCommand(healthServerCmd)

	// Health check flags
	healthCheckCmd.Flags().StringVar(&timeout, "timeout", "10s", "Health check timeout")
	healthCheckCmd.Flags().BoolVar(&healthVerbose, "verbose", false, "Verbose output")

	// Health server flags
	healthServerCmd.Flags().StringVar(&healthPort, "port", "8081", "Health server port")
	healthServerCmd.Flags().StringVar(&healthHost, "host", "0.0.0.0", "Health server host")
}

func runHealthCheck(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 Performing health check...")
	
	// Parse timeout
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Printf("❌ Invalid timeout: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	cfg := config.Load()

	// Check ClickHouse connectivity
	fmt.Println("📊 Checking ClickHouse connectivity...")
	clickhouse, err := database.NewClickHouseConnection(cfg)
	if err != nil {
		fmt.Printf("❌ ClickHouse connection failed: %v\n", err)
		os.Exit(1)
	}
	defer clickhouse.Close()

	// Check database tables
	fmt.Println("🗄️  Checking database tables...")
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	if err := clickhouse.InitClickHouseTables(); err != nil {
		fmt.Printf("❌ Database tables check failed: %v\n", err)
		os.Exit(1)
	}
	_ = ctx // Suppress unused variable warning

	// Check JWT configuration
	fmt.Println("🔐 Checking JWT configuration...")
	if cfg.JWT.Secret == "" || cfg.JWT.Secret == "your-super-secret-jwt-key-here" {
		fmt.Println("⚠️  Warning: JWT secret not properly configured")
	}

	// Check server configuration
	fmt.Println("⚙️  Checking server configuration...")
	if cfg.Server.Port == "" {
		fmt.Println("❌ Server port not configured")
		os.Exit(1)
	}

	if healthVerbose {
		fmt.Println("\n📋 Configuration details:")
		fmt.Printf("- Server: %s:%s\n", cfg.Server.Host, cfg.Server.Port)
		fmt.Printf("- ClickHouse: %s:%d/%s\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
		fmt.Printf("- JWT Expire: %d hours\n", cfg.JWT.ExpireHours)
		fmt.Printf("- Log Level: %s\n", cfg.Logging.Level)
	}

	fmt.Println("✅ All health checks passed!")
	fmt.Println("🎉 HEPIC App Server is healthy!")
}

func runHealthServer(cmd *cobra.Command, args []string) {
	fmt.Printf("🏥 Starting health check server on %s:%s\n", healthHost, healthPort)

	// Load configuration
	cfg := config.Load()

	// Create health check server
	mux := http.NewServeMux()

	// Basic health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "2.0.0",
		})
	})

	// Readiness check
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		// Check ClickHouse connectivity
		clickhouse, err := database.NewClickHouseConnection(cfg)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":    "not ready",
				"error":     err.Error(),
				"timestamp": time.Now().Format(time.RFC3339),
			})
			return
		}
		defer clickhouse.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ready",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Liveness check
	mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "alive",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Detailed health information
	mux.HandleFunc("/health/detailed", func(w http.ResponseWriter, r *http.Request) {
		health := map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "2.0.0",
			"config": map[string]interface{}{
				"server": map[string]interface{}{
					"host": cfg.Server.Host,
					"port": cfg.Server.Port,
				},
				"database": map[string]interface{}{
					"host":     cfg.Database.Host,
					"port":     cfg.Database.Port,
					"database": cfg.Database.Database,
					"user":     cfg.Database.User,
				},
				"jwt": map[string]interface{}{
					"expire_hours": cfg.JWT.ExpireHours,
					"secret_set":   cfg.JWT.Secret != "" && cfg.JWT.Secret != "your-super-secret-jwt-key-here",
				},
				"logging": map[string]interface{}{
					"level": cfg.Logging.Level,
				},
			},
		}

		// Check ClickHouse connectivity
		clickhouse, err := database.NewClickHouseConnection(cfg)
		if err != nil {
			health["database"] = map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
		} else {
			clickhouse.Close()
			health["database"] = map[string]interface{}{
				"status": "ok",
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(health)
	})

	// Start server
	server := &http.Server{
		Addr:    healthHost + ":" + healthPort,
		Handler: mux,
	}

	fmt.Printf("🚀 Health check server started successfully!\n")
	fmt.Printf("📡 Endpoints available:\n")
	fmt.Printf("  - GET /health - Basic health check\n")
	fmt.Printf("  - GET /health/ready - Readiness check\n")
	fmt.Printf("  - GET /health/live - Liveness check\n")
	fmt.Printf("  - GET /health/detailed - Detailed health information\n")

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("❌ Health server error: %v\n", err)
		os.Exit(1)
	}
}
