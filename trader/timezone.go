package trader

import (
	"time"
	"yukimaterrace/andaman/broker"
)

// CreateWeekdayTimeZone is a factory method to create tradable time zone for monday to friday
func CreateWeekdayTimeZone() *TradableTimeZone {
	return &TradableTimeZone{
		Name: "Weekday",

		OK: func(timeExtractor broker.TimeExtractor) bool {
			t := time.Unix(timeExtractor.Time(), 0)

			cond1 := t.Weekday() == time.Monday && t.Hour() >= 7
			cond2 := time.Tuesday <= t.Weekday() && t.Weekday() <= time.Friday
			cond3 := t.Weekday() == time.Saturday && t.Hour() <= 5

			return cond1 || cond2 || cond3
		},
	}
}

// CreateTokyoAMTimeZone is a factory method to create Tokyo AM time zone
func CreateTokyoAMTimeZone() *TradableTimeZone {
	return &TradableTimeZone{
		Name: "TokyoAM",

		OK: func(timeExtractor broker.TimeExtractor) bool {
			t := time.Unix(timeExtractor.Time(), 0)

			cond1 := 7 <= t.Hour() && t.Hour() < 12
			cond2 := time.Monday <= t.Weekday() && t.Weekday() <= time.Friday

			return cond1 && cond2
		},
	}
}

// CreateTokyoPMTimeZone is a factory method to create Tokyo PM time zone
func CreateTokyoPMTimeZone() *TradableTimeZone {
	return &TradableTimeZone{
		Name: "TokyoPM",

		OK: func(timeExtractor broker.TimeExtractor) bool {
			t := time.Unix(timeExtractor.Time(), 0)

			cond1 := 12 <= t.Hour() && t.Hour() < 15
			cond2 := time.Monday <= t.Weekday() && t.Weekday() <= time.Friday

			return cond1 && cond2
		},
	}
}

// CreateLondonAMTimeZone is a factory method to create London AM time zone
func CreateLondonAMTimeZone() *TradableTimeZone {
	return &TradableTimeZone{
		Name: "LondonAM",

		OK: func(timeExtractor broker.TimeExtractor) bool {
			t := time.Unix(timeExtractor.Time(), 0)

			cond1 := 15 <= t.Hour() && t.Hour() < 20
			cond2 := time.Monday <= t.Weekday() && t.Weekday() <= time.Friday

			return cond1 && cond2
		},
	}
}

// CreateLondonPMTimeZone is a factory method to create London PM time zone
func CreateLondonPMTimeZone() *TradableTimeZone {
	return &TradableTimeZone{
		Name: "LondonPM",

		OK: func(timeExtractor broker.TimeExtractor) bool {
			t := time.Unix(timeExtractor.Time(), 0)

			cond1 := 20 <= t.Hour() && t.Hour() < 22
			cond2 := time.Monday <= t.Weekday() && t.Weekday() <= time.Friday

			return cond1 && cond2
		},
	}
}

// CreateNewYorkAMTimeZone is a factory method to create NewYork AM time zone
func CreateNewYorkAMTimeZone() *TradableTimeZone {
	return &TradableTimeZone{
		Name: "NewYorkAM",

		OK: func(timeExtractor broker.TimeExtractor) bool {
			t := time.Unix(timeExtractor.Time(), 0)

			cond1 := (t.Weekday() != time.Saturday && 22 <= t.Hour()) || (t.Weekday() != time.Monday && t.Hour() < 3)
			cond2 := time.Monday <= t.Weekday() && t.Weekday() <= time.Saturday

			return cond1 && cond2
		},
	}
}

// CreateNewYorkPMTimeZone is a factory method to create NewYork PM time zone
func CreateNewYorkPMTimeZone() *TradableTimeZone {
	return &TradableTimeZone{
		Name: "NewYorkPM",

		OK: func(timeExtractor broker.TimeExtractor) bool {
			t := time.Unix(timeExtractor.Time(), 0)

			cond1 := 3 <= t.Hour() && t.Hour() < 7
			cond2 := time.Tuesday <= t.Weekday() && t.Weekday() <= time.Saturday

			return cond1 && cond2
		},
	}
}
