package main

import (
	"log"
	"fmt"

	"github.com/HoustonMiles/gmailScraper/internal/database"
	"github.com/HoustonMiles/gmailScraper/internal/gmail"
	"github.com/HoustonMiles/gmailScraper/internal/ui"
)

func main() {
	fmt.Println("Starting application...")
	
	// Initialize database
	fmt.Println("Connecting to database...")
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables if they don't exist
	fmt.Println("Creating tables...")
	err = database.CreateTables(db)
	if err != nil {
		log.Fatal(err)
	}

	// Get Gmail client
	fmt.Println("Getting Gmail client...")
	client, err := gmail.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	// Launch UI
	fmt.Println("Launching UI...")
	app := ui.NewApp(db, client)
	app.Run()
}
