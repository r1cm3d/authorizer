package internal

const (
	creation = EventKind("creation")
	transaction = EventKind("transaction")
)

type (
	EventKind string
	Violation string
	Account   struct {
		ActiveCard bool
		AvailableLimit int
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
	t.events = append(t.events, Event{
		kind: creation,
		Account:    acc,
		Violations: []Violation{""},
	})
}

func (t Timeline) Events() []Event {
	return t.events
}

