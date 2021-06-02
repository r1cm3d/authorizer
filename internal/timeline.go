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
	}

	t.events = append(t.events, OutputEvent{
		Event: Event{
			Account:     &newAccountState,
			Transaction: nil,
		},
		Violations: violations,
	})
}

func (t *Timeline) ProcessTransaction(tr Transaction) {
	lastValidAccountState := t.lastInitializedAccount()
	availableLimit := 0
	if lastValidAccountState != nil {
		availableLimit = lastValidAccountState.AvailableLimit
	}
	violations := t.checkTransactionViolations(tr, availableLimit)

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

func (t Timeline) checkTransactionViolations(tr Transaction, availableLimit int) []Violation {
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

	if lastInitializedAccount := t.lastInitializedAccount(); lastInitializedAccount == nil {
		return append(violations, accountNotInitialized)
	}

	if lastCardActive := t.lastAccountWithActiveCard(); lastCardActive == nil {
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


func (t Timeline) lastInitializedAccount() *Account {
	return t.lastAccountByPredicate(func(events []OutputEvent, i int) bool {
		return events[i].Account != nil
	})
}

func (t Timeline) lastAccountWithActiveCard() *Account {
	return t.lastAccountByPredicate(func(events []OutputEvent, i int) bool {
		return events[i].Account != nil && events[i].ActiveCard
	})
}

func (t Timeline) lastAccountByPredicate(pred func(events []OutputEvent, i int) bool) *Account {
	if len(t.events) <= 0 {
		return nil
	}
	sortedEvents := make([]OutputEvent, len(t.events))
	copy(sortedEvents, t.events)
	sort.Slice(sortedEvents, func(i, j int) bool { return i > j	})
	noViolPred := func(j int) bool {
		return pred(sortedEvents, j) && len(sortedEvents[j].Violations) == 0
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

func (e Event) isInitializationEvent() bool {
	return e.Account != nil
}
