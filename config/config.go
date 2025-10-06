package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	App      AppConfig
	Security SecurityConfig
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
	APISecret string
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
			APISecret: getEnv("API_SECRET", "default-secret-key"),
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
