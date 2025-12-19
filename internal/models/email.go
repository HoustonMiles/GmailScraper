package models

import "time"

type Email struct {
	ID	string
	From	string
	Subject	string
	Body	string
	Date	time.Time
}
