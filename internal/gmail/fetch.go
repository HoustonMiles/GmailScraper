package gmail

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HoustonMiles/gmailScraper/internal/models"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// FetchEmails retrieves emails from Gmail
func FetchEmails(client *http.Client, maxResults int64) ([]models.Email, error) {
	ctx := context.Background()

	// Create Gmail service
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail service: %v", err)
	}

	// Get emails
	user := "me"
	r, err := srv.Users.Messages.List(user).MaxResults(maxResults).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve messages: %v", err)
	}

	var emails []models.Email

	// Process each message
	for _, msg := range r.Messages {
		message, err := srv.Users.Messages.Get(user, msg.Id).Format("full").Do()
		if err != nil {
			fmt.Printf("Unable to retrieve message %s: %v\n", msg.Id, err)
			continue
		}

		email := models.Email{
			ID: msg.Id,
		}

		// Extract headers
		for _, header := range message.Payload.Headers {
			switch header.Name {
			case "From":
				email.From = header.Value
			case "Subject":
				email.Subject = header.Value
			case "Date":
				email.Date = header.Value
			}
		}

		// Extract body (simplified - gets snippet)
		email.Body = message.Snippet

		emails = append(emails, email)
	}

	return emails, nil
}
