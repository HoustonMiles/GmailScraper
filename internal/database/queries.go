package database

import (
	"context"
	"fmt"

	"github.com/HoustonMiles/gmailScraper/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SaveEmails saves a batch of emails to the database
func SaveEmails(pool *pgxpool.Pool, emails []models.Email) error {
	ctx := context.Background()

	for _, email := range emails {
		// Use UPSERT to avoid duplicate key errors
		query := `
		INSERT INTO emails (id, from_address, subject, body, date_received)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			from_address = EXCLUDED.from_address,
			subject = EXCLUDED.subject,
			body = EXCLUDED.body,
			date_received = EXCLUDED.date_received
		`

		_, err := pool.Exec(ctx, query,
			email.ID,
			email.From,
			email.Subject,
			email.Body,
			email.Date,
		)

		if err != nil {
			return fmt.Errorf("error saving email %s: %v", email.ID, err)
		}
	}

	fmt.Printf("Successfully saved %d emails\n", len(emails))
	return nil
}

// GetAllEmails retrieves all emails from the database
func GetAllEmails(pool *pgxpool.Pool) ([]models.Email, error) {
	ctx := context.Background()

	query := `
	SELECT id, from_address, subject, body, date_received
	FROM emails
	ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying emails: %v", err)
	}
	defer rows.Close()

	var emails []models.Email
	for rows.Next() {
		var email models.Email
		err := rows.Scan(
			&email.ID,
			&email.From,
			&email.Subject,
			&email.Body,
			&email.Date,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning email: %v", err)
		}
		emails = append(emails, email)
	}

	return emails, nil
}

// GetEmailsByFrom retrieves emails from a specific sender
func GetEmailsByFrom(pool *pgxpool.Pool, fromAddress string) ([]models.Email, error) {
	ctx := context.Background()

	query := `
	SELECT id, from_address, subject, body, date_received
	FROM emails
	WHERE from_address LIKE $1
	ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query, "%"+fromAddress+"%")
	if err != nil {
		return nil, fmt.Errorf("error querying emails: %v", err)
	}
	defer rows.Close()

	var emails []models.Email
	for rows.Next() {
		var email models.Email
		err := rows.Scan(
			&email.ID,
			&email.From,
			&email.Subject,
			&email.Body,
			&email.Date,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning email: %v", err)
		}
		emails = append(emails, email)
	}

	return emails, nil
}
