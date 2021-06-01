package internal

import "time"

const (
	accountAlreadyInitialized = Violation("account-already-initialized")
)

type (
	Timer interface {
		Now() time.Time
	}
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
		*Account
		*Transaction
	}
	InputEvent struct {
		Event
	}
	OutputEvent struct {
		Event
		Violations []Violation
	}
	Timeline struct {
		events []OutputEvent
		timer Timer
	}
)

func NewTimeline() Timeline {
	return Timeline{events: make([]OutputEvent, 0)}
}

func (t *Timeline) ProcessEvent(ie InputEvent) {
	if ie.isInitializationEvent() {
		t.InitializeAccount(*ie.Account)
		return
	}

	t.ProcessTransaction(*ie.Transaction)
}

func (t *Timeline) InitializeAccount(acc Account) {
	violations := make([]Violation, 0)
	newAccountState := acc
	if len(t.events) > 0 {
		violations = append(violations, accountAlreadyInitialized)
		newAccountState = *t.events[0].Account
	} else {
		violations = append(violations, "")
	}

	t.events = append(t.events, OutputEvent{
		Event: Event{
			Account:     &newAccountState,
			Transaction: &Transaction{"ISSUER", newAccountState.AvailableLimit, t.timer.Now()},
		},
		Violations: violations,
	})
}

func (t *Timeline) ProcessTransaction(tr Transaction) {
	violations := []Violation{""}
	currentAccountState := t.events[len(t.events)-1].Account
	newAccountState := Account{
		ActiveCard:     true,
		AvailableLimit: currentAccountState.AvailableLimit - tr.Amount,
	}

	t.events = append(t.events, OutputEvent{
		Event: Event{
			Account:     &newAccountState,
			Transaction: &tr,
		},
		Violations: violations,
	})
}

func (t Timeline) Events() []OutputEvent {
	return t.events
}

func (e Event) isInitializationEvent() bool {
	return e.Account != nil
}

