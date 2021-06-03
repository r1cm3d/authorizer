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
		{"Account", accJSON, accEvent},
		{"Transaction", trJSON, trEvent},
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
		{"without Account", tewoAcc, woAcc},
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
	accJSON  = `{"Account":{"active-card":true,"available-limit":666}}`
	accEvent = Event{
		Account: &Account{
			ActiveCard:     true,
			AvailableLimit: 666,
		},
		Transaction: nil,
	}

	trJSON  = `{"Transaction":{"merchant":"Montreal Canadiens","amount":666,"time":"2019-02-13T11:00:00.000Z"}}`
	trEvent = Event{
		Account: nil,
		Transaction: &Transaction{
			Merchant: "Montreal Canadiens",
			Amount:   666,
			Time:     datetime(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC)),
		},
	}

	tewoAcc = TimelineEvent{
		Event: Event{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Vegas Golden Knights",
				Amount:   142,
				Time:     datetime(time.Now()),
			},
		},
		Violations: []violation{accountNotInitialized},
	}
	woAcc = `{"Account":{},"violations":["Account-not-initialized"]}`

	tew1Vio = TimelineEvent{
		Event: Event{
			Account: &Account{
				ActiveCard:     false,
				AvailableLimit: 666,
			},
			Transaction: nil,
		},
		Violations: []violation{accountAlreadyInitialized},
	}
	w1Vio = `{"Account":{"active-card":false,"available-limit":666},"violations":["Account-already-initialized"]}`

	tew2Vio = TimelineEvent{
		Event: Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 175,
			},
			Transaction: &Transaction{
				Merchant: "Pittsburgh Penguins",
				Amount:   175,
				Time:     datetime(time.Now()),
			},
		},
		Violations: []violation{
			insufficientLimit,
			doubleTransaction,
		},
	}
	w2Vio = `{"Account":{"active-card":true,"available-limit":175},"violations":["insufficient-limit","double-Transaction"]}`

	tewoVio = TimelineEvent{
		Event: Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 666,
			},
			Transaction: nil,
		},
		Violations: make([]violation, 0),
	}
	woVio = `{"Account":{"active-card":true,"available-limit":666},"violations":[]}`
)
