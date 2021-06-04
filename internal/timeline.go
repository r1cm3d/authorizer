package internal

import (
	"sort"
	"time"
)

const (
	accountAlreadyInitialized = violation("Account-already-initialized")
	accountNotInitialized     = violation("Account-not-initialized")
	cardNotActive             = violation("card-not-active")
	insufficientLimit         = violation("insufficient-limit")
	highFrequency             = violation("high-frequency-small-interval")
	doubleTransaction         = violation("double-Transaction")
)

type (
	// Timeline is the kernel of this application. It has all events present into Timeline.
	// It is NOT thread safe and SHOULD NOT be used in concurrent environments.
	Timeline struct {
		// events has the internal state of the Timeline.
		// All TimelineEvent are keep in this slice.
		// This property is not thread safe, does not have any synchronization and SHOULD NOT
		// be used in concurrent environments.
		events []TimelineEvent
	}
)

// NewTimeline creates a new Timeline.
func NewTimeline() Timeline {
	return Timeline{events: make([]TimelineEvent, 0)}
}

// Events returns all TimelineEvent stored in Timeline.
func (t Timeline) Events() []TimelineEvent {
	return t.events
}

// Last returns the last TimelineEvent.
// This method does not have neither lock strategy nor any kind of synchronization
// and SHOULD NOT be used in concurrent environments.
func (t Timeline) Last() *TimelineEvent {
	if t.events == nil || len(t.events) <= 0 {
		return nil
	}

	return &t.events[len(t.events)-1]
}

//TODO:
// - Improve README
// - Take a look at documentation one more time and find any overlooked

// Process adds an Event into Timeline. It could be either an initialization Event or a Transaction Event.
func (t *Timeline) Process(ie Event) {
	if !ie.isTransaction() {
		t.init(*ie.Account)
		return
	}

	t.add(*ie.Transaction)
}

// init handles initialization Event. Those Event should not have Transaction, only Account.
// If an initialization was done before, it will put it into TimelineEvent with an accountAlreadyInitialized violation
// plus the last valid Account state.
func (t *Timeline) init(acc Account) {
	violations := make([]violation, 0)

	newState := acc
	if initAcc := t.state(); initAcc != nil {
		violations = append(violations, accountAlreadyInitialized)
		newState = *initAcc
	}

	t.events = append(t.events, TimelineEvent{
		Event: Event{
			Account:     &newState,
			Transaction: nil,
		},
		Violations: violations,
	})
}

// add handles Transaction Event.
// It performs a series of validations before put it into TimelineEvent.
// See README.md for more details.
func (t *Timeline) add(tr Transaction) {
	lastState := t.state()
	availableLimit := 0
	if lastState != nil {
		availableLimit = lastState.AvailableLimit
	}
	violations := t.validate(tr, availableLimit)

	if len(violations) > 0 {
		oe := TimelineEvent{
			Event: Event{
				Account:     lastState,
				Transaction: &tr,
			},
			Violations: violations,
		}
		t.events = append(t.events, oe)
		return
	}

	newState := Account{
		ActiveCard:     true,
		AvailableLimit: availableLimit - tr.Amount,
	}
	oe := TimelineEvent{
		Event: Event{
			Account:     &newState,
			Transaction: &tr,
		},
		Violations: violations,
	}
	t.events = append(t.events, oe)
}

// validate performs a series of validations in the Transaction Event.
// See README.md for more details.
func (t Timeline) validate(tr Transaction, availableLimit int) []violation {
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
	violations := make([]violation, 0)

	if lastInitializedAccount := t.state(); lastInitializedAccount == nil {
		return append(violations, accountNotInitialized)
	}

	if lastCardActive := t.activeState(); lastCardActive == nil {
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

// count returns how many valid Transaction are inside the Timeline according the given function filter.
func (t Timeline) count(filter func(event Event) bool) (count int) {
	for _, outputEvent := range t.events {
		if outputEvent.isTransaction() && !outputEvent.hasViolation() && filter(outputEvent.Event) {
			count++
		}
	}

	return
}

// state returns the current Account state. It could be either active or inactive.
func (t Timeline) state() *Account {
	return t.stateByFilter(func(events []TimelineEvent, i int) bool {
		return events[i].Account != nil
	})
}


// activeState returns the last active state.
func (t Timeline) activeState() *Account {
	return t.stateByFilter(func(te []TimelineEvent, i int) bool {
		return te[i].Account != nil && te[i].ActiveCard
	})
}

// stateByFilter returns a state according a given function filter.
// It makes a copy of the current timeline, sort its in descending order and returns the first match.
// It returns nil if no state is found.
func (t Timeline) stateByFilter(filter func(te []TimelineEvent, i int) bool) *Account {
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
