package ui

import (
	"fmt"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/HoustonMiles/gmailScraper/internal/database"
	"github.com/HoustonMiles/gmailScraper/internal/ui/components"
	"github.com/HoustonMiles/gmailScraper/internal/ui/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	fyneApp     fyne.App
	mainWindow  fyne.Window
	db          *pgxpool.Pool
	gmailClient *http.Client
	
	emailList   *components.EmailList
	emailView   *components.EmailView
	senderList  *widget.Select
	viewMode    string // "all" or "sender"
}

func NewApp(db *pgxpool.Pool, gmailClient *http.Client) *App {
	a := &App{
		fyneApp:     app.New(),
		db:          db,
		gmailClient: gmailClient,
		viewMode:    "all",
	}
	
	a.mainWindow = a.fyneApp.NewWindow("Gmail Manager")
	a.mainWindow.Resize(fyne.NewSize(1000, 600))
	
	a.setupUI()
	
	return a
}

func (a *App) setupUI() {
	// Create components
	a.emailView = components.NewEmailView()
	a.emailList = components.NewEmailList(a.db, a.emailView, a)
	
	// Load senders for dropdown
	senders, err := database.GetAllSenders(a.db)
	if err != nil {
		log.Printf("Error loading senders: %v", err)
		senders = []string{}
	}
	
	// Add "All Emails" option
	senderOptions := append([]string{"All Emails"}, senders...)
	
	a.senderList = widget.NewSelect(senderOptions, func(selected string) {
		if selected == "All Emails" {
			a.viewMode = "all"
			a.emailList.LoadAllEmails()
		} else {
			a.viewMode = "sender"
			a.emailList.LoadEmailsBySender(selected)
		}
	})
	a.senderList.SetSelected("All Emails")
	
	// Create toolbar
	toolbar := a.createToolbar()
	
	// Create split view
	split := container.NewHSplit(
		a.emailList.Container,
		a.emailView.Container,
	)
	split.SetOffset(0.4)
	
	// Main layout
	content := container.NewBorder(
		toolbar,  // top
		nil,      // bottom
		nil,      // left
		nil,      // right
		split,    // center
	)
	
	a.mainWindow.SetContent(content)
}

func (a *App) createToolbar() *fyne.Container {
	// Sync button
	syncBtn := widget.NewButton("Sync Emails", func() {
		a.syncEmails()
	})
	
	// Delete selected button
	deleteBtn := widget.NewButton("Delete Selected", func() {
		a.deleteSelected()
	})
	
	// Refresh button
	refreshBtn := widget.NewButton("Refresh", func() {
		a.refreshView()
	})
	
	return container.NewHBox(
		syncBtn,
		deleteBtn,
		refreshBtn,
		widget.NewLabel("Filter by sender:"),
		a.senderList,
	)
}

func (a *App) syncEmails() {
	// Show progress dialog
	progress := dialog.NewProgressInfinite("Syncing", "Fetching emails from Gmail...", a.mainWindow)
	progress.Show()
	
	go func() {
		err := handlers.SyncEmails(a.gmailClient, a.db, 50) // Fetch 50 at a time
		progress.Hide()
		
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
		} else {
			dialog.ShowInformation("Success", "Emails synced successfully!", a.mainWindow)
			a.refreshView()
		}
	}()
}

func (a *App) deleteSelected() {
	selectedIDs := a.emailList.GetSelectedIDs()
	if len(selectedIDs) == 0 {
		dialog.ShowInformation("No Selection", "Please select emails to delete", a.mainWindow)
		return
	}
	
	msg := fmt.Sprintf("Are you sure you want to delete %d email(s)?", len(selectedIDs))
	dialog.ShowConfirm("Confirm Delete", msg, func(confirmed bool) {
		if confirmed {
			err := database.DeleteEmails(a.db, selectedIDs)
			if err != nil {
				dialog.ShowError(err, a.mainWindow)
			} else {
				dialog.ShowInformation("Success", fmt.Sprintf("Deleted %d emails", len(selectedIDs)), a.mainWindow)
				a.refreshView()
			}
		}
	}, a.mainWindow)
}

func (a *App) refreshView() {
	// Reload senders
	senders, err := database.GetAllSenders(a.db)
	if err != nil {
		log.Printf("Error loading senders: %v", err)
	} else {
		senderOptions := append([]string{"All Emails"}, senders...)
		a.senderList.Options = senderOptions
		a.senderList.Refresh()
	}
	
	// Reload email list
	if a.viewMode == "all" {
		a.emailList.LoadAllEmails()
	} else if a.senderList.Selected != "All Emails" {
		a.emailList.LoadEmailsBySender(a.senderList.Selected)
	}
}

func (a *App) Run() {
	a.mainWindow.ShowAndRun()
}
