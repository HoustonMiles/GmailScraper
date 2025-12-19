package components

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/HoustonMiles/gmailScraper/internal/database"
	"github.com/HoustonMiles/gmailScraper/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmailList struct {
	Container  *fyne.Container
	db         *pgxpool.Pool
	emailView  *EmailView
	emails     []models.Email
	list       *widget.List
	checkboxes map[int]*widget.Check
	app        interface{} // Reference to parent app for refresh
}

func NewEmailList(db *pgxpool.Pool, emailView *EmailView, app interface{}) *EmailList {
	el := &EmailList{
		db:         db,
		emailView:  emailView,
		emails:     []models.Email{},
		checkboxes: make(map[int]*widget.Check),
		app:        app,
	}
	
	el.list = widget.NewList(
		func() int {
			return len(el.emails)
		},
		func() fyne.CanvasObject {
			check := widget.NewCheck("", nil)
			label := widget.NewLabel("Template")
			return container.NewBorder(nil, nil, check, nil, label)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			c := obj.(*fyne.Container)
			check := c.Objects[0].(*widget.Check)
			label := c.Objects[1].(*widget.Label)
			
			email := el.emails[id]
			label.SetText(fmt.Sprintf("%s - %s", email.From, email.Subject))
			
			// Store checkbox reference
			el.checkboxes[id] = check
			
			// Clear previous callback
			check.OnChanged = nil
			check.Checked = false
			check.Refresh()
		},
	)
	
	el.list.OnSelected = func(id widget.ListItemID) {
		el.emailView.ShowEmail(el.emails[id])
	}
	
	el.Container = container.NewBorder(
		widget.NewLabel("Emails"),
		nil, nil, nil,
		el.list,
	)
	
	el.LoadAllEmails()
	
	return el
}

func (el *EmailList) LoadAllEmails() {
	emails, err := database.GetAllEmails(el.db)
	if err != nil {
		log.Printf("Error loading emails: %v", err)
		return
	}
	el.emails = emails
	el.checkboxes = make(map[int]*widget.Check)
	el.list.Refresh()
}

func (el *EmailList) LoadEmailsBySender(sender string) {
	emails, err := database.GetEmailsByFrom(el.db, sender)
	if err != nil {
		log.Printf("Error loading emails: %v", err)
		return
	}
	el.emails = emails
	el.checkboxes = make(map[int]*widget.Check)
	el.list.Refresh()
}

func (el *EmailList) GetSelectedIDs() []string {
	var selectedIDs []string
	for idx, check := range el.checkboxes {
		if check.Checked {
			selectedIDs = append(selectedIDs, el.emails[idx].ID)
		}
	}
	return selectedIDs
}
