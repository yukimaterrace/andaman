package market

type (
	// Instrument is an interface for instrument
	Instrument interface{}

	// GbpUsd is a currency pair for GBP and USD
	GbpUsd Instrument
	// EurUsd is a currency pair for EUR and USD
	EurUsd Instrument
	// AudUsd is a currency pair for AUD and USD
	AudUsd Instrument
	// AudJpy is a currency pair for AUD and JPY
	AudJpy Instrument
	// GbpAud is a currency pair for GBP and AUD
	GbpAud Instrument
	// EurAud is a currency pair for EUR and AUD
	EurAud Instrument
	// UsdJpy is a currency pair for UDS and JPY
	UsdJpy Instrument
	// GbpJpy is a currency pair for GBP and JPY
	GbpJpy Instrument
	// EurJpy is a currency pair for EUR and JPY
	EurJpy Instrument
)

type (
	// Time is an interface for market time
	Time interface{}

	// S5 is 5 sec
	S5 Time
	// S15 is 15 sec
	S15 Time
	// M1 is 1 min
	M1 Time
	// M5 is 5 min
	M5 Time
	// M15 is 15 min
	M15 Time
	// H1 is 1 hour
	H1 Time
	// H4 is 4 hour
	H4 Time
)

type (
	// PriceType is an interface for price type
	PriceType interface{}

	// Bid is a price type for bid
	Bid PriceType
	// Ask is a price type for ask
	Ask PriceType
)

type (
	// OrderType is an interface for order type
	OrderType interface{}

	// Buy is an order type for long
	Buy OrderType
	// Sell is an order type for short
	Sell OrderType
)
