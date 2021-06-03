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

var (
	accJson  = `{"account": {"active-card": true, "available-limit": 666}}`
	accEvent = Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 666,
			},
			Transaction: nil,
	}

	trJson  = `{"transaction": {"merchant": "Montreal Canadiens", "amount": 666, "time": "2019-02-13T11:00:00.000Z"}}`
	trEvent = Event{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Montreal Canadiens",
				Amount:   666,
				Time: Time(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC)),
			},
	}
)
