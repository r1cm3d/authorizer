package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestTimeline_ProcessEvent(t *testing.T) {
	cases := []struct {
		name string
		in   []InputEvent
		want []OutputEvent
	}{
		{"successful-initialization", initializeAccountInput, initializeAccountOutput},
		{"successful-transaction", successfulTransactionInput, successfulTransactionOutput},
		{"account-already-initialized", accountAlreadyInitializedInput, accountAlreadyInitializedOutput},
		{"account-not-initialized", accountNotInitializedInput, accountNotInitializedOutput},
		{"card-not-active", cardNotActiveInput, cardNotActiveOutput},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			timeline := NewTimeline()
			timeline.timer = mockTimer{}

			for _, ie := range c.in {
				timeline.ProcessEvent(ie)
			}

			if got := timeline.Events(); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

var (
	now                    = time.Now()
	initializeAccountInput = []InputEvent{{
		Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 750,
			},
			Transaction: nil,
		},
	}}
	initializeAccountOutput = []OutputEvent{{
		Event: Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 750,
			},
			Transaction: nil,
		},
		Violations: make([]Violation, 0),
	}}
	accountAlreadyInitializedInput = []InputEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 175,
				},
				Transaction: nil,
			},
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 350,
				},
				Transaction: nil,
			},
		},
	}
	accountAlreadyInitializedOutput = []OutputEvent{{
		Event: Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 175,
			},
			Transaction: nil,
		},
		Violations: make([]Violation, 0)},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 175,
				},
				Transaction: nil,
			},
			Violations: []Violation{
				accountAlreadyInitialized,
			}},
	}

	trTime                     = time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC)
	successfulTransactionInput = []InputEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 100,
				},
				Transaction: nil,
			},
		},
		{
			Event: Event{
				Account: nil,
				Transaction: &Transaction{
					Merchant: "New York Yankees",
					Amount:   20,
					Time:     trTime,
				},
			},
		},
	}
	successfulTransactionOutput = []OutputEvent{{
		Event: Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			},
			Transaction: nil,
		},
		Violations: make([]Violation, 0)},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 80,
				},
				Transaction: &Transaction{
					Merchant: "New York Yankees",
					Amount:   20,
					Time:     trTime,
				},
			},
			Violations: make([]Violation, 0)},
	}

	accountNotInitializedInput = []InputEvent{
		{
			Event: Event{
				Account: nil,
				Transaction: &Transaction{
					Merchant: "San Francisco Giants",
					Amount:   36,
					Time:     trTime,
				},
			},
		},
		{
			Event: Event{
				Account: nil,
				Transaction: &Transaction{
					Merchant: "Tampa Bay Rays",
					Amount:   20,
					Time:     trTime,
				},
			},
		},
		{
			Event: Event{
				Account: nil,
				Transaction: &Transaction{
					Merchant: "San Diego Padres",
					Amount:   15,
					Time:     trTime,
				},
			},
		},
	}
	accountNotInitializedOutput = []OutputEvent{
		{
			Event: Event{
				Account: nil,
				Transaction: &Transaction{
					Merchant: "San Francisco Giants",
					Amount:   36,
					Time:     trTime,
				},
			},
			Violations: []Violation{
				accountNotInitialized,
				cardNotActive,
			},
		},
		{
			Event: Event{
				Account: nil,
				Transaction: &Transaction{
					Merchant: "Tampa Bay Rays",
					Amount:   20,
					Time:     trTime,
				},
			},
			Violations: []Violation{
				accountNotInitialized,
				cardNotActive,
			},
		},
		{
			Event: Event{
				Account: nil,
				Transaction: &Transaction{
					Merchant: "San Diego Padres",
					Amount:   15,
					Time:     trTime,
				},
			},
			Violations: []Violation{
				accountNotInitialized,
				cardNotActive,
			},
		},
	}

	cardNotActiveInput = []InputEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     false,
					AvailableLimit: 100,
				},
				Transaction: nil,
			},
		},
		{
			Event: Event{
				Account: nil,
				Transaction: &Transaction{
					Merchant: "New York Yankees",
					Amount:   20,
					Time:     trTime,
				},
			},
		},
	}
	cardNotActiveOutput = []OutputEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     false,
					AvailableLimit: 100,
				},
				Transaction: nil,
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     false,
					AvailableLimit: 100,
				},
				Transaction: &Transaction{
					Merchant: "New York Yankees",
					Amount:   20,
					Time:     trTime,
				},
			},
			Violations: []Violation{
				cardNotActive,
			},
		},
	}
)

type mockTimer struct{}

func (m mockTimer) Now() time.Time {
	return now
}
