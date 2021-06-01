package internal

import (
	"reflect"
	"testing"
	"time"
)

var now = time.Now()

type mockTimer struct {}

func (m mockTimer) Now() time.Time  {
	return now
}

// TODO: merge this two tests
func TestCreateAnAccount(t *testing.T) {
	in := Account{
		ActiveCard: true,
		AvailableLimit: 750,
	}
	want := OutputEvent{
		Event: Event{
			Account: &Account{
				ActiveCard: true,
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
	}

	timeline := NewTimeline()
	timeline.timer = mockTimer{}
	timeline.InitializeAccount(in)
	got := timeline.Events()[0]

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

// TODO: merge this two tests
func TestCreateAnAccountAlreadyInitialized(t *testing.T) {
	firstIn := Account{
		ActiveCard: true,
		AvailableLimit: 175,
	}
	secondIn := Account{
		ActiveCard: true,
		AvailableLimit: 350,
	}
	want := []OutputEvent{{
		Event: Event{
			Account: &Account{
				ActiveCard: true,
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
					ActiveCard: true,
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

	timeline := NewTimeline()
	timeline.timer = mockTimer{}
	timeline.InitializeAccount(firstIn)
	timeline.InitializeAccount(secondIn)
	got := timeline.Events()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestCreateProcessSuccessfulTransaction(t *testing.T) {
	account := Account{
		ActiveCard: true,
		AvailableLimit: 100,
	}
	trTime := time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC)
	transaction := Transaction{
		Merchant: "Burger King",
		Amount: 20,
		Time: trTime,
	}
	want := []OutputEvent{{
		Event: Event{
			Account: &Account{
				ActiveCard: true,
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
					ActiveCard: true,
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

	timeline := NewTimeline()
	timeline.timer = mockTimer{}
	timeline.InitializeAccount(account)
	timeline.ProcessTransaction(transaction)
	got := timeline.Events()

	if !reflect.DeepEqual(got[1], want[1]) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}