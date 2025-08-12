package config

import (
	"os"
	"strconv"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Host string
	Port string

	// Database configuration
	DatabaseURL string

	// JWT configuration
	JWTSecret     string
	JWTExpiration time.Duration

	// Application configuration
	Environment string
	Debug       bool

	// CORS configuration
	CORSAllowOrigins []string
	CORSAllowMethods []string
	CORSAllowHeaders []string

	// Rate limiting
	RateLimitMax      int
	RateLimitDuration time.Duration
	Midtrans          MidtransConfig
}

type MidtransConfig struct {
	ServerKey   string
	ClientKey   string
	Environment midtrans.EnvironmentType
	Client      snap.Client
}

// Load environment variables
func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

// LoadConfig loads configuration from environment variables with defaults

func LoadConfig() *Config {
	// Initialize Midtrans environment
	midtransEnv := midtrans.Sandbox
	if getEnv("MIDTRANS_ENV", "sandbox") == "production" {
		midtransEnv = midtrans.Production
	}

	// Initialize Snap client
	snapClient := snap.Client{}
	snapClient.New(getEnv("MIDTRANS_SERVER_KEY", ""), midtransEnv)

	config := &Config{
		// Server defaults
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnv("PORT", "8000"),

		// Database defaults
		DatabaseURL: getEnv("DATABASE_URL", ""),

		// JWT defaults
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),

		// Application defaults
		Environment: getEnv("ENVIRONMENT", ""),
		Debug:       getBoolEnv("DEBUG", true),

		// CORS defaults
		CORSAllowOrigins: []string{"*"},
		CORSAllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		CORSAllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},

		// Rate limiting defaults
		RateLimitMax:      getIntEnv("RATE_LIMIT_MAX", 100),
		RateLimitDuration: getDurationEnv("RATE_LIMIT_DURATION", time.Minute),
		// Midtrans configuration
		Midtrans: MidtransConfig{
			ServerKey:   getEnv("MIDTRANS_SERVER_KEY", ""),
			ClientKey:   getEnv("MIDTRANS_CLIENT_KEY", ""),
			Environment: midtransEnv,
			Client:      snapClient,
		},
	}

	return config
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getBoolEnv gets a boolean environment variable with a default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if result, err := strconv.ParseBool(value); err == nil {
			return result
		}
	}
	return defaultValue
}

// getIntEnv gets an integer environment variable with a default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if result, err := strconv.Atoi(value); err == nil {
			return result
		}
	}
	return defaultValue
}

// getDurationEnv gets a duration environment variable with a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if result, err := time.ParseDuration(value); err == nil {
			return result
		}
	}
	return defaultValue
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
