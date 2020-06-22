package market

type (
	// Price is a definition for price
	Price struct {
		Value float64
		Time  int64
	}

	// PriceSequence is a definition for price sequence
	PriceSequence struct {
		Instrument  Instrument
		Granularity Granularity
		Prices      []Price
	}

	// PriceDetail is a definition for price detail
	PriceDetail struct {
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

	// PriceDetailStatus is a definition for price detail status
	PriceDetailStatus struct {
		PriceDetail *PriceDetail
		Err         error
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
		replyTo     chan<- *PriceSequence
	}

	latestPriceRequest struct {
		instrument Instrument
		replyTo    chan<- *PriceDetail
	}

	ordersRequest struct {
		instrument Instrument
		replyTo    chan<- *Orders
	}

	assetRequest struct {
		replyTo chan<- *AssetStatus
	}

	makeOrderRequest struct {
		instrument Instrument
		orderType  OrderType
		unit       float64
		replyTo    chan<- *MakeOrderStatus
	}

	closerOrderRequest struct {
		orderID string
		replyTo chan<- *CloseOrderStatus
	}
)

type (
	// Instrument is an interface for instrument
	Instrument interface {
		String()
	}

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
	// UsdJpy is a currency pair for USD and JPY
	UsdJpy Instrument
	// GbpJpy is a currency pair for GBP and JPY
	GbpJpy Instrument
	// EurJpy is a currency pair for EUR and JPY
	EurJpy Instrument
)

func (gbpUsd GbpUsd) String() string {
	return "GBP_USD"
}

func (eurUsd EurUsd) String() string {
	return "EUR_USD"
}

func (audUsd AudUsd) String() string {
	return "AUD_USD"
}

func (audJpy AudJpy) String() string {
	return "AUD_JPY"
}

func (gbpAud GbpAud) String() string {
	return "GBP_AUD"
}

func (eurAud EurAud) String() string {
	return "EUR_AUD"
}

func (usdJpy UsdJpy) String() string {
	return "USD_JPY"
}

func (gbpJpy GbpJpy) String() string {
	return "GBP_JPY"
}

func (eurJpy EurJpy) String() string {
	return "EUR_JPY"
}

type (
	// Granularity is an interface for time granularity
	Granularity interface {
		String()
	}

	// S5 is 5 sec
	S5 Granularity
	// S15 is 15 sec
	S15 Granularity
	// M1 is 1 min
	M1 Granularity
	// M5 is 5 min
	M5 Granularity
	// M15 is 15 min
	M15 Granularity
	// H1 is 1 hour
	H1 Granularity
	// H4 is 4 hour
	H4 Granularity
)

func (s5 S5) String() string {
	return "S5"
}

func (s15 S15) String() string {
	return "S15"
}

func (m1 M1) String() string {
	return "M1"
}

func (m5 M5) String() string {
	return "M5"
}

func (h1 H1) String() string {
	return "H1"
}

func (h4 H4) String() string {
	return "H4"
}

type (
	// PriceType is an interface for price type
	PriceType interface {
		String() string
	}

	// Bid is a price type for bid
	Bid PriceType
	// Ask is a price type for ask
	Ask PriceType
	// Mid is a price type for mid
	Mid PriceType
)

func (bid Bid) String() string {
	return "B"
}

func (ask Ask) String() string {
	return "A"
}

func (mid Mid) String() string {
	return "M"
}

type (
	// OrderType is an interface for order type
	OrderType interface {
		String() string
	}

	// Buy is an order type for long
	Buy OrderType
	// Sell is an order type for short
	Sell OrderType
)

func (buy Buy) String() string {
	return "BUY"
}

func (sell Sell) String() string {
	return "SELL"
}
