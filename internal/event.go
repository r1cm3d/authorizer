package internal

import (
	"encoding/json"
	"strings"
	"time"
)

type (
	// Account groups information about an Account.
	Account   struct {
		// ActiveCard when is true indicates that is possible to transact with this Account.
		ActiveCard     bool `json:"active-card"`
		// AvailableLimit indicates how much limit this account can transact.
		AvailableLimit int  `json:"available-limit"`
	}
	// Transaction groups information about an Transaction.
	Transaction struct {
		// Merchant is the name of the Merchant that sent Transaction through acquirer.
		Merchant string   `json:"merchant"`
		// Amount is the value of the Transaction without any cents.
		Amount   int      `json:"amount"`
		// Time is the datetime of the Transaction in UTC.
		Time     datetime `json:"time"`
	}
	// Event represents an input Event.
	Event struct {
		// Account related to the Event.
		*Account     `json:"Account"`
		// Transaction related to the Event.
		*Transaction `json:"Transaction"`
	}
	// TimelineEvent represents each event of the Timeline.
	// It could be a valid Event (Violations empty) or a invalid Event.
	TimelineEvent struct {
		// Event related to the TimelineEvent.
		// It is NOT the same input Event. When it is a invalid Event, Account has the last valid state.
		Event
		// Violations has all Violations of this TimelineEvents.
		// It is never nil. When this is empty, the TimelineEvent is valid.
		Violations []violation
	}

	// datetime is a wrapper type created to implement UnmarshalJSON.
	datetime      time.Time
	// violation is a type created to abstract all constants violations.
	violation     string
	// outputAccount is a structured created to represent a TimelineEvent.
	// This new structure is need because properties must be pointers to be compliance with functional requirements.
	outputAccount struct {
		// ActiveCard when is true indicates that is possible to transact with this Account.
		// When it is nil, must be omitted in JSON.
		ActiveCard     *bool `json:"active-card,omitempty"`
		// AvailableLimit indicates how much limit this account can transact.
		// When it is nil, must be omitted in JSON.
		AvailableLimit *int  `json:"available-limit,omitempty"`
	}
	// output is the output.
	// This new structure is need to avoid print Transaction in standard output.
	output struct {
		outputAccount `json:"Account"`
		// Violations has all Violations of this TimelineEvents.
		// It is never nil.
		Violations    []violation `json:"violations"`
	}
)

// Parse receives a JSON input in string format and parses it into an Event.
// For simplicity, it does not handle invalid input or breach of contract.
func Parse(input string) Event {
	var ie Event
	json.Unmarshal([]byte(input), &ie)

	return ie
}

// UnmarshalJSON receives a []byte datetime and parses it into RFC-3339 datetime standard.
// For simplicity, it does not handle invalid input or breach of contract.
func (it *datetime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, _ := time.Parse(time.RFC3339, s)

	*it = datetime(t)
	return nil
}

// String maps TimelineEvent into output that is compliance with functional requirements.
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

// isTransaction is true when Event is a Transaction.
func (e Event) isTransaction() bool {
	return e.Transaction != nil
}

// hasViolation is true when TimelineEvent has any violation.
func (te TimelineEvent) hasViolation() bool {
	return len(te.Violations) > 0
}
