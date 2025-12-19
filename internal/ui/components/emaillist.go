package components

import (
	"log"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/HoustonMiles/gmailScraper/internal/database"
	"github.com/HoustonMiles/gmailScraper/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmailList struct {
	Container    *fyne.Container
	db           *pgxpool.Pool
	emailView    *EmailView
	senderGroups []*SenderGroup
	app          interface{}
}

func NewEmailList(db *pgxpool.Pool, emailView *EmailView, app interface{}) *EmailList {
	el := &EmailList{
		db:           db,
		emailView:    emailView,
		senderGroups: []*SenderGroup{},
		app:          app,
	}

	el.Container = container.NewVBox()
	el.LoadAllEmails("date_newest")

	return el
}

func (el *EmailList) LoadAllEmails(sortBy string) {
	emails, err := database.GetAllEmails(el.db, sortBy)
	if err != nil {
		log.Printf("Error loading emails: %v", err)
		return
	}

	el.groupAndDisplay(emails)
}

func (el *EmailList) LoadEmailsBySender(sender string, sortBy string) {
	emails, err := database.GetEmailsByFrom(el.db, sender, sortBy)
	if err != nil {
		log.Printf("Error loading emails: %v", err)
		return
	}

	el.groupAndDisplay(emails)
}

func (el *EmailList) groupAndDisplay(emails []models.Email) {
	// Group emails by sender
	grouped := make(map[string][]models.Email)
	for _, email := range emails {
		grouped[email.From] = append(grouped[email.From], email)
	}

	// Clear existing groups
	el.senderGroups = []*SenderGroup{}
	el.Container.Objects = []fyne.CanvasObject{}

	// Sort senders alphabetically
	senders := make([]string, 0, len(grouped))
	for sender := range grouped {
		senders = append(senders, sender)
	}
	sort.Strings(senders)

	// Create a sender group for each sender
	for _, sender := range senders {
		senderEmails := grouped[sender]
		group := NewSenderGroup(sender, senderEmails, el.emailView)
		el.senderGroups = append(el.senderGroups, group)
		
		el.Container.Add(group.Container)
		el.Container.Add(widget.NewSeparator())
	}

	el.Container.Refresh()
}

func (el *EmailList) GetSelectedIDs() []string {
	var allSelected []string
	for _, group := range el.senderGroups {
		allSelected = append(allSelected, group.GetSelectedIDs()...)
	}
	return allSelected
}
