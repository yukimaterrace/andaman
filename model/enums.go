package model

import (
	"errors"
	"fmt"
	"time"
)

// ErrNumber is an error for number
var ErrNumber = errors.New("invalid number")

// TradeSetType is trade set type enums
type TradeSetType int

const (
	// Trade is a trade type
	Trade TradeSetType = iota
	// Simulation is simulation type
	Simulation
	// GridSearch is a grid search type
	GridSearch
)

// IsValid is a method to validate
func (_type TradeSetType) IsValid() error {
	switch _type {
	case Trade, Simulation, GridSearch:
		return nil
	default:
		return ErrNumber
	}
}

// TradeRunType is trade run type enums
type TradeRunType int

const (
	// OandaSimulation is oanda simulation type
	OandaSimulation TradeRunType = iota
	// OandaTrade is oanda trade type
	OandaTrade
)

// IsValid is a method to validate
func (_type TradeRunType) IsValid() error {
	switch _type {
	case OandaSimulation, OandaTrade:
		return nil
	default:
		return ErrNumber
	}
}

// TradeRunState is trade set state enums
type TradeRunState int

const (
	// Pending is a pending state
	Pending TradeRunState = iota
	// Running is a running state
	Running
	// Finished is a finished state
	Finished
)

// TradePair is trade pair enums
type TradePair int

const (
	// GbpUsd is a currency pair for GBP and USD
	GbpUsd TradePair = iota
	// EurUsd is a currency pair for EUR and USD
	EurUsd
	// AudUsd is a currency pair for AUD and USD
	AudUsd
	// AudJpy is a currency pair for AUD and JPY
	AudJpy
	// GbpAud is a currency pair for GBP and AUD
	GbpAud
	// EurAud is a currency pair for EUR and AUD
	EurAud
	// UsdJpy is a currency pair for USD and JPY
	UsdJpy
	// GbpJpy is a currency pair for GBP and JPY
	GbpJpy
	// EurJpy is a currency pair for EUR and JPY
	EurJpy
	// EurGbp is a currency pair for EUR and GBP
	EurGbp
)

// OandaInstrument is a method to get OandaInstrument
func (tradePair TradePair) OandaInstrument() OandaInstrument {
	switch tradePair {
	case GbpUsd:
		return OandaGbpUsd
	case EurUsd:
		return OandaEurUsd
	case AudUsd:
		return OandaAudUsd
	case AudJpy:
		return OandaAudJpy
	case GbpAud:
		return OandaGbpAud
	case EurAud:
		return OandaEurAud
	case UsdJpy:
		return OandaUsdJpy
	case GbpJpy:
		return OandaGbpJpy
	case EurJpy:
		return OandaEurJpy
	case EurGbp:
		return OandaEurGbp
	default:
		return "unknown"
	}
}

// PricePerPip returns price per pips
func (tradePair TradePair) PricePerPip() float64 {
	switch tradePair {
	case GbpUsd, EurUsd, AudUsd, GbpAud, EurAud, EurGbp:
		return 0.0001
	case AudJpy, UsdJpy, GbpJpy, EurJpy:
		return 0.01
	default:
		panic(fmt.Sprintf("unknown tradepair: %d", tradePair))
	}
}

// TradePairIterator is a struct for trade pair iterator
type TradePairIterator struct {
	current TradePair
}

// Next is a method to know if the next item exists
func (iterator TradePairIterator) Next() bool {
	if iterator.current <= EurGbp {
		return true
	}
	return false
}

// Value is a method to get the value
func (iterator TradePairIterator) Value() TradePair {
	value := iterator.current
	iterator.current++
	return value
}

// OandaInstrument is a definition for oanda instrument
type OandaInstrument string

