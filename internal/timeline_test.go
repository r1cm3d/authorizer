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
		want []TimelineEvent
	}{
		{"successful-initialization", iaInput, iaOutput},
		{"successful-Transaction", sfInput, sfOutput},
		{"Account-already-initialized", aaiInput, aaiOutput},
		{"Account-not-initialized", aniInput, aniOutput},
		{"card-not-active", cnaInput, cnaOutput},
		{"insufficient-limit", ilInput, ilOutput},
		{"high-frequency-small-interval", hfInput, hfOutput},
		{"double-Transaction", dtInput, dtOutput},
		{"successful-transactions-after-hf-violation", stavInput, stavOutput},
		{"successful-transactions-after-dt-violation", stadtvInput, stadtvOutput},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			timeline := NewTimeline()

			for _, ie := range c.in {
				timeline.Process(ie)
			}

			if got := timeline.Events(); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestTimeline_Last(t *testing.T) {
	cases := []struct {
		name string
		in   []TimelineEvent
		want *TimelineEvent
	}{
		{"with nil timeline", nil, nil},
		{"without any Event", make([]TimelineEvent, 0), nil},
		{"with one Event", tlw1Event, &tlFirstEvent},
		{"with more than one", tlw2Events, &tlLastEvent},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tl := Timeline{
				events: c.in,
			}
			if got := tl.Last(); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %s, got: %s", c.name, c.want, got)
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
	iaOutput = []TimelineEvent{{
		Event: Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 750,
			},
			Transaction: nil,
		},
		Violations: make([]violation, 0),
	}}

	aaiInput = []Event{
		{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Boston Bruins",
				Amount:   666,
				Time:     datetime(now),
			},
		},
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
	aaiOutput = []TimelineEvent{{
		Event: Event{
			Account: nil,
			Transaction: &Transaction{
				Merchant: "Boston Bruins",
				Amount:   666,
				Time:     datetime(now),
			},
		},
		Violations: []violation{
			accountNotInitialized,
		}},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 175,
				},
				Transaction: nil,
			},
			Violations: make([]violation, 0),
		},
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 175,
				},
				Transaction: nil,
			},
			Violations: []violation{
				accountAlreadyInitialized,
			}},
	}

	trTime  = datetime(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC))
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
	sfOutput = []TimelineEvent{{
		Event: Event{
			Account: &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			},
			Transaction: nil,
		},
		Violations: make([]violation, 0)},
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
			Violations: make([]violation, 0)},
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
	aniOutput = []TimelineEvent{
		{
			Event: Event{
				Account: nil,
				Transaction: &Transaction{
					Merchant: "San Francisco Giants",
					Amount:   36,
					Time:     trTime,
				},
			},
			Violations: []violation{
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
			Violations: []violation{
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
			Violations: []violation{
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
	cnaOutput = []TimelineEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     false,
					AvailableLimit: 100,
				},
				Transaction: nil,
			},
			Violations: make([]violation, 0),
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
			Violations: []violation{
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
	ilOutput = []TimelineEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 100,
				},
				Transaction: nil,
			},
			Violations: make([]violation, 0),
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
			Violations: []violation{
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
			Violations: make([]violation, 0),
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
			Violations: []violation{
				insufficientLimit,
			},
		},
	}

	hfTime  = datetime(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC))
	hfTime2 = datetime(time.Date(2019, time.February, 13, 11, 1, 0, 0, time.UTC))
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
	hfOutput = []TimelineEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 100,
				},
				Transaction: nil,
			},
			Violations: make([]violation, 0),
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
			Violations: make([]violation, 0),
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
			Violations: make([]violation, 0),
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
			Violations: make([]violation, 0),
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
			Violations: []violation{
				highFrequency,
			},
		},
	}

	dtTime  = datetime(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC))
	dtTime2 = datetime(time.Date(2019, time.February, 13, 11, 1, 0, 0, time.UTC))
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
	dtOutput = []TimelineEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 100,
				},
				Transaction: nil,
			},
			Violations: make([]violation, 0),
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
			Violations: make([]violation, 0),
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
			Violations: make([]violation, 0),
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
			Violations: []violation{
				doubleTransaction,
			},
		},
	}

	stavTime  = datetime(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC))
	stavTime2 = datetime(time.Date(2019, time.February, 13, 11, 0, 1, 0, time.UTC))
	stavTime3 = datetime(time.Date(2019, time.February, 13, 11, 1, 1, 0, time.UTC))
	stavTime4 = datetime(time.Date(2019, time.February, 13, 11, 1, 31, 0, time.UTC))
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
	stavOutput = []TimelineEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 1000,
				},
				Transaction: nil,
			},
			Violations: make([]violation, 0),
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
			Violations: []violation{
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
			Violations: []violation{
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
			Violations: make([]violation, 0),
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
			Violations: make([]violation, 0),
		},
	}

	stadtvTime  = datetime(time.Date(2019, time.February, 13, 11, 0, 0, 0, time.UTC))
	stadtvTime2 = datetime(time.Date(2019, time.February, 13, 11, 0, 1, 0, time.UTC))
	stadtvTime3 = datetime(time.Date(2019, time.February, 13, 11, 1, 1, 0, time.UTC))
	stadtvTime4 = datetime(time.Date(2019, time.February, 13, 11, 1, 31, 0, time.UTC))
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
	stadtvOutput = []TimelineEvent{
		{
			Event: Event{
				Account: &Account{
					ActiveCard:     true,
					AvailableLimit: 1000,
				},
				Transaction: nil,
			},
			Violations: make([]violation, 0),
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
			Violations: []violation{
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
			Violations: []violation{
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
			Violations: make([]violation, 0),
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
			Violations: make([]violation, 0),
		},
	}

	tlFirstEvent = TimelineEvent{
		Event: Event{
			Account:     &Account{
				ActiveCard:     true,
				AvailableLimit: 666,
			},
			Transaction: nil,
		},
		Violations: make([]violation, 0),
	}
	tlLastEvent = TimelineEvent{
		Event: Event{
			Account:     &Account{
				ActiveCard:     true,
				AvailableLimit: 555,
			},
			Transaction: &Transaction{
				Merchant: "New York Rangers",
				Amount:   111,
				Time:     datetime(time.Now()),
			},
		},
		Violations: make([]violation, 0),
	}
	tlw1Event = []TimelineEvent{
		tlFirstEvent,
	}
	tlw2Events = []TimelineEvent{
		tlFirstEvent,
		tlLastEvent,
	}
)