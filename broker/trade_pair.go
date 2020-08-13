package broker

import "fmt"

// TradePair is market instrument
type TradePair string

const (
	// GbpUsd is a currency pair for GBP and USD
	GbpUsd TradePair = "GBP_USD"
	// EurUsd is a currency pair for EUR and USD
	EurUsd TradePair = "EUR_USD"
	// AudUsd is a currency pair for AUD and USD
	AudUsd TradePair = "AUD_USD"
	// AudJpy is a currency pair for AUD and JPY
	AudJpy TradePair = "AUD_JPY"
	// GbpAud is a currency pair for GBP and AUD
	GbpAud TradePair = "GBP_AUD"
	// EurAud is a currency pair for EUR and AUD
	EurAud TradePair = "EUR_AUD"
	// UsdJpy is a currency pair for USD and JPY
	UsdJpy TradePair = "USD_JPY"
	// GbpJpy is a currency pair for GBP and JPY
	GbpJpy TradePair = "GBP_JPY"
	// EurJpy is a currency pair for EUR and JPY
	EurJpy TradePair = "EUR_JPY"
	// EurGbp is a currency pair for EUR and GBP
	EurGbp TradePair = "EUR_GBP"
)

func priceGap2Pips(tradePair TradePair) float64 {
	switch tradePair {
	case GbpUsd, EurUsd, AudUsd, GbpAud, EurAud, EurGbp:
		return 0.0001
	case AudJpy, UsdJpy, GbpJpy, EurJpy:
		return 0.01
	default:
		panic(fmt.Sprintf("unknown tradepair: %s", string(tradePair)))
	}
}
