package internal

import (
	"strings"
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

func (it *Time) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, _ := time.Parse(time.RFC3339, s)

	*it = Time(t)
	return nil
}


