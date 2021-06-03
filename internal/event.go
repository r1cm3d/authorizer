package internal

import (
	"encoding/json"
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
		ActiveCard     bool `json:"active-card,omitempty"`
		AvailableLimit int `json:"available-limit,omitempty"`
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

func Parse(input string) Event {
	var ie Event
	json.Unmarshal([]byte(input), &ie)

	return ie
}

func (it *Time) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, _ := time.Parse(time.RFC3339, s)

	*it = Time(t)
	return nil
}

func (oe OutputEvent) String() string {
	//TODO: extract it to top of the file
	type Account struct {
		ActiveCard     *bool `json:"active-card,omitempty"`
		AvailableLimit *int `json:"available-limit,omitempty"`
	}
	type output struct {
		Account `json:"account"`
		Violations []Violation `json:"violations"`
	}
	op := output{
		Account:   Account{
			ActiveCard:     nil,
			AvailableLimit: nil,
		},
		Violations: make([]Violation, 0),
	}

	if oe.Account != nil {
		op.ActiveCard = &oe.ActiveCard
		op.AvailableLimit = &oe.AvailableLimit
	}

	if oe.hasViolation() {
		op.Violations = oe.Violations
	}

	str, _ := json.Marshal(op)

	return string(str)
}

func (e Event) isTransaction() bool {
	return e.Transaction != nil
}

func (oe OutputEvent) hasViolation() bool {
	return len(oe.Violations) > 0
}