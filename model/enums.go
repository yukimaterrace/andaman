package model

import (
	"errors"
	"fmt"
	"time"
)

// ErrUnknownType is an error for unknown type
var ErrUnknownType = errors.New("unknown type")

// UnknownString is a string for unknown
const UnknownString = "unknown"

// MarshalJSON is a method to marshal JSON
func MarshalJSON(t fmt.Stringer) ([]byte, error) {
	s := t.String()
	if s == UnknownString {
		return nil, ErrUnknownType
	}
	return []byte(s), nil
}

// TradeSetType is trade set type enums
type TradeSetType int

const (
	// Trade is a trade type
	Trade TradeSetType = iota + 1
	// Simulation is simulation type
	Simulation
	// GridSearch is a grid search type
	GridSearch
)

const (
	tradeString      = "trade"
	simulationString = "simulation"
	gridSearchString = "grid_search"
)

func (t *TradeSetType) String() string {
	switch *t {
	case Trade:
		return tradeString
	case Simulation:
		return simulationString
	case GridSearch:
		return gridSearchString
	default:
		return UnknownString
	}
}

// UnmarshalParam is a method to unmarshal param for trade set type
func (t *TradeSetType) UnmarshalParam(param string) error {
	switch param {
	case tradeString:
		*t = Trade
	case simulationString:
		*t = Simulation
	case gridSearchString:
		*t = GridSearch
	default:
		return ErrUnknownType
	}
	return nil
}

// MarshalJSON is a method to marshal JSON for trade set type
func (t *TradeSetType) MarshalJSON() ([]byte, error) {
	return MarshalJSON(t)
}

// UnmarshalJSON is a method to unmarshal JSON for trade set type
func (t *TradeSetType) UnmarshalJSON(b []byte) error {
	return t.UnmarshalParam(string(b))
}

// TradeRunType is trade run type enums
type TradeRunType int

const (
	// OandaSimulation is oanda simulation type
	OandaSimulation TradeRunType = iota + 1
	// OandaTrade is oanda trade type
	OandaTrade
)

const (
	oandaSimulationString = "oanda_simulation"
	oandaTradeString      = "oanda_trade"
)

func (t *TradeRunType) String() string {
	switch *t {
	case OandaSimulation:
		return oandaSimulationString
	case OandaTrade:
		return oandaTradeString
	default:
		return UnknownString
	}
}

// UnmarshalParam is a method to unmarshal param
func (t *TradeRunType) UnmarshalParam(param string) error {
	switch param {
	case oandaSimulationString:
		*t = OandaSimulation
	case oandaTradeString:
		*t = OandaTrade
	default:
		return ErrUnknownType
	}
	return nil
}

// MarhsalJSON is a method to marshal JSON for trade run type
func (t *TradeRunType) MarhsalJSON() ([]byte, error) {
	return MarshalJSON(t)
}

// UnmarshalJSON is a method to unmarshal JSON for trade run type
func (t *TradeRunType) UnmarshalJSON(b []byte) error {
	return t.UnmarshalParam(string(b))
}

// TradeRunState is trade set state enums
type TradeRunState int

const (
	// Pending is a pending state
	Pending TradeRunState = iota + 1
	// Running is a running state
	Running
	// Finished is a finished state
	Finished
)

const (
	pendingString  = "pending"
	runningString  = "running"
	finishedString = "finished"
)

func (t *TradeRunState) String() string {
	switch *t {
	case Pending:
		return pendingString
	case Running:
		return runningString
	case Finished:
		return finishedString
	default:
		return UnknownString
	}
}

// UnmarshalParam is a method to unmarshal param for trade run state
func (t *TradeRunState) UnmarshalParam(param string) error {
	switch param {
	case pendingString:
		*t = Pending
	case runningString:
		*t = Running
	case finishedString:
		*t = Finished
	default:
		return ErrUnknownType
	}
	return nil
}

// MarhsalJSON is a method to marshal JSON for trade run state
func (t *TradeRunState) MarhsalJSON() ([]byte, error) {
	return MarshalJSON(t)
}

// UnmarshalJSON is a method to unmarshal JSON for trade run state
func (t *TradeRunState) UnmarshalJSON(b []byte) error {
	return t.UnmarshalParam(string(b))
}

// TradePair is trade pair enums
type TradePair int

