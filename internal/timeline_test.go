package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestTimeline_Process(t *testing.T) {
	cases := []struct {
		name string
		in   []Event
		want []OutputEvent
	}{
		{"successful-initialization", iaInput, iaOutput},
		{"successful-transaction", sfInput, sfOutput},
		{"account-already-initialized", aaiInput, aaiOutput},
		{"account-not-initialized", aniInput, aniOutput},
		{"card-not-active", cnaInput, cnaOutput},
		{"insufficient-limit", ilInput, ilOutput},
		{"high-frequency-small-interval", hfInput, hfOutput},
		{"double-transaction", dtInput, dtOutput},
		{"successful-transactions-after-hf-violation", stavInput, stavOutput},
		{"successful-transactions-after-dt-violation", stadtvInput, stadtvOutput},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			timeline := NewTimeline()
			timeline.timer = mockTimer{}

			for _, ie := range c.in {
				timeline.Process(ie)
			}

			if got := timeline.Events(); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

var (
	now = time.Now()

	iaInput = []Event{{
		Account: &Account{
			ActiveCard:     true,
			AvailableLimit: 750,
		},
		Transaction: nil,
	},
	}
	iaOutput = []OutputEvent{{
		Event: Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 750,
			},
			Transaction: nil,
		},
		Violations: make([]Violation, 0),
	}}

	aaiInput = []Event{
		{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 175,
			},
			Transaction: nil,
		},
		{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 350,
			},
			Transaction: nil,
		},
	}
	aaiOutput = []OutputEvent{{
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

	trTime  = Time(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC))
	sfInput = []Event{
		{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			},
			Transaction: nil,
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "New York Yankees",
				Amount:   20,
				Time:     trTime,
			},
		},
	}
	sfOutput = []OutputEvent{{
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

	aniInput = []Event{
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "San Francisco Giants",
				Amount:   36,
				Time:     trTime,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Tampa Bay Rays",
				Amount:   20,
				Time:     trTime,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "San Diego Padres",
				Amount:   15,
				Time:     trTime,
			},
		},
	}
	aniOutput = []OutputEvent{
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
			},
		},
	}

	cnaInput = []Event{
		{
			Account: &Account{
				ActiveCard:     false,
				AvailableLimit: 100,
			},
			Transaction: nil,
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "New York Yankees",
				Amount:   20,
				Time:     trTime,
			},
		},
	}
	cnaOutput = []OutputEvent{
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

	ilInput = []Event{
		{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			},
			Transaction: nil,
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Pittsburgh Pirates",
				Amount:   101,
				Time:     trTime,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Chicago White Sox",
				Amount:   98,
				Time:     trTime,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "St. Louis Cardinals",
				Amount:   5,
				Time:     trTime,
			},
		},
	}
	ilOutput = []OutputEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 100,
				},
				Transaction: nil,
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 100,
				},
				Transaction: &Transaction{
					Merchant: "Pittsburgh Pirates",
					Amount:   101,
					Time:     trTime,
				},
			},
			Violations: []Violation{
				insufficientLimit,
			},
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 2,
				},
				Transaction: &Transaction{
					Merchant: "Chicago White Sox",
					Amount:   98,
					Time:     trTime,
				},
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 2,
				},
				Transaction: &Transaction{
					Merchant: "St. Louis Cardinals",
					Amount:   5,
					Time:     trTime,
				},
			},
			Violations: []Violation{
				insufficientLimit,
			},
		},
	}

	hfTime  = Time(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC))
	hfTime2 = Time(time.Date(2019, time.February, 13, 11, 1, 0, 0, time.UTC))
	hfInput = []Event{
		{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			},
			Transaction: nil,
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Atlanta Braves",
				Amount:   10,
				Time:     hfTime,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "New York Mets",
				Amount:   11,
				Time:     hfTime,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Los Angeles Dodgers",
				Amount:   12,
				Time:     hfTime2,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Boston Red Sox",
				Amount:   16,
				Time:     hfTime2,
			},
		},
	}
	hfOutput = []OutputEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 100,
				},
				Transaction: nil,
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 90,
				},
				Transaction: &Transaction{
					Merchant: "Atlanta Braves",
					Amount:   10,
					Time:     hfTime,
				},
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 79,
				},
				Transaction: &Transaction{
					Merchant: "New York Mets",
					Amount:   11,
					Time:     hfTime,
				},
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 67,
				},
				Transaction: &Transaction{
					Merchant: "Los Angeles Dodgers",
					Amount:   12,
					Time:     hfTime2,
				},
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 67,
				},
				Transaction: &Transaction{
					Merchant: "Boston Red Sox",
					Amount:   16,
					Time:     hfTime2,
				},
			},
			Violations: []Violation{
				highFrequency,
			},
		},
	}

	dtTime  = Time(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC))
	dtTime2 = Time(time.Date(2019, time.February, 13, 11, 1, 0, 0, time.UTC))
	dtInput = []Event{
		{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			},
			Transaction: nil,
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Toronto Blue Jays",
				Amount:   10,
				Time:     dtTime,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Los Angeles Angels",
				Amount:   11,
				Time:     dtTime,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Los Angeles Angels",
				Amount:   12,
				Time:     dtTime2,
			},
		},
	}
	dtOutput = []OutputEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 100,
				},
				Transaction: nil,
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 90,
				},
				Transaction: &Transaction{
					Merchant: "Toronto Blue Jays",
					Amount:   10,
					Time:     dtTime,
				},
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 79,
				},
				Transaction: &Transaction{
					Merchant: "Los Angeles Angels",
					Amount:   11,
					Time:     dtTime,
				},
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 79,
				},
				Transaction: &Transaction{
					Merchant: "Los Angeles Angels",
					Amount:   12,
					Time:     dtTime2,
				},
			},
			Violations: []Violation{
				doubleTransaction,
			},
		},
	}

	stavTime  = Time(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC))
	stavTime2 = Time(time.Date(2019, time.February, 13, 11, 0, 1, 0, time.UTC))
	stavTime3 = Time(time.Date(2019, time.February, 13, 11, 1, 1, 0, time.UTC))
	stavTime4 = Time(time.Date(2019, time.February, 13, 11, 1, 31, 0, time.UTC))
	stavInput = []Event{
		{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 1000,
			},
			Transaction: nil,
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Philadelphia Phillies",
				Amount:   1250,
				Time:     stavTime,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Cleveland Indians",
				Amount:   2500,
				Time:     stavTime2,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Milwaukee Brewers",
				Amount:   800,
				Time:     stavTime3,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Cincinnati Reds",
				Amount:   80,
				Time:     stavTime4,
			},
		},
	}
	stavOutput = []OutputEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 1000,
				},
				Transaction: nil,
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 1000,
				},
				Transaction: &Transaction{
					Merchant: "Philadelphia Phillies",
					Amount:   1250,
					Time:     stavTime,
				},
			},
			Violations: []Violation{
				insufficientLimit,
			},
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 1000,
				},
				Transaction: &Transaction{
					Merchant: "Cleveland Indians",
					Amount:   2500,
					Time:     stavTime2,
				},
			},
			Violations: []Violation{
				insufficientLimit,
			},
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 200,
				},
				Transaction: &Transaction{
					Merchant: "Milwaukee Brewers",
					Amount:   800,
					Time:     stavTime3,
				},
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 120,
				},
				Transaction: &Transaction{
					Merchant: "Cincinnati Reds",
					Amount:   80,
					Time:     stavTime4,
				},
			},
			Violations: make([]Violation, 0),
		},
	}

	stadtvTime  = Time(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC))
	stadtvTime2 = Time(time.Date(2019, time.February, 13, 11, 0, 1, 0, time.UTC))
	stadtvTime3 = Time(time.Date(2019, time.February, 13, 11, 1, 1, 0, time.UTC))
	stadtvTime4 = Time(time.Date(2019, time.February, 13, 11, 1, 31, 0, time.UTC))
	stadtvInput = []Event{
		{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 1000,
			},
			Transaction: nil,
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Philadelphia Phillies",
				Amount:   1250,
				Time:     stadtvTime,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Cleveland Indians",
				Amount:   2500,
				Time:     stadtvTime2,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Cleveland Indians",
				Amount:   800,
				Time:     stadtvTime3,
			},
		},
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Philadelphia Phillies",
				Amount:   80,
				Time:     stadtvTime4,
			},
		},
	}
	stadtvOutput = []OutputEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 1000,
				},
				Transaction: nil,
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 1000,
				},
				Transaction: &Transaction{
					Merchant: "Philadelphia Phillies",
					Amount:   1250,
					Time:     stavTime,
				},
			},
			Violations: []Violation{
				insufficientLimit,
			},
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 1000,
				},
				Transaction: &Transaction{
					Merchant: "Cleveland Indians",
					Amount:   2500,
					Time:     stavTime2,
				},
			},
			Violations: []Violation{
				insufficientLimit,
			},
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 200,
				},
				Transaction: &Transaction{
					Merchant: "Cleveland Indians",
					Amount:   800,
					Time:     stavTime3,
				},
			},
			Violations: make([]Violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 120,
				},
				Transaction: &Transaction{
					Merchant: "Philadelphia Phillies",
					Amount:   80,
					Time:     stavTime4,
				},
			},
			Violations: make([]Violation, 0),
		},
	}
)

type mockTimer struct{}

func (m mockTimer) Now() time.Time {
	return now
}
