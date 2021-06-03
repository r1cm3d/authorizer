package internal

import (
	"time"
)

type (
	Time time.Time
	Timer interface {
		Now() time.Time
	}
	Violation string
	Account   struct {
		ActiveCard     bool `json:"active-card"`
		AvailableLimit int `json:"available-limit"`
	}
	Transaction struct {
		Merchant string `json:"merchant"`
		Amount   int `json:"amount"`
		Time Time `json:"time"`
	}
	Event struct {
		*Account `json:"account"`
		*Transaction `json:"transaction"`
	}
	OutputEvent struct {
		Event
		Violations []Violation
	}
)

func (e Event) isTransaction() bool {
	return e.Transaction != nil
}

func (e OutputEvent) hasViolation() bool {
	return len(e.Violations) > 0
}



