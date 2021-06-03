package internal

import (
	"sort"
	"time"
)

const (
	accountAlreadyInitialized = Violation("account-already-initialized")
	accountNotInitialized     = Violation("account-not-initialized")
	cardNotActive             = Violation("card-not-active")
	insufficientLimit         = Violation("insufficient-limit")
	highFrequency             = Violation("high-frequency-small-interval")
	doubleTransaction         = Violation("double-transaction")
)

type (
	Timeline struct {
		events []TimelineEvent
		timer  Timer
	}
)

func NewTimeline() Timeline {
	return Timeline{events: make([]TimelineEvent, 0)}
}

func (t Timeline) Events() []TimelineEvent {
	return t.events
}

//TODO:
// - Implement unmarshal for TimelineEvent
// - Implement LastEvent method
// - Implement integration test for the application
// - Implement acceptance tests for the application
// - Pass golinter
// - Add documentation
// - Create docker infrastructure
// - Improve README
// - Take a look at documentation one more time and find any overlooked

func (t *Timeline) Process(ie Event) {
	if !ie.isTransaction() {
		t.init(*ie.Account)
		return
	}

	t.add(*ie.Transaction)
}

func (t *Timeline) init(acc Account) {
	violations := make([]Violation, 0)

	newAccountState := acc
	if initAcc := t.lastInitAcc(); initAcc != nil {
		violations = append(violations, accountAlreadyInitialized)
		newAccountState = *initAcc
	}

	t.events = append(t.events, TimelineEvent{
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
		oe := TimelineEvent{
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
	oe := TimelineEvent{
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
		diff := time.Time(tr.Time).Sub(time.Time(e.Time))
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

func (t Timeline) count(filter func(event Event) bool) (count int) {
	for _, outputEvent := range t.events {
		if outputEvent.isTransaction() && !outputEvent.hasViolation() && filter(outputEvent.Event) {
			count++
		}
	}

	return
}

func (t Timeline) lastInitAcc() *Account {
	return t.lastAcctByFilter(func(events []TimelineEvent, i int) bool {
		return events[i].Account != nil
	})
}

func (t Timeline) lastActiveAcc() *Account {
	return t.lastAcctByFilter(func(te []TimelineEvent, i int) bool {
		return te[i].Account != nil && te[i].ActiveCard
	})
}

func (t Timeline) lastAcctByFilter(filter func(te []TimelineEvent, i int) bool) *Account {
	if len(t.events) <= 0 {
		return nil
	}
	sortedEvents := make([]TimelineEvent, len(t.events))
	copy(sortedEvents, t.events)
	sort.Slice(sortedEvents, func(i, j int) bool { return i > j })
	filterValidEvents := func(j int) bool {
		return filter(sortedEvents, j) && !sortedEvents[j].hasViolation()
	}

	i := 0
	if !filterValidEvents(0) {
		i = sort.Search(len(sortedEvents), filterValidEvents)
	}

	if i == len(sortedEvents) {
		return nil
	}
	return sortedEvents[i].Account
}
