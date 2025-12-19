package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/HoustonMiles/gmailScraper/internal/models"
)

type SenderGroup struct {
	Container   *fyne.Container
	sender      string
	emails      []models.Email
	emailView   *EmailView
	selectAll   *widget.Check
	checkboxes  map[int]*widget.Check
	emailList   *widget.List
}

func NewSenderGroup(sender string, emails []models.Email, emailView *EmailView) *SenderGroup {
	sg := &SenderGroup{
		sender:     sender,
		emails:     emails,
		emailView:  emailView,
		checkboxes: make(map[int]*widget.Check),
	}

	// Select all checkbox
	sg.selectAll = widget.NewCheck(fmt.Sprintf("%s (%d emails)", sender, len(emails)), func(checked bool) {
		// Select/deselect all checkboxes in this group
		for _, check := range sg.checkboxes {
			check.Checked = checked
			check.Refresh()
		}
	})

	// Email list for this sender
	sg.emailList = widget.NewList(
		func() int {
			return len(sg.emails)
		},
		func() fyne.CanvasObject {
			check := widget.NewCheck("", nil)
			label := widget.NewLabel("Template")
			return container.NewHBox(check, label)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			c := obj.(*fyne.Container)
			check := c.Objects[0].(*widget.Check)
			label := c.Objects[1].(*widget.Label)

			email := sg.emails[id]
			label.SetText(email.Subject)

			// Store checkbox reference
			sg.checkboxes[id] = check

			// Clear previous callback
			check.OnChanged = nil
			check.Checked = false
			check.Refresh()
		},
	)

	sg.emailList.OnSelected = func(id widget.ListItemID) {
		sg.emailView.ShowEmail(sg.emails[id])
	}

	// Container with header and list
	sg.Container = container.NewBorder(
		sg.selectAll, // top
		nil, nil, nil,
		sg.emailList, // center
	)

	return sg
}

func (sg *SenderGroup) GetSelectedIDs() []string {
	var selectedIDs []string
	for idx, check := range sg.checkboxes {
		if check.Checked {
			selectedIDs = append(selectedIDs, sg.emails[idx].ID)
		}
	}
	return selectedIDs
}
