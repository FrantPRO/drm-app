package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	Pool *pgxpool.Pool
	DB   *sqlx.DB
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func LoadConfig() (*Config, error) {
	portStr := getEnv("DB_PORT", "5432")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT value: %s", portStr)
	}

	config := &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "postgres"),
	}

	// Validate required fields
	if config.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	return config, nil
}

func NewDatabase() (*Database, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	dsn := buildDSN(config)

	// Try to connect with retries
	maxRetries := 5
	retryDelay := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// Create connection pool
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			cancel()
			if attempt == maxRetries {
				return nil, fmt.Errorf("failed to create connection pool after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(retryDelay)
			continue
		}

		// Test connection
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			cancel()
			if attempt == maxRetries {
				return nil, fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(retryDelay)
			continue
		}

		// Create sqlx connection for easier queries
		db, err := sqlx.Open("pgx", dsn)
		if err != nil {
			pool.Close()
			cancel()
			if attempt == maxRetries {
				return nil, fmt.Errorf("failed to open sqlx connection after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(retryDelay)
			continue
		}

		// Test sqlx connection
		if err := db.PingContext(ctx); err != nil {
			pool.Close()
			db.Close()
			cancel()
			if attempt == maxRetries {
				return nil, fmt.Errorf("failed to ping database via sqlx after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(retryDelay)
			continue
		}

		cancel()
		return &Database{
			Pool: pool,
			DB:   db,
		}, nil
	}

	return nil, fmt.Errorf("exhausted all %d connection attempts", maxRetries)
}

func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
	if d.DB != nil {
		d.DB.Close()
	}
}

func (d *Database) Health(ctx context.Context) error {
	return d.Pool.Ping(ctx)
}

func buildDSN(config *Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
