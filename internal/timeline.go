package internal

import "time"

const (
	creation = EventKind("creation")
	transaction = EventKind("transaction")

	accountAlreadyInitialized = Violation("account-already-initialized")
)

type (
	EventKind string
	Violation string
	Account   struct {
		ActiveCard bool
		AvailableLimit int
	}
	Transaction struct {
		Merchant string
		Amount int
		time.Time
	}
	Event struct {
		kind EventKind
		Account
		Violations []Violation
	}
	Timeline struct {
		events []Event
	}
)

func NewTimeline() Timeline {
	return Timeline{events: make([]Event, 0)}
}

func (t *Timeline) AddCreationEvent(acc Account) {
	violations := make([]Violation, 0)
	currentState := acc
	if len(t.events) > 0 {
		violations = append(violations, accountAlreadyInitialized)
		currentState = t.events[0].Account
	} else {
		violations = append(violations, "")
	}

	t.events = append(t.events, Event{
		kind: creation,
		Account:    currentState,
		Violations: violations,
	})
}

func (t *Timeline) ProcessTransaction(tr Transaction) {
	violations := []Violation{""}
	newAvailableLimit := t.events[len(t.events)-1].Account.AvailableLimit - tr.Amount

	t.events = append(t.events, Event{
		kind: transaction,
		Account:    Account{
			ActiveCard:     true,
			AvailableLimit: newAvailableLimit,
		},
		Violations: violations,
	})
}


func (t Timeline) Events() []Event {
	return t.events
}

