package internal

import (
	"reflect"
	"testing"
	"time"
)

// TODO: merge this two tests
func TestCreateAnAccount(t *testing.T) {
	in := Account{
		ActiveCard: true,
		AvailableLimit: 750,
	}
	want := Event{
		kind: EventKind("creation"),
		Account: Account{
			ActiveCard: true,
			AvailableLimit: 750,
		},
		Violations: []Violation{
			Violation(""),
		},
	}

	timeline := NewTimeline()
	timeline.AddCreationEvent(in)
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
	want := []Event{{
		kind: EventKind("creation"),
		Account: Account{
			ActiveCard:     true,
			AvailableLimit: 175,
		},
		Violations: []Violation{
			Violation(""),
		}},
		{
			kind: EventKind("creation"),
			Account: Account{
				ActiveCard:     true,
				AvailableLimit: 175,
			},
			Violations: []Violation{
				Violation("account-already-initialized"),
			}},
	}

	timeline := NewTimeline()
	timeline.AddCreationEvent(firstIn)
	timeline.AddCreationEvent(secondIn)
	got := timeline.Events()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestCreateProcessSucessfulTransaction(t *testing.T) {
	account := Account{
		ActiveCard: true,
		AvailableLimit: 100,
	}
	transaction := Transaction{
		Merchant: "Burger King",
		Amount: 20,
		Time: time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC),
	}
	want := []Event{{
		kind: EventKind("creation"),
		Account: Account{
			ActiveCard:     true,
			AvailableLimit: 100,
		},
		Violations: []Violation{
			Violation(""),
		}},
		{
			kind: EventKind("transaction"),
			Account: Account{
				ActiveCard:     true,
				AvailableLimit: 80,
			},
			Violations: []Violation{
				Violation(""),
			}},
	}

	timeline := NewTimeline()
	timeline.AddCreationEvent(account)
	timeline.ProcessTransaction(transaction)
	got := timeline.Events()

	if !reflect.DeepEqual(got[1], want[1]) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}