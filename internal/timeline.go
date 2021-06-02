package internal

import (
	"sort"
	"time"
)

const (
	accountAlreadyInitialized = Violation("account-already-initialized")
	accountNotInitialized = Violation("account-not-initialized")
	cardNotActive = Violation("card-not-active")
	insufficientLimit = Violation("insufficient-limit")
	highFrequency = Violation("high-frequency-small-interval")
	doubleTransaction = Violation("double-transaction")
)

type (
	Timer interface {
		Now() time.Time
	}
	Violation string
	Account   struct {
		ActiveCard     bool
		AvailableLimit int
	}
	Transaction struct {
		Merchant string
		Amount   int
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
		timer  Timer
	}
)

func NewTimeline() Timeline {
	return Timeline{events: make([]OutputEvent, 0)}
}

func (t Timeline) Events() []OutputEvent {
	return t.events
}

func (t *Timeline) Process(ie InputEvent) {
	if ie.isInitEvent() {
		t.init(*ie.Account)
		return
	}

	t.add(*ie.Transaction)
}

func (t *Timeline) init(acc Account) {
	violations := make([]Violation, 0)
	newAccountState := acc
	if len(t.events) > 0 {
		violations = append(violations, accountAlreadyInitialized)
		newAccountState = *t.events[0].Account
	}

	t.events = append(t.events, OutputEvent{
		Event: Event{
			Account:     &newAccountState,
			Transaction: nil,
		},
		Violations: violations,
	})
}

func (t *Timeline) add(tr Transaction) {
	lastValidAccountState := t.lastInitAcc()
	availableLimit := 0
	if lastValidAccountState != nil {
		availableLimit = lastValidAccountState.AvailableLimit
	}
	violations := t.validate(tr, availableLimit)

	if len(violations) > 0 {
		oe := OutputEvent{
			Event: Event{
				Account:     lastValidAccountState,
				Transaction: &tr,
			},
			Violations: violations,
		}
		t.events = append(t.events, oe)
		return
	}

	newAccountState := Account{
		ActiveCard:     true,
		AvailableLimit: availableLimit - tr.Amount,
	}
	oe := OutputEvent{
		Event: Event{
			Account:     &newAccountState,
			Transaction: &tr,
		},
		Violations: violations,
	}
	t.events = append(t.events, oe)
}

func (t Timeline) validate(tr Transaction, availableLimit int) []Violation {
	const maxAllowedHF = 3
	const maxAllowedDT = 1
	const minIntervalAllowed = 2
	betweenFilter := func(e Event) bool {
		diff := tr.Time.Sub(e.Time)
		return diff.Minutes() <= minIntervalAllowed
	}
	betweenFilterSameMerchant := func(e Event) bool {
		return betweenFilter(e) && e.Merchant == tr.Merchant
	}
	violations := make([]Violation, 0)

	if lastInitializedAccount := t.lastInitAcc(); lastInitializedAccount == nil {
		return append(violations, accountNotInitialized)
	}

	if lastCardActive := t.lastActiveAcc(); lastCardActive == nil {
		return append(violations, cardNotActive)
	}

	if tr.Amount > availableLimit {
		violations = append(violations, insufficientLimit)
	}

	if t.count(betweenFilter) >= maxAllowedHF {
		violations = append(violations, highFrequency)
	}

	if t.count(betweenFilterSameMerchant) >= maxAllowedDT {
		violations = append(violations, doubleTransaction)
	}

	return violations
}

func (t Timeline) count(filter func(event Event) bool) (count int){
	for _, event := range t.events {
		// TODO: These two checks could be event methods
		if event.Transaction == nil || len(event.Violations) > 0 {
			continue
		}

		if filter(event.Event) {
			count++
		}
	}

	return
}


func (t Timeline) lastInitAcc() *Account {
	return t.lastAcctByFilter(func(events []OutputEvent, i int) bool {
		return events[i].Account != nil
	})
}

func (t Timeline) lastActiveAcc() *Account {
	return t.lastAcctByFilter(func(events []OutputEvent, i int) bool {
		return events[i].Account != nil && events[i].ActiveCard
	})
}

func (t Timeline) lastAcctByFilter(filter func(events []OutputEvent, i int) bool) *Account {
	if len(t.events) <= 0 {
		return nil
	}
	sortedEvents := make([]OutputEvent, len(t.events))
	copy(sortedEvents, t.events)
	sort.Slice(sortedEvents, func(i, j int) bool { return i > j	})
	noViolPred := func(j int) bool {
		return filter(sortedEvents, j) && len(sortedEvents[j].Violations) == 0
	}

	// TODO: change this weird name noViolPred
	i := 0
	if !noViolPred(0) {
		i = sort.Search(len(sortedEvents), noViolPred)
	}

	if i == len(sortedEvents) {
		return nil
	}
	return sortedEvents[i].Account
}

func (e Event) isInitEvent() bool {
	return e.Account != nil
}
