package flow

// Instrument is market instrument
type Instrument int

const (
	// GbpUsd is a currency pair for GBP and USD
	GbpUsd Instrument = iota
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
	// EurGbp is a currency pair for EUR and USD
	EurGbp
)

func (instrument Instrument) String() string {
	switch instrument {
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
		return "Unknown"
	}
}
