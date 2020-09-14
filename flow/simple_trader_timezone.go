package flow

import (
	"time"
	"yukimaterrace/andaman/broker"
)

// CreateWeekdayTradableTimeZone is a factory method to create tradable time zone for monday to friday
func CreateWeekdayTradableTimeZone() *TradableTimeZone {
	return &TradableTimeZone{
		Name: "Weekday",

		OK: func(timeExtractor broker.TimeExtractor) bool {
			t := time.Unix(timeExtractor.Time(), 0)

			cond1 := t.Weekday() == time.Monday && t.Hour() >= 7
			cond2 := time.Thursday <= t.Weekday() && t.Weekday() <= time.Friday
			cond3 := t.Weekday() == time.Saturday && t.Hour() <= 5

			return cond1 || cond2 || cond3
		},
	}
}
