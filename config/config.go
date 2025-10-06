package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	App      AppConfig
	Security SecurityConfig
	Database DatabaseConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port    string
	GinMode string
}

// AppConfig holds application metadata
type AppConfig struct {
	Name    string
	Version string
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	APISecret     string
	JWTSecret     string
	JWTExpiration int // in hours
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Driver   string
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		App: AppConfig{
			Name:    getEnv("APP_NAME", "GateHide API"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
		Security: SecurityConfig{
			APISecret:     getEnv("API_SECRET", "default-secret-key"),
			JWTSecret:     getEnv("JWT_SECRET", "jwt-secret-key-change-in-production"),
			JWTExpiration: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "gatehide"),
			SSLMode:  getEnv("DB_SSLMODE", "false"),
			Driver:   getEnv("DB_DRIVER", "mysql"),
		},
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves an environment variable as integer or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	switch c.Database.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.DBName,
		)
	case "postgres":
		sslMode := "disable"
		if c.Database.SSLMode == "true" {
			sslMode = "require"
		}
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Host,
			c.Database.Port,
			c.Database.User,
			c.Database.Password,
			c.Database.DBName,
			sslMode,
		)
	default:
		return ""
	}
}

// GetServerDSN returns the database server connection string (without database name)
func (c *Config) GetServerDSN() string {
	switch c.Database.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&loc=Local",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
		)
	case "postgres":
		sslMode := "disable"
		if c.Database.SSLMode == "true" {
			sslMode = "require"
		}
		return fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s",
			c.Database.Host,
			c.Database.Port,
			c.Database.User,
			c.Database.Password,
			sslMode,
		)
	default:
		return ""
	}
}
