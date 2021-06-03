package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Event
	}{
		{"account", accJson, accEvent},
		{"transaction", trJson, trEvent},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Parse(c.in); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestTimelineEvent_String(t *testing.T) {
	cases := []struct {
		name string
		in   TimelineEvent
		want string
	}{
		{"without account", tewoAcc, woAcc},
		{"with one violation", tew1Vio, w1Vio},
		{"with two violation", tew2Vio, w2Vio},
		{"without violation", tewoVio, woVio},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.in.String(); got != c.want {
				t.Errorf("%s, want: %s, got: %s", c.name, c.want, got)
			}
		})
	}
}

var (
	accJson  = `{"account":{"active-card":true,"available-limit":666}}`
	accEvent = Event{
		Account: &Account{
			ActiveCard:     true,
			AvailableLimit: 666,
		},
		Transaction: nil,
	}

	trJson  = `{"transaction":{"merchant":"Montreal Canadiens","amount":666,"time":"2019-02-13T11:00:00.000Z"}}`
	trEvent = Event{
		Account: nil,
		Transaction: &Transaction{
			Merchant: "Montreal Canadiens",
			Amount:   666,
			Time:     Time(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC)),
		},
	}

	tewoAcc = TimelineEvent{
		Event: Event{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Vegas Golden Knights",
				Amount:   142,
				Time:     Time(time.Now()),
			},
		},
		Violations: []Violation{accountNotInitialized},
	}
	woAcc = `{"account":{},"violations":["account-not-initialized"]}`

	tew1Vio = TimelineEvent{
		Event: Event{
			Account: &Account{
				ActiveCard:     false,
				AvailableLimit: 666,
			},
			Transaction: nil,
		},
		Violations: []Violation{accountAlreadyInitialized},
	}
	w1Vio = `{"account":{"active-card":false,"available-limit":666},"violations":["account-already-initialized"]}`

	tew2Vio = TimelineEvent{
		Event: Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 175,
			},
			Transaction: &Transaction{
				Merchant: "Pittsburgh Penguins",
				Amount:   175,
				Time:     Time(time.Now()),
			},
		},
		Violations: []Violation{
			insufficientLimit,
			doubleTransaction,
		},
	}
	w2Vio = `{"account":{"active-card":true,"available-limit":175},"violations":["insufficient-limit","double-transaction"]}`

	tewoVio = TimelineEvent{
		Event: Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 666,
			},
			Transaction: nil,
		},
		Violations: make([]Violation, 0),
	}
	woVio = `{"account":{"active-card":true,"available-limit":666},"violations":[]}`
)
