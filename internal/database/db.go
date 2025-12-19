package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// InitDB initializes the database connection pool
func InitDB() (*pgxpool.Pool, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Get connection string from environment
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	// Test the connection
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	fmt.Println("Successfully connected to database")
	return pool, nil
}

// CreateTables creates the necessary database tables
func CreateTables(pool *pgxpool.Pool) error {
	ctx := context.Background()

	query := `
	CREATE TABLE IF NOT EXISTS emails (
		id VARCHAR(255) PRIMARY KEY,
		from_address TEXT NOT NULL,
		subject TEXT,
		body TEXT,
		date_received TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_emails_from ON emails(from_address);
	CREATE INDEX IF NOT EXISTS idx_emails_date ON emails(date_received);
	`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	fmt.Println("Tables created successfully")
	return nil
}
