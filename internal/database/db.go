package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitDB() (*pgxpool.Pool, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://myappuser:password@localhost:5432/gmailscraper?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	fmt.Println("Successfully connected to database")
	return pool, nil
}

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
