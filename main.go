package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Read credentials.json file
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials.json: %v", err)
	}

	// Create OAuth config
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse credentials: %v", err)
	}

	// Get token
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	client := config.Client(ctx, tok)

	// Create Gmail service
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	// Get first email
	user := "me"
	r, err := srv.Users.Messages.List(user).MaxResults(1).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	if len(r.Messages) == 0 {
		fmt.Println("No messages found.")
		return
	}

	// Get full message details
	msg, err := srv.Users.Messages.Get(user, r.Messages[0].Id).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve message: %v", err)
	}

	// Print email details
	fmt.Println("First Email:")
	fmt.Println("ID:", msg.Id)
	for _, header := range msg.Payload.Headers {
		if header.Name == "From" || header.Name == "Subject" || header.Name == "Date" {
			fmt.Printf("%s: %s\n", header.Name, header.Value)
		}
	}
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to this link in your browser:\n%v\n\n", authURL)
	fmt.Print("Enter authorization code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving token to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
