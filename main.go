package main

import (
	"fmt"
	"log"

	"github.com/HoustonMiles/gmailScraper/internal/database"
	"github.com/HoustonMiles/gmailScraper/internal/gmail"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables if they don't exist
	err = database.CreateTables(db)
	if err != nil {
		log.Fatal(err)
	}

	// Get Gmail client
	client, err := gmail.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	// Fetch emails
	fmt.Println("Fetching emails from Gmail...")
	emails, err := gmail.FetchEmails(client, 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched %d emails\n", len(emails))

	// Save to database
	err = database.SaveEmails(db, emails)
	if err != nil {
		log.Fatal(err)
	}

	// Display emails from database
	fmt.Println("\nEmails in database:")
	dbEmails, err := database.GetAllEmails(db)
	if err != nil {
		log.Fatal(err)
	}

	for i, email := range dbEmails {
		fmt.Printf("%d. From: %s | Subject: %s\n", i+1, email.From, email.Subject)
	}
}
