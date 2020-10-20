package model

import (
	"errors"
	"fmt"
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

// TradeSetState is trade set state enums
type TradeSetState int

const (
	// Running is a running state
	Running TradeSetState = iota
	// Stopped is a stopped state
	Stopped
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

// OandaString is a string representation for OANDA
func (tradePair TradePair) OandaString() string {
	switch tradePair {
	case GbpUsd:
		return "GBP_USD"
	case EurUsd:
		return "EUR_USD"
	case AudUsd:
		return "AUD_USD"
	case AudJpy:
		return "AUD_JPY"
	case GbpAud:
		return "GBP_AUD"
	case EurAud:
		return "EUR_AUD"
	case UsdJpy:
		return "USD_JPY"
	case GbpJpy:
		return "GBP_JPY"
	case EurJpy:
		return "EUR_JPY"
	case EurGbp:
		return "EUR_GBP"
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
