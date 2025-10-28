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
	Server       ServerConfig
	App          AppConfig
	Security     SecurityConfig
	Database     DatabaseConfig
	Notification NotificationConfig
	FileStorage  FileStorageConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host    string
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

// NotificationConfig holds notification-related configuration
type NotificationConfig struct {
	Email EmailConfig
	SMS   SMSConfig
}

// EmailConfig holds email SMTP configuration
type EmailConfig struct {
	Enabled   bool
	SMTPHost  string
	SMTPPort  int
	SMTPUser  string
	SMTPPass  string
	FromEmail string
	FromName  string
	UseTLS    bool
	UseSSL    bool
}

// SMSConfig holds SMS configuration for Kavenegar
type SMSConfig struct {
	Enabled    bool
	APIKey     string
	Sender     string
	TestMode   bool
	MaxRetries int
}

// FileStorageConfig holds file storage configuration
type FileStorageConfig struct {
	UploadPath   string
	MaxFileSize  int64 // in bytes
	AllowedTypes []string
	PublicURL    string
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Host:    getEnv("HOST", "0.0.0.0"),
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
		Notification: NotificationConfig{
			Email: EmailConfig{
				Enabled:   getEnvBool("EMAIL_ENABLED", true),
				SMTPHost:  getEnv("MAIL_HOST", getEnv("SMTP_HOST", "localhost")),
				SMTPPort:  getEnvInt("MAIL_PORT", getEnvInt("SMTP_PORT", 587)),
				SMTPUser:  getEnv("MAIL_USERNAME", getEnv("SMTP_USER", "")),
				SMTPPass:  getEnv("MAIL_PASSWORD", getEnv("SMTP_PASS", "")),
				FromEmail: getEnv("FROM_EMAIL", "noreply@gatehide.com"),
				FromName:  getEnv("FROM_NAME", "GateHide"),
				UseTLS:    getEnvBool("SMTP_USE_TLS", true),
				UseSSL:    getEnvBool("SMTP_USE_SSL", false),
			},
			SMS: SMSConfig{
				Enabled:    getEnvBool("SMS_ENABLED", false),
				APIKey:     getEnv("KAVENEGAR_API_KEY", ""),
				Sender:     getEnv("SMS_SENDER", "10008663"),
				TestMode:   getEnvBool("SMS_TEST_MODE", true),
				MaxRetries: getEnvInt("SMS_MAX_RETRIES", 3),
			},
		},
		FileStorage: FileStorageConfig{
			UploadPath:   getEnv("UPLOAD_PATH", "./uploads"),
			MaxFileSize:  getEnvInt64("MAX_FILE_SIZE", 10*1024*1024), // 10MB default
			AllowedTypes: []string{".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"},
			PublicURL:    getEnv("PUBLIC_URL", "http://localhost:8080"),
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

// getEnvBool retrieves an environment variable as boolean or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvInt64 retrieves an environment variable as int64 or returns a default value
func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	switch c.Database.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&multiStatements=true",
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
