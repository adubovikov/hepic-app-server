package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HEPIC App Server",
	Long: `Start the HEPIC App Server with ClickHouse integration.

This command starts the REST API server with the following features:
- ClickHouse analytics database
- Structured JSON logging
- JWT authentication
- Health monitoring
- Graceful shutdown

Examples:
  hepic-app-server serve
  hepic-app-server serve --config /path/to/config.json
  hepic-app-server serve --port 8080 --host 0.0.0.0
  hepic-app-server serve --log-level debug`,
	Run: runServe,
}

// Variables are now in root.go

func init() {
	rootCmd.AddCommand(serveCmd)

	// Server flags (using variables from root.go)
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "port to listen on")
	serveCmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "host to bind to")
	serveCmd.Flags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error)")
	serveCmd.Flags().StringVar(&logFormat, "log-format", "json", "log format (json, text)")

	// Database flags
	serveCmd.Flags().String("db-host", "localhost", "ClickHouse host")
	serveCmd.Flags().Int("db-port", 9000, "ClickHouse port")
	serveCmd.Flags().String("db-user", "default", "ClickHouse user")
	serveCmd.Flags().String("db-password", "", "ClickHouse password")
	serveCmd.Flags().String("db-database", "hepic_analytics", "ClickHouse database")
	serveCmd.Flags().Bool("db-compress", true, "Enable ClickHouse compression")

	// JWT flags
	serveCmd.Flags().String("jwt-secret", "", "JWT secret key")
	serveCmd.Flags().Int("jwt-expire-hours", 24, "JWT token expiration in hours")

	// Bind flags to viper
	viper.BindPFlag("server.port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.host", serveCmd.Flags().Lookup("host"))
	viper.BindPFlag("database.host", serveCmd.Flags().Lookup("db-host"))
	viper.BindPFlag("database.port", serveCmd.Flags().Lookup("db-port"))
	viper.BindPFlag("database.user", serveCmd.Flags().Lookup("db-user"))
	viper.BindPFlag("database.password", serveCmd.Flags().Lookup("db-password"))
	viper.BindPFlag("database.database", serveCmd.Flags().Lookup("db-database"))
	viper.BindPFlag("database.compress", serveCmd.Flags().Lookup("db-compress"))
	viper.BindPFlag("jwt.secret", serveCmd.Flags().Lookup("jwt-secret"))
	viper.BindPFlag("jwt.expire_hours", serveCmd.Flags().Lookup("jwt-expire-hours"))
}

// runServe function is now in root.go
