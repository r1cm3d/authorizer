package internal

import (
	"reflect"
	"testing"
)

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

	if !reflect.DeepEqual(got[1], want[1]) {
		t.Errorf("want: %v, got: %v", want[1], got[1])
	}
}

