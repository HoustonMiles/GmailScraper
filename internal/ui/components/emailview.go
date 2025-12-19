package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/HoustonMiles/gmailScraper/internal/models"
)

type EmailView struct {
	Container *fyne.Container
	fromLabel *widget.Label
	dateLabel *widget.Label
	subject   *widget.Label
	body      *widget.Label
}

func NewEmailView() *EmailView {
	ev := &EmailView{
		fromLabel: widget.NewLabel("From: "),
		dateLabel: widget.NewLabel("Date: "),
		subject:   widget.NewLabel("Subject: "),
		body:      widget.NewLabel("Select an email to view"),
	}
	
	ev.body.Wrapping = fyne.TextWrapWord
	
	ev.Container = container.NewBorder(
		container.NewVBox(
			ev.fromLabel,
			ev.dateLabel,
			ev.subject,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		container.NewScroll(ev.body),
	)
	
	return ev
}

func (ev *EmailView) ShowEmail(email models.Email) {
	ev.fromLabel.SetText("From: " + email.From)
	ev.dateLabel.SetText("Date: " + email.Date)
	ev.subject.SetText("Subject: " + email.Subject)
	ev.body.SetText(email.Body)
}
