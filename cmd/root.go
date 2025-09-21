package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hepic-app-server/v2/config"
	"hepic-app-server/v2/database"
	appMiddleware "hepic-app-server/v2/middleware"
	"hepic-app-server/v2/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	// Server flags
	port      string
	host      string
	logLevel  string
	logFormat string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hepic-app-server",
	Short: "HEPIC App Server v2 - Advanced REST API Server for Analytics",
	Long: `HEPIC App Server v2 is a high-performance REST API server designed for 
analytics and monitoring of HEP (Homer Encapsulation Protocol) data.

Features:
- ClickHouse integration for analytics
- Structured JSON logging with slog
- JWT authentication
- Real-time statistics
- Docker support
- Health monitoring

Built with Go and Echo framework.`,
	Version: "2.0.0",
	Run:     runServe, // Set serve as default command
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("log-format", "json", "log format (json, text)")

	// Server flags (for default serve command)
	rootCmd.Flags().StringVarP(&port, "port", "p", "8080", "port to listen on")
	rootCmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "host to bind to")
	rootCmd.Flags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.Flags().StringVar(&logFormat, "log-format", "json", "log format (json, text)")

	// Database flags
	rootCmd.Flags().String("db-host", "localhost", "ClickHouse host")
	rootCmd.Flags().Int("db-port", 9000, "ClickHouse port")
	rootCmd.Flags().String("db-user", "default", "ClickHouse user")
	rootCmd.Flags().String("db-password", "", "ClickHouse password")
	rootCmd.Flags().String("db-database", "hepic_analytics", "ClickHouse database")
	rootCmd.Flags().Bool("db-compress", true, "Enable ClickHouse compression")

	// JWT flags
	rootCmd.Flags().String("jwt-secret", "", "JWT secret key")
	rootCmd.Flags().Int("jwt-expire-hours", 24, "JWT token expiration in hours")

	// Bind flags to viper
	viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("logging.format", rootCmd.PersistentFlags().Lookup("log-format"))
	viper.BindPFlag("server.port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.host", rootCmd.Flags().Lookup("host"))
	viper.BindPFlag("database.host", rootCmd.Flags().Lookup("db-host"))
	viper.BindPFlag("database.port", rootCmd.Flags().Lookup("db-port"))
	viper.BindPFlag("database.user", rootCmd.Flags().Lookup("db-user"))
	viper.BindPFlag("database.password", rootCmd.Flags().Lookup("db-password"))
	viper.BindPFlag("database.database", rootCmd.Flags().Lookup("db-database"))
	viper.BindPFlag("database.compress", rootCmd.Flags().Lookup("db-compress"))
	viper.BindPFlag("jwt.secret", rootCmd.Flags().Lookup("jwt-secret"))
	viper.BindPFlag("jwt.expire_hours", rootCmd.Flags().Lookup("jwt-expire-hours"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config files in order of preference
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("/etc/hepic-app-server")
		viper.AddConfigPath("$HOME/.hepic-app-server")
	}

	// Environment variables
	viper.SetEnvPrefix("HEPIC")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}

func runServe(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg := config.Load()

	// Override with command line flags if provided
	if port != "8080" {
		cfg.Server.Port = port
	}
	if host != "0.0.0.0" {
		cfg.Server.Host = host
	}

	// Setup logger
	setupLogger(logLevel, logFormat)

	// Create Echo instance
	e := echo.New()

	// Setup validator
	appMiddleware.SetupValidator(e)

	// Setup middleware
	setupMiddleware(e)

	// Connect to ClickHouse
	clickhouse, err := database.NewClickHouseConnection(cfg)
	if err != nil {
		slog.Error("Failed to connect to ClickHouse", "error", err)
		os.Exit(1)
	}
	defer clickhouse.Close()

	// Initialize ClickHouse tables
	if err := clickhouse.InitClickHouseTables(); err != nil {
		slog.Error("Failed to initialize ClickHouse tables", "error", err)
		os.Exit(1)
	}

	// Setup routes
	routes.SetupRoutes(e, clickhouse, cfg.JWT.Secret)

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	slog.Info("Starting HEPIC App Server v2",
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
		"version", "2.0.0",
	)

	// Graceful shutdown
	go func() {
		if err := e.Start(serverAddr); err != nil {
			slog.Error("Server startup error", "error", err)
		}
	}()

	// Wait for signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}
}

func setupLogger(level, format string) {
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	var handler slog.Handler
	if format == "text" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slogLevel,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slogLevel,
		})
	}

	slog.SetDefault(slog.New(handler))
}

func setupMiddleware(e *echo.Echo) {
	// CORS
	e.Use(middleware.CORS())

	// Slog logging middleware
	e.Use(appMiddleware.Slog())

	// Slog error logging
	e.Use(appMiddleware.SlogError())

	// Slog panic recovery
	e.Use(appMiddleware.SlogRecover())

	// Response compression
	e.Use(middleware.Gzip())

	// Request body size limit
	e.Use(middleware.BodyLimit("10M"))

	// Timeouts
	e.Use(middleware.Timeout())

	// Security headers
	e.Use(middleware.Secure())

	// Remove trailing slash
	e.Pre(middleware.RemoveTrailingSlash())

	// Add request ID
	e.Use(middleware.RequestID())
}
