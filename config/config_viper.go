package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database ClickHouseConfig `mapstructure:"database"`
	Server   ServerConfig     `mapstructure:"server"`
	JWT      JWTConfig        `mapstructure:"jwt"`
	Logging  LoggingConfig    `mapstructure:"logging"`
}

type ClickHouseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
	Compress bool   `mapstructure:"compress"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

func Load() *Config {
	// Configure Viper
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/hepic-app-server/")
	viper.AddConfigPath("$HOME/.hepic-app-server")

	// Automatic reading of environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	setDefaults()

	// Read configuration from file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults and environment variables")
		} else {
			log.Printf("Warning: Error reading config file: %v", err)
		}
	} else {
		log.Printf("Config file loaded: %s", viper.ConfigFileUsed())
	}

	// Read .env file if exists
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	if err := viper.MergeInConfig(); err == nil {
		log.Println("Environment file (.env) loaded")
	}

	// Create configuration structure
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}

	log.Println("Configuration loaded successfully")
	logConfig(&config)
	return &config
}

func setDefaults() {
	// Database defaults (ClickHouse)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 9000)
	viper.SetDefault("database.user", "default")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "hepic_analytics")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.compress", true)

	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-super-secret-jwt-key-here")
	viper.SetDefault("jwt.expire_hours", 24)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
}

func validateConfig(config *Config) error {
	// Validate required fields
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}
	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if config.JWT.Secret == "" || config.JWT.Secret == "your-super-secret-jwt-key-here" {
		return fmt.Errorf("JWT secret must be set to a secure value")
	}
	if config.JWT.ExpireHours <= 0 {
		return fmt.Errorf("JWT expire hours must be greater than 0")
	}

	return nil
}

// ValidateConfig validates the configuration
func ValidateConfig(cfg *Config) error {
	return validateConfig(cfg)
}

// GetString returns string configuration value
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns integer configuration value
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns boolean configuration value
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// IsSet checks if value is set
func IsSet(key string) bool {
	return viper.IsSet(key)
}

// logConfig logs loaded configuration (without secrets)
func logConfig(config *Config) {
	log.Printf("ClickHouse: %s@%s:%d/%s (SSL: %s, Compress: %t)",
		config.Database.User,
		config.Database.Host,
		config.Database.Port,
		config.Database.Database,
		config.Database.SSLMode,
		config.Database.Compress)

	log.Printf("Server: %s:%s", config.Server.Host, config.Server.Port)
	log.Printf("JWT: expire_hours=%d, secret_set=%t",
		config.JWT.ExpireHours,
		config.JWT.Secret != "" && config.JWT.Secret != "your-super-secret-jwt-key-here")
	log.Printf("Logging: level=%s", config.Logging.Level)
}

// LoadFromEnv loads configuration only from environment variables
func LoadFromEnv() *Config {
	// Reset Viper
	viper.Reset()

	// Setup for reading only ENV
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	setDefaults()

	// Create configuration structure
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config from ENV: %v", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}

	log.Println("Configuration loaded from environment variables")
	logConfig(&config)
	return &config
}

// LoadFromFile loads configuration only from file
func LoadFromFile(filename string) *Config {
	// Reset Viper
	viper.Reset()

	// Setup for reading file
	viper.SetConfigFile(filename)

	// Set default values
	setDefaults()

	// Read file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file %s: %v", filename, err)
	}

	// Create configuration structure
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}

	log.Printf("Configuration loaded from file: %s", filename)
	logConfig(&config)
	return &config
}

// GetConfigSource returns configuration source
func GetConfigSource() string {
	if viper.ConfigFileUsed() != "" {
		return "file: " + viper.ConfigFileUsed()
	}
	return "environment variables and defaults"
}
