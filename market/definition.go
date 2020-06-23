package market

type (
	// Price is a definition for price
	Price struct {
		Value    float64
		Time     int64
		Complete bool
	}

	// PriceSequence is a definition for price sequence
	PriceSequence struct {
		Instrument  Instrument
		Granularity Granularity
		Type        PriceType
		Prices      []Price
	}

	// TradePrice is a definition for trade price
	TradePrice struct {
		Bid  float64
		Ask  float64
		Time int64
	}

	// Order is a definition for order
	Order struct {
		OrderID    string
		Instrument Instrument
		OrderType  OrderType
		OpenPrice  float64
		Profit     float64
	}

	// Orders is a definition for orders
	Orders struct {
		Orders      []Orders
		TotalProfit float64
	}

	// MadeOrder is a definition for made order
	MadeOrder struct {
		OrderID    string
		Instrument Instrument
		OrderType  OrderType
		OpenPrice  float64
		Profit     float64
	}

	// ClosedOrder is a definition for closed order
	ClosedOrder struct {
		OrderID    string
		Instument  Instrument
		OrderType  OrderType
		ClosePrice float64
		Profit     float64
	}

	// AssetSummary is a definition for asset summary
	AssetSummary struct {
		unrealizedProfit float64
		balance          float64
	}

	// PriceSequenceStatus is a definition for price sequence status
	PriceSequenceStatus struct {
		PriceSequence *PriceSequence
		Err           error
	}

	// TradePriceStatus is a definition for trade price status
	TradePriceStatus struct {
		TradePrice *TradePrice
		Err        error
	}

	// OrdersStatus is a definition for orders status
	OrdersStatus struct {
		Orders *Orders
		Err    error
	}

	// AssetStatus is a definition for asset status
	AssetStatus struct {
		AssetSummary *AssetSummary
		Err          error
	}

	// MadeOrderStatus is a definition for made order status
	MadeOrderStatus struct {
		MadeOrder *MadeOrder
		Err       error
	}

	// ClosedOrderStatus is a definition for close order status
	ClosedOrderStatus struct {
		ClosedOrder ClosedOrder
		Err         error
	}
)

type (
	request interface{}

	pricesRequest struct {
		instrument  Instrument
		granularity Granularity
		count       int
		from        int64
		replyTo     chan<- *PriceSequenceStatus
	}

	latestPriceRequest struct {
		instrument Instrument
		replyTo    chan<- *TradePriceStatus
	}

	ordersRequest struct {
		instrument Instrument
		replyTo    chan<- *OrdersStatus
	}

	assetRequest struct {
		replyTo chan<- *AssetStatus
	}

	makeOrderRequest struct {
		instrument Instrument
		orderType  OrderType
		unit       float64
		replyTo    chan<- *MadeOrderStatus
	}

	closerOrderRequest struct {
		orderID string
		replyTo chan<- *ClosedOrderStatus
	}
)

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
	default:
		return "Unknown"
	}
}

// Granularity is market time granularity
type Granularity int

const (
	// S5 is 5 sec
	S5 Granularity = iota
	// S15 is 15 sec
	S15
	// M1 is 1 min
	M1
	// M5 is 5 min
	M5
	// M15 is 15 min
	M15
	// H1 is 1 hour
	H1
	// H4 is 4 hour
	H4
)

func (granularity Granularity) String() string {
	switch granularity {
	case S5:
		return "S5"
	case S15:
		return "S15"
	case M1:
		return "M1"
	case M5:
		return "M5"
	case H1:
		return "H1"
	case H4:
		return "H4"
	default:
		return "Unknown"
	}
}

// PriceType is market price type
type PriceType int

const (
	// Bid is a price type for bid
	Bid PriceType = iota
	// Ask is a price type for ask
	Ask
	// Mid is a price type for mid
	Mid
)

func (priceType PriceType) String() string {
	switch priceType {
	case Bid:
		return "B"
	case Ask:
		return "A"
	case Mid:
		return "M"
	default:
		return "Unknown"
	}
}

// OrderType is market order type
type OrderType int

const (
	// Buy is an order type for long
	Buy OrderType = iota
	// Sell is an order type for short
	Sell
)

func (orderType OrderType) String() string {
	switch orderType {
	case Buy:
		return "BUY"
	case Sell:
		return "SELL"
	default:
		return "Unkown"
	}
}
