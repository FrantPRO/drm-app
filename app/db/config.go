package db

import (
	"log"
)

// LogConfig logs the current database configuration (without password)
func LogConfig() {
	config, err := LoadConfig()
	if err != nil {
		log.Printf("Failed to load database config: %v", err)
		return
	}

	log.Printf("Database Configuration:")
	log.Printf("  Host: %s", config.Host)
	log.Printf("  Port: %d", config.Port)
	log.Printf("  User: %s", config.User)
	log.Printf("  Database: %s", config.DBName)
	log.Printf("  Password: %s", maskPassword(config.Password))
}

// ValidateConfig checks if all required environment variables are set
func ValidateConfig() error {
	_, err := LoadConfig()
	return err
}

func maskPassword(password string) string {
	if len(password) == 0 {
		return "<empty>"
	}
	if len(password) <= 3 {
		return "***"
	}
	return password[:2] + "***"
}
