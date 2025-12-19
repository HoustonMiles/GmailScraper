package database

import (
	"context"
	"fmt"

	"github.com/HoustonMiles/gmailScraper/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SaveEmails saves emails to database
func SaveEmails(pool *pgxpool.Pool, emails []models.Email) error {
	ctx := context.Background()

	for _, email := range emails {
		query := `
		INSERT INTO emails (id, from_address, subject, body, date_received)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			from_address = EXCLUDED.from_address,
			subject = EXCLUDED.subject,
			body = EXCLUDED.body,
			date_received = EXCLUDED.date_received
		`

		_, err := pool.Exec(ctx, query, email.ID, email.From, email.Subject, email.Body, email.Date)
		if err != nil {
			return fmt.Errorf("error saving email %s: %v", email.ID, err)
		}
	}

	fmt.Printf("Successfully saved %d emails\n", len(emails))
	return nil
}

// GetAllEmails gets all emails with optional sorting
func GetAllEmails(pool *pgxpool.Pool, sortBy string) ([]models.Email, error) {
	ctx := context.Background()

	// Determine sort order
	var orderClause string
	switch sortBy {
	case "date_newest":
		orderClause = "ORDER BY date_received DESC"
	case "date_oldest":
		orderClause = "ORDER BY date_received ASC"
	case "sender_asc":
		orderClause = "ORDER BY from_address ASC"
	case "sender_desc":
		orderClause = "ORDER BY from_address DESC"
	default:
		orderClause = "ORDER BY created_at DESC"
	}

	query := fmt.Sprintf(`
	SELECT id, from_address, subject, body, date_received
	FROM emails
	%s
	`, orderClause)

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying emails: %v", err)
	}
	defer rows.Close()

	var emails []models.Email
	for rows.Next() {
		var email models.Email
		err := rows.Scan(&email.ID, &email.From, &email.Subject, &email.Body, &email.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning email: %v", err)
		}
		emails = append(emails, email)
	}

	return emails, nil
}

// GetEmailsByFrom retrieves emails from a specific sender with sorting
func GetEmailsByFrom(pool *pgxpool.Pool, fromAddress string, sortBy string) ([]models.Email, error) {
	ctx := context.Background()

	// Determine sort order
	var orderClause string
	switch sortBy {
	case "date_newest":
		orderClause = "ORDER BY date_received DESC"
	case "date_oldest":
		orderClause = "ORDER BY date_received ASC"
	case "sender_asc":
		orderClause = "ORDER BY from_address ASC"
	case "sender_desc":
		orderClause = "ORDER BY from_address DESC"
	default:
		orderClause = "ORDER BY created_at DESC"
	}

	query := fmt.Sprintf(`
	SELECT id, from_address, subject, body, date_received
	FROM emails
	WHERE from_address LIKE $1
	%s
	`, orderClause)

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

// GetEmailsBySender gets all emails from a specific sender with sorting
func GetEmailsBySender(pool *pgxpool.Pool, sender string, sortBy string) ([]models.Email, error) {
	return GetEmailsByFrom(pool, sender, sortBy)
}

// GetAllSenders gets a list of unique senders
func GetAllSenders(pool *pgxpool.Pool) ([]string, error) {
	ctx := context.Background()

	query := `
	SELECT DISTINCT from_address
	FROM emails
	ORDER BY from_address
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying senders: %v", err)
	}
	defer rows.Close()

	var senders []string
	for rows.Next() {
		var sender string
		err := rows.Scan(&sender)
		if err != nil {
			return nil, fmt.Errorf("error scanning sender: %v", err)
		}
		senders = append(senders, sender)
	}

	return senders, nil
}

// DeleteEmail deletes a single email by ID
func DeleteEmail(pool *pgxpool.Pool, emailID string) error {
	ctx := context.Background()

	query := `DELETE FROM emails WHERE id = $1`
	
	result, err := pool.Exec(ctx, query, emailID)
	if err != nil {
		return fmt.Errorf("error deleting email: %v", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("email not found")
	}

	fmt.Printf("Deleted email: %s\n", emailID)
	return nil
}

// DeleteEmails deletes multiple emails by their IDs
func DeleteEmails(pool *pgxpool.Pool, emailIDs []string) error {
	ctx := context.Background()

	for _, id := range emailIDs {
		query := `DELETE FROM emails WHERE id = $1`
		_, err := pool.Exec(ctx, query, id)
		if err != nil {
			return fmt.Errorf("error deleting email %s: %v", id, err)
		}
	}

	fmt.Printf("Deleted %d emails\n", len(emailIDs))
	return nil
}

// DeleteEmailsBySender deletes all emails from a specific sender
func DeleteEmailsBySender(pool *pgxpool.Pool, sender string) error {
	ctx := context.Background()

	query := `DELETE FROM emails WHERE from_address LIKE $1`
	
	result, err := pool.Exec(ctx, query, "%"+sender+"%")
	if err != nil {
		return fmt.Errorf("error deleting emails: %v", err)
	}

	rowsAffected := result.RowsAffected()
	fmt.Printf("Deleted %d emails from %s\n", rowsAffected, sender)
	return nil
}
