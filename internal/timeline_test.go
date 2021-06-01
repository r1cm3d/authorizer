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
		kind: "creation",
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
