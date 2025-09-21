package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"hepic-app-server/v2/config"
	"hepic-app-server/v2/database"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
	Long: `Configuration management commands for HEPIC App Server.

This command provides utilities for managing configuration files,
validating settings, and generating example configurations.`,
}

// configValidateCmd represents the config validate command
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long: `Validate the configuration file for syntax errors and required fields.

This command checks:
- JSON/YAML syntax validity
- Required fields presence
- Data type validation
- Value range validation
- ClickHouse connectivity (optional)

Examples:
  hepic-app-server config validate
  hepic-app-server config validate --config /path/to/config.json
  hepic-app-server config validate --check-db`,
	Run: runConfigValidate,
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long: `Display the current configuration with sensitive data masked.

This command shows the loaded configuration with:
- Masked passwords and secrets
- Source of each setting (file, env, default)
- Validation status
- ClickHouse connection info (without credentials)

Examples:
  hepic-app-server config show
  hepic-app-server config show --config /path/to/config.json`,
	Run: runConfigShow,
}

// configGenerateCmd represents the config generate command
var configGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate example configuration",
	Long: `Generate example configuration files in different formats.

This command creates example configuration files:
- JSON format (config.json)
- YAML format (config.yaml)
- Environment variables (.env)
- Docker Compose (docker-compose.yml)

Examples:
  hepic-app-server config generate
  hepic-app-server config generate --format yaml
  hepic-app-server config generate --output /path/to/config/`,
	Run: runConfigGenerate,
}

var (
	checkDB     bool
	showSecrets bool
	format      string
	output      string
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configGenerateCmd)

	// Validate command flags
	configValidateCmd.Flags().BoolVar(&checkDB, "check-db", false, "Check ClickHouse connectivity")

	// Show command flags
	configShowCmd.Flags().BoolVar(&showSecrets, "show-secrets", false, "Show sensitive data (passwords, secrets)")

	// Generate command flags
	configGenerateCmd.Flags().StringVar(&format, "format", "json", "Output format (json, yaml, env, docker)")
	configGenerateCmd.Flags().StringVar(&output, "output", ".", "Output directory")
}

func runConfigValidate(cmd *cobra.Command, args []string) {
	fmt.Println("Validating configuration...")

	// Load configuration
	cfg := config.Load()

	// Basic validation
	if err := config.ValidateConfig(cfg); err != nil {
		fmt.Printf("❌ Configuration validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Configuration is valid")

	// Check ClickHouse connectivity if requested
	if checkDB {
		fmt.Println("Checking ClickHouse connectivity...")
		
		clickhouse, err := database.NewClickHouseConnection(cfg)
		if err != nil {
			fmt.Printf("❌ ClickHouse connection failed: %v\n", err)
			os.Exit(1)
		}
		defer clickhouse.Close()

		fmt.Println("✅ ClickHouse connection successful")
	}

	fmt.Println("🎉 All validations passed!")
}

func runConfigShow(cmd *cobra.Command, args []string) {
	fmt.Println("Current configuration:")
	fmt.Println("====================")

	// Load configuration
	cfg := config.Load()

	// Create a copy for display
	displayCfg := *cfg

	// Mask sensitive data unless requested
	if !showSecrets {
		if displayCfg.Database.Password != "" {
			displayCfg.Database.Password = "***"
		}
		if displayCfg.JWT.Secret != "" {
			displayCfg.JWT.Secret = "***"
		}
	}

	// Convert to JSON for display
	jsonData, err := json.MarshalIndent(displayCfg, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling config: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))

	// Show configuration source
	fmt.Println("\nConfiguration source:")
	fmt.Printf("- Config file: %s\n", viper.ConfigFileUsed())
	fmt.Printf("- Environment variables: %s\n", viper.GetString("HEPIC_DATABASE_HOST"))
}

func runConfigGenerate(cmd *cobra.Command, args []string) {
	fmt.Printf("Generating example configuration in %s format...\n", format)

	switch format {
	case "json":
		generateJSONConfig()
	case "yaml":
		generateYAMLConfig()
	case "env":
		generateEnvConfig()
	case "docker":
		generateDockerConfig()
	default:
		fmt.Printf("❌ Unsupported format: %s\n", format)
		os.Exit(1)
	}

	fmt.Println("✅ Example configuration generated successfully!")
}

