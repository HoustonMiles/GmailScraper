package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"gmailScraper/internal/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func FetchEmail(client *gmail.Service) ([]models.Email, error) {
	// Create Gmail service
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	// Get emails
	user := "me"
	r, err := srv.Users.Messages.List(user).MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	var senders []string

	// Print email details
	for _, msg := range r.Messages {
		message, err := srv.Users.Messages.Get(user, msg.Id).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve message %s: %v", msg.Id, err)
			continue
		}
		//fmt.Printf("%s", msg)
		for _, header := range message.Payload.Headers {
			if header.Name == "From" {
				senders = append(senders, header.Value)
				break
			}
		}
	}
}