const (
	// GbpUsd is a currency pair for GBP and USD
	GbpUsd TradePair = iota + 1
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

const (
	gbpUsdString = "GBP/USD"
	eurUsdString = "EUR/USD"
	audUsdString = "AUD/USD"
	audJpyString = "AUD/JPY"
	gbpAudString = "GBP/AUD"
	eurAudString = "EUR/AUD"
	usdJpyString = "USD/JPY"
	gbpJpyString = "GBP/JPY"
	eurJpyString = "EUR/JPY"
	eurGbpString = "EUR/GBP"
)

func (t *TradePair) String() string {
	switch *t {
	case GbpUsd:
		return gbpUsdString
	case EurUsd:
		return eurUsdString
	case AudUsd:
		return audUsdString
	case AudJpy:
		return audJpyString
	case GbpAud:
		return gbpAudString
	case EurAud:
		return eurAudString
	case UsdJpy:
		return usdJpyString
	case GbpJpy:
		return gbpJpyString
	case EurJpy:
		return eurJpyString
	case EurGbp:
		return eurGbpString
	default:
		return UnknownString
	}
}

// UnmarshalParam is a method to unmarshal param for trade pair
func (t *TradePair) UnmarshalParam(param string) error {
	switch param {
	case gbpUsdString:
		*t = GbpUsd
	case eurUsdString:
		*t = EurUsd
	case audUsdString:
		*t = AudUsd
	case audJpyString:
		*t = AudJpy
	case gbpAudString:
		*t = GbpAud
	case eurAudString:
		*t = EurAud
	case usdJpyString:
		*t = UsdJpy
	case gbpJpyString:
		*t = GbpJpy
	case eurJpyString:
		*t = EurJpy
	case eurGbpString:
		*t = EurGbp
	default:
		return ErrUnknownType
	}
	return nil
}

// MarshalJSON is a method to marshal JSON for trade pair
func (t *TradePair) MarshalJSON() ([]byte, error) {
	return MarshalJSON(t)
}

// UnmarshalJSON is a method to unmarshal JSON for trade pair
func (t *TradePair) UnmarshalJSON(b []byte) error {
	return t.UnmarshalParam(string(b))
}

// OandaInstrument is a method to get OandaInstrument
func (t TradePair) OandaInstrument() OandaInstrument {
	switch t {
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
		return UnknownString
	}
}

// PricePerPip returns price per pips
func (t TradePair) PricePerPip() float64 {
	switch t {
	case GbpUsd, EurUsd, AudUsd, GbpAud, EurAud, EurGbp:
		return 0.0001
	case AudJpy, UsdJpy, GbpJpy, EurJpy:
		return 0.01
	default:
		panic(fmt.Sprintf("unknown tradepair: %d", t))
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
	TokyoAM Timezone = iota + 1
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

const (
	tokyoAMString   = "tokyo_am"
	tokyoPMString   = "tokyo_pm"
	londonAMString  = "london_am"
	londonPMString  = "london_pm"
	newyorkAMString = "newyork_am"
	newyorkPMString = "newyork_pm"
)

func (t *Timezone) String() string {
	switch *t {
	case TokyoAM:
		return tokyoAMString
	case TokyoPM:
		return tokyoPMString
	case LondonAM:
		return londonAMString
	case LondonPM:
		return londonPMString
	case NewYorkAM:
		return newyorkAMString
	case NewYorkPM:
		return newyorkPMString
	default:
		return UnknownString
	}
}

// UnmarshalParam is a method to unmarshal param for timezone
func (t *Timezone) UnmarshalParam(param string) error {
	switch param {
	case tokyoAMString:
		*t = TokyoAM
	case tokyoPMString:
		*t = TokyoPM
	case londonAMString:
		*t = LondonAM
	case londonPMString:
		*t = LondonPM
	case newyorkAMString:
		*t = NewYorkAM
	case newyorkPMString:
		*t = NewYorkPM
	default:
		return ErrUnknownType
	}
	return nil
}

// MarshalJSON is a method to marhsal JSON for timezone
func (t *Timezone) MarshalJSON() ([]byte, error) {
	return MarshalJSON(t)
}

// UnmarshalJSON is a method to unmarshal JSON for timezone
func (t *Timezone) UnmarshalJSON(b []byte) error {
	return t.UnmarshalParam(string(b))
}

// OK is a method to valudate unix time for the timezone
func (t Timezone) OK(unix int64) bool {
	tm := time.Unix(unix, 0)

	switch t {
	case TokyoAM:
		cond1 := 7 <= tm.Hour() && tm.Hour() < 12
		cond2 := time.Monday <= tm.Weekday() && tm.Weekday() <= time.Friday
		return cond1 && cond2

	case TokyoPM:
		cond1 := 12 <= tm.Hour() && tm.Hour() < 15
		cond2 := time.Monday <= tm.Weekday() && tm.Weekday() <= time.Friday
		return cond1 && cond2

	case LondonAM:
		cond1 := 15 <= tm.Hour() && tm.Hour() < 20
		cond2 := time.Monday <= tm.Weekday() && tm.Weekday() <= time.Friday
		return cond1 && cond2

	case LondonPM:
		cond1 := 20 <= tm.Hour() && tm.Hour() < 22
		cond2 := time.Monday <= tm.Weekday() && tm.Weekday() <= time.Friday
		return cond1 && cond2

	case NewYorkAM:
		cond1 := (tm.Weekday() != time.Saturday && 22 <= tm.Hour()) || (tm.Weekday() != time.Monday && tm.Hour() < 3)
		cond2 := time.Monday <= tm.Weekday() && tm.Weekday() <= time.Saturday
		return cond1 && cond2

	case NewYorkPM:
		cond1 := 3 <= tm.Hour() && tm.Hour() < 7
		cond2 := time.Tuesday <= tm.Weekday() && tm.Weekday() <= time.Saturday
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
	Frame TradeAlgorithmType = iota + 1
)

const (
	frameString = "frame"
)

func (t *TradeAlgorithmType) String() string {
	switch *t {
	case Frame:
		return frameString
	default:
		return UnknownString
	}
}

// UnmarshalParam is a method to unmarshal param for trade algorithm type
func (t *TradeAlgorithmType) UnmarshalParam(param string) error {
	switch param {
	case frameString:
		*t = Frame
	default:
		return ErrUnknownType
	}
	return nil
}

// MarshalJSON is a method to marhsal JSON for trade algorithm type
func (t *TradeAlgorithmType) MarshalJSON() ([]byte, error) {
	return MarshalJSON(t)
}

// UnmarshalJSON is a method to unmarshal JSON for trade algorithm type
func (t *TradeAlgorithmType) UnmarshalJSON(b []byte) error {
	return t.UnmarshalParam(string(b))
}

// TradeDirection is trade direction enums
type TradeDirection int

const (
	// Long is long trade direction
	Long TradeDirection = iota + 1
	// Short is short trade direction
	Short
)

const (
	longString  = "long"
	shortString = "short"
)

func (t *TradeDirection) String() string {
	switch *t {
	case Long:
		return longString
	case Short:
		return shortString
	default:
		return UnknownString
	}
}

// UnmarshalParam is a method to unmarshal param for trade direction
func (t *TradeDirection) UnmarshalParam(param string) error {
	switch param {
	case longString:
		*t = Long
	case shortString:
		*t = Short
	default:
		return ErrUnknownType
	}
	return nil
}

// MarshalJSON is a method to marshal JSON for trade direction
func (t *TradeDirection) MarshalJSON() ([]byte, error) {
	return MarshalJSON(t)
}

// UnmarshalJSON is a method to unmarshal JSON for trade direction
func (t *TradeDirection) UnmarshalJSON(b []byte) error {
	return t.UnmarshalParam(string(b))
}

// OrderState is order state enums
type OrderState int

const (
	// Open is open state
	Open OrderState = iota + 1
	// Closed is closed state
	Closed
)

const (
	openString   = "open"
	closedString = "closed"
)

func (t *OrderState) String() string {
	switch *t {
	case Open:
		return openString
	case Closed:
		return closedString
	default:
		return UnknownString
	}
}

// UnmarshalParam is a method to unmarshal param for order state
func (t *OrderState) UnmarshalParam(param string) error {
	switch param {
	case openString:
		*t = Open
	case closedString:
		*t = Closed
	default:
		return ErrUnknownType
	}
	return nil
}

// MarshalJSON is a method to marshal JSON for order state
func (t *OrderState) MarshalJSON() ([]byte, error) {
	return MarshalJSON(t)
}

// UnmarshalJSON is a method to unmarshal JSON for order state
func (t *OrderState) UnmarshalJSON(b []byte) error {
	return t.UnmarshalParam(string(b))
}