func generateJSONConfig() {
	config := `{
  "database": {
    "host": "localhost",
    "port": 9000,
    "user": "default",
    "password": "",
    "database": "hepic_analytics",
    "sslmode": "disable",
    "compress": true
  },
  "server": {
    "port": "8080",
    "host": "0.0.0.0"
  },
  "jwt": {
    "secret": "your-super-secret-jwt-key-here-change-in-production",
    "expire_hours": 24
  },
  "logging": {
    "level": "info"
  }
}`

	filename := output + "/config.json"
	if err := os.WriteFile(filename, []byte(config), 0644); err != nil {
		fmt.Printf("❌ Failed to write config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📄 Generated: %s\n", filename)
}

func generateYAMLConfig() {
	config := `database:
  host: localhost
  port: 9000
  user: default
  password: ""
  database: hepic_analytics
  sslmode: disable
  compress: true

server:
  port: "8080"
  host: "0.0.0.0"

jwt:
  secret: "your-super-secret-jwt-key-here-change-in-production"
  expire_hours: 24

logging:
  level: info`

	filename := output + "/config.yaml"
	if err := os.WriteFile(filename, []byte(config), 0644); err != nil {
		fmt.Printf("❌ Failed to write config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📄 Generated: %s\n", filename)
}

func generateEnvConfig() {
	config := `# HEPIC App Server v2 Configuration
# Database Configuration (ClickHouse)
HEPIC_DATABASE_HOST=localhost
HEPIC_DATABASE_PORT=9000
HEPIC_DATABASE_USER=default
HEPIC_DATABASE_PASSWORD=
HEPIC_DATABASE_DATABASE=hepic_analytics
HEPIC_DATABASE_SSLMODE=disable
HEPIC_DATABASE_COMPRESS=true

# Server Configuration
HEPIC_SERVER_PORT=8080
HEPIC_SERVER_HOST=0.0.0.0

# JWT Configuration
HEPIC_JWT_SECRET=your-super-secret-jwt-key-here-change-in-production
HEPIC_JWT_EXPIRE_HOURS=24

# Logging
HEPIC_LOGGING_LEVEL=info`

	filename := output + "/.env"
	if err := os.WriteFile(filename, []byte(config), 0644); err != nil {
		fmt.Printf("❌ Failed to write env file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📄 Generated: %s\n", filename)
}

func generateDockerConfig() {
	config := `version: '3.8'

services:
  clickhouse:
    image: clickhouse/clickhouse-server:latest
    container_name: hepic-clickhouse
    environment:
      CLICKHOUSE_DB: hepic_analytics
      CLICKHOUSE_USER: default
      CLICKHOUSE_PASSWORD: ""
    ports:
      - "9000:9000"
      - "8123:8123"
    volumes:
      - clickhouse_data:/var/lib/clickhouse
    networks:
      - hepic-network

  hepic-app-server:
    build: .
    container_name: hepic-app-server-v2
    ports:
      - "8080:8080"
    environment:
      # ClickHouse
      - HEPIC_DATABASE_HOST=clickhouse
      - HEPIC_DATABASE_PORT=9000
      - HEPIC_DATABASE_USER=default
      - HEPIC_DATABASE_PASSWORD=
      - HEPIC_DATABASE_DATABASE=hepic_analytics
      - HEPIC_DATABASE_SSLMODE=disable
      - HEPIC_DATABASE_COMPRESS=true
      
      # Server
      - HEPIC_SERVER_PORT=8080
      - HEPIC_SERVER_HOST=0.0.0.0
      
      # JWT
      - HEPIC_JWT_SECRET=your-super-secret-jwt-key-here-change-in-production
      - HEPIC_JWT_EXPIRE_HOURS=24
      
      # Logging
      - HEPIC_LOGGING_LEVEL=info
    depends_on:
      - clickhouse
    networks:
      - hepic-network

volumes:
  clickhouse_data:

networks:
  hepic-network:
    driver: bridge`

	filename := output + "/docker-compose.yml"
	if err := os.WriteFile(filename, []byte(config), 0644); err != nil {
		fmt.Printf("❌ Failed to write docker file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📄 Generated: %s\n", filename)
}
