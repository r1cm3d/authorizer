package internal

import (
	"encoding/json"
	"strings"
	"time"
)

type (
	Time  time.Time
	Violation string
	Account   struct {
		ActiveCard     bool `json:"active-card"`
		AvailableLimit int  `json:"available-limit"`
	}
	Transaction struct {
		Merchant string `json:"merchant"`
		Amount   int    `json:"amount"`
		Time     Time   `json:"time"`
	}
	Event struct {
		*Account     `json:"account"`
		*Transaction `json:"transaction"`
	}
	TimelineEvent struct {
		Event
		Violations []Violation
	}
	OutputAccount struct {
		ActiveCard     *bool `json:"active-card,omitempty"`
		AvailableLimit *int  `json:"available-limit,omitempty"`
	}
	Output struct {
		OutputAccount `json:"account"`
		Violations []Violation   `json:"violations"`
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

func (te TimelineEvent) String() string {
	//TODO: extract it to top of the file

	op := Output{
		OutputAccount: OutputAccount{
			ActiveCard:     nil,
			AvailableLimit: nil,
		},
		Violations: make([]Violation, 0),
	}

	if te.Account != nil {
		op.ActiveCard = &te.ActiveCard
		op.AvailableLimit = &te.AvailableLimit
	}

	if te.hasViolation() {
		op.Violations = te.Violations
	}

	str, _ := json.Marshal(op)

	return string(str)
}

func (e Event) isTransaction() bool {
	return e.Transaction != nil
}

func (te TimelineEvent) hasViolation() bool {
	return len(te.Violations) > 0
}
