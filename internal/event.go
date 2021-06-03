package internal

import (
	"encoding/json"
	"strings"
	"time"
)

type (
	Account   struct {
		ActiveCard     bool `json:"active-card"`
		AvailableLimit int  `json:"available-limit"`
	}
	Transaction struct {
		Merchant string   `json:"merchant"`
		Amount   int      `json:"amount"`
		Time     datetime `json:"time"`
	}
	Event struct {
		*Account     `json:"Account"`
		*Transaction `json:"Transaction"`
	}
	TimelineEvent struct {
		Event
		Violations []violation
	}

	datetime      time.Time
	violation     string
	outputAccount struct {
		ActiveCard     *bool `json:"active-card,omitempty"`
		AvailableLimit *int  `json:"available-limit,omitempty"`
	}
	output struct {
		outputAccount `json:"Account"`
		Violations    []violation `json:"violations"`
	}
)

func Parse(input string) Event {
	var ie Event
	json.Unmarshal([]byte(input), &ie)

	return ie
}

func (it *datetime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, _ := time.Parse(time.RFC3339, s)

	*it = datetime(t)
	return nil
}

func (te TimelineEvent) String() string {
	op := output{
		outputAccount: outputAccount{
			ActiveCard:     nil,
			AvailableLimit: nil,
		},
		Violations: make([]violation, 0),
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
