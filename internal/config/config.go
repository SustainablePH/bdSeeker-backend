package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	JWTSecret        string `mapstructure:"JWT_SECRET"`
	JWTExpiry        string `mapstructure:"JWT_EXPIRY"`
	JWTRefreshExpiry string `mapstructure:"JWT_REFRESH_EXPIRY"`

	ServerPort string `mapstructure:"SERVER_PORT"`
	ServerHost string `mapstructure:"SERVER_HOST"`

	Environment string `mapstructure:"ENV"`
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables and config files using Viper
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (for backward compatibility)
	// This loads variables into the environment, then Viper reads them
	if err := loadEnvFile(".env"); err != nil {
		log.Println("No .env file found, using system environment variables and config files")
	}

	// Set config file properties (optional config file)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set default values
	setDefaults()

	// Read config file (optional - won't error if not found)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, using environment variables and defaults")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Environment variables take precedence over config file
	viper.AutomaticEnv()

	// Unmarshal config into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Validate required fields
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	AppConfig = &config
	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "bdseeker")
	viper.SetDefault("DB_SSLMODE", "disable")

	// JWT defaults
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("JWT_EXPIRY", "24h")
	viper.SetDefault("JWT_REFRESH_EXPIRY", "168h")

	// Server defaults
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")

	// Environment default
	viper.SetDefault("ENV", "development")
}

// validateConfig validates that required configuration fields are set
func validateConfig(c *Config) error {
	if c.JWTSecret == "your-secret-key" && c.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be set in production environment")
	}

	if c.DBPassword == "postgres" && c.Environment == "production" {
		log.Println("WARNING: Using default database password in production is not recommended")
	}

	return nil
}

// GetDSN returns the PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// loadEnvFile loads environment variables from a .env file
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Split on first = only
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			os.Setenv(key, value)
		}
	}
	return scanner.Err()
}
