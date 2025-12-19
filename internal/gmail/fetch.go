package gmail

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HoustonMiles/gmailScraper/internal/models"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// FetchEmails retrieves emails from Gmail with optional limit
// If maxResults is 0, it fetches ALL emails (with pagination)
func FetchEmails(client *http.Client, maxResults int64) ([]models.Email, error) {
	ctx := context.Background()

	// Create Gmail service
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail service: %v", err)
	}

	var allEmails []models.Email
	pageToken := ""
	user := "me"

	// If maxResults is 0, fetch ALL emails
	fetchAll := maxResults == 0
	
	for {
		// Gmail API max is 500 per request
		batchSize := int64(500)
		if !fetchAll && maxResults < 500 {
			batchSize = maxResults
		}

		// Build the request
		req := srv.Users.Messages.List(user).MaxResults(batchSize)
		if pageToken != "" {
			req = req.PageToken(pageToken)
		}

		// Execute the request
		r, err := req.Do()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve messages: %v", err)
		}

		fmt.Printf("Fetched %d message IDs (total so far: %d)...\n", len(r.Messages), len(allEmails)+len(r.Messages))

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

			allEmails = append(allEmails, email)

			// If we have a specific limit and reached it, stop
			if !fetchAll && int64(len(allEmails)) >= maxResults {
				return allEmails, nil
			}
		}

		// Check if there are more pages
		pageToken = r.NextPageToken
		if pageToken == "" {
			// No more pages
			break
		}

		// If not fetching all and we got enough, stop
		if !fetchAll && int64(len(allEmails)) >= maxResults {
			break
		}
	}

	fmt.Printf("Finished fetching. Total emails: %d\n", len(allEmails))
	return allEmails, nil
}

// FetchAllEmails is a convenience function to fetch all emails
func FetchAllEmails(client *http.Client) ([]models.Email, error) {
	return FetchEmails(client, 0)
}
