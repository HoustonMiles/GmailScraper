package handlers

import (
	"fmt"
	"net/http"

	"github.com/HoustonMiles/gmailScraper/internal/database"
	"github.com/HoustonMiles/gmailScraper/internal/gmail"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SyncEmails(gmailClient *http.Client, db *pgxpool.Pool, maxResults int64) error {
	// Fetch emails from Gmail
	emails, err := gmail.FetchEmails(gmailClient, maxResults)
	if err != nil {
		return fmt.Errorf("failed to fetch emails: %v", err)
	}
	
	// Save to database
	err = database.SaveEmails(db, emails)
	if err != nil {
		return fmt.Errorf("failed to save emails: %v", err)
	}
	
	return nil
}