const (
	// OandaGbpUsd is a const for oanda GBP/USD
	OandaGbpUsd OandaInstrument = "GBP_USD"
	// OandaEurUsd is a const for oanda EUR/USD
	OandaEurUsd OandaInstrument = "EUR_USD"
	// OandaAudUsd is a const for oanda AUD/USD
	OandaAudUsd OandaInstrument = "AUD_USD"
	// OandaAudJpy is a const for oanda AUD/JPY
	OandaAudJpy OandaInstrument = "AUD_JPY"
	// OandaGbpAud is a const for oanda GBP/AUD
	OandaGbpAud OandaInstrument = "GBP_AUD"
	// OandaEurAud is a const for oanda EUR/AUD
	OandaEurAud OandaInstrument = "EUR_AUD"
	// OandaUsdJpy is a const for oanda USD/JPY
	OandaUsdJpy OandaInstrument = "USD_JPY"
	// OandaGbpJpy is a const for oanda GBP/JPY
	OandaGbpJpy OandaInstrument = "GBP_JPY"
	// OandaEurJpy is a cosnt for oanda EUR/JPY
	OandaEurJpy OandaInstrument = "EUR_JPY"
	// OandaEurGbp is a const for oanda EUR/GBP
	OandaEurGbp OandaInstrument = "EUR_GBP"
)

// TradePair is a method to get trade pair
func (instrument OandaInstrument) TradePair() TradePair {
	switch instrument {
	case OandaGbpUsd:
		return GbpUsd
	case OandaEurUsd:
		return EurUsd
	case OandaAudUsd:
		return AudUsd
	case OandaAudJpy:
		return AudJpy
	case OandaGbpAud:
		return GbpAud
	case OandaEurAud:
		return EurAud
	case OandaUsdJpy:
		return UsdJpy
	case OandaGbpJpy:
		return GbpJpy
	case OandaEurJpy:
		return EurJpy
	case OandaEurGbp:
		return EurGbp
	default:
		return -1
	}
}

// Timezone is timezone enums
type Timezone int

const (
	// TokyoAM is Tokyo AM timezone
	TokyoAM Timezone = iota
	// TokyoPM is Tokyo PM timezone
	TokyoPM
	// LondonAM is London AM timezone
	LondonAM
	// LondonPM is London PM timezone
	LondonPM
	// NewYorkAM is NewYork AM timezone
	NewYorkAM
	// NewYorkPM is NewYork PM timezone
	NewYorkPM
)

// OK is a method to valudate unix time for the timezone
func (timezone Timezone) OK(unix int64) bool {
	t := time.Unix(unix, 0)

	switch timezone {
	case TokyoAM:
		cond1 := 7 <= t.Hour() && t.Hour() < 12
		cond2 := time.Monday <= t.Weekday() && t.Weekday() <= time.Friday
		return cond1 && cond2

	case TokyoPM:
		cond1 := 12 <= t.Hour() && t.Hour() < 15
		cond2 := time.Monday <= t.Weekday() && t.Weekday() <= time.Friday
		return cond1 && cond2

	case LondonAM:
		cond1 := 15 <= t.Hour() && t.Hour() < 20
		cond2 := time.Monday <= t.Weekday() && t.Weekday() <= time.Friday
		return cond1 && cond2

	case LondonPM:
		cond1 := 20 <= t.Hour() && t.Hour() < 22
		cond2 := time.Monday <= t.Weekday() && t.Weekday() <= time.Friday
		return cond1 && cond2

	case NewYorkAM:
		cond1 := (t.Weekday() != time.Saturday && 22 <= t.Hour()) || (t.Weekday() != time.Monday && t.Hour() < 3)
		cond2 := time.Monday <= t.Weekday() && t.Weekday() <= time.Saturday
		return cond1 && cond2

	case NewYorkPM:
		cond1 := 3 <= t.Hour() && t.Hour() < 7
		cond2 := time.Tuesday <= t.Weekday() && t.Weekday() <= time.Saturday
		return cond1 && cond2

	default:
		return false
	}
}

// TimezoneIterator is a struct for timezone iterator
type TimezoneIterator struct {
	current Timezone
}

// Next is a method to know if the next item exists
func (iterator TimezoneIterator) Next() bool {
	if iterator.current <= NewYorkPM {
		return true
	}
	return false
}

// Value is a method to get the value
func (iterator TimezoneIterator) Value() Timezone {
	value := iterator.current
	iterator.current++
	return value
}

// TradeAlgorithmType is trade algorithm type enums
type TradeAlgorithmType int

const (
	// Frame is frame trade algorithm type
	Frame TradeAlgorithmType = iota
)

// TradeDirection is trade direction enums
type TradeDirection int

const (
	// Long is long trade direction
	Long TradeDirection = iota
	// Short is short trade direction
	Short
)

// OrderState is order state enums
type OrderState int

const (
	// Open is open state
	Open OrderState = iota
	// Closed is closed state
	Closed
)
