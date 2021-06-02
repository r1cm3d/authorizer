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
		{"Initialize an account", initializeAccountInput, initializeAccountOutput},
		{"Account already initialized", accountAlreadyInitializedInput, accountAlreadyInitializedOutput},
		{"Successful transaction", successfulTransactionInput, successfulTransactionOutput},
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
	now = time.Now()
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
			Transaction: &Transaction{
				Merchant: "ISSUER",
				Amount:   750,
				Time:     now,
			},
		},
		Violations: []Violation{
			Violation(""),
		},
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
			Transaction: &Transaction{
				Merchant: "ISSUER",
				Amount:   175,
				Time:     now,
			},
		},
		Violations: []Violation{
			Violation(""),
		}},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 175,
				},
				Transaction: &Transaction{
					Merchant: "ISSUER",
					Amount:   175,
					Time:     now,
				},
			},
			Violations: []Violation{
				Violation("account-already-initialized"),
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
					Merchant: "Burger King",
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
			Transaction: &Transaction{
				Merchant: "ISSUER",
				Amount:   100,
				Time:     now,
			},
		},
		Violations: []Violation{
			Violation(""),
		}},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 80,
				},
				Transaction: &Transaction{
					Merchant: "Burger King",
					Amount:   20,
					Time:     trTime,
				},
			},
			Violations: []Violation{
				Violation(""),
			}},
	}
)

type mockTimer struct{}

func (m mockTimer) Now() time.Time {
	return now
}