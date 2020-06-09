package oanda

// Accounts is a definition for accounts
type Accounts struct {
	Accounts []AccountProperties `json:"accounts"`
}

// AccountProperties is a definition for account properties
type AccountProperties struct {
	ID string `json:"id"`
}

// Account is a definition for account
type Account struct {
	LastTransactionID string `json:"lastTransactionID"`
}

// AccountChanges is a definition for account changes
type AccountChanges struct {
	State             AccountState `json:"state"`
	LastTransactionID string       `json:"lastTransactionID"`
}

// AccountState is a definition for account state
type AccountState struct {
	UnrealizedPL float64 `json:"unrealizedPL,string"`
	NAV          float64 `json:"NAV,string"`
}

// Candles is a definition for candle
type Candles struct {
	Instrument  string        `json:"instrument"`
	Granularity string        `json:"granularity"`
	Candles     []CandleStick `json:"candles"`
}

// CandleStick is a definition for candle stick
type CandleStick struct {
	Time     float64         `json:"time,string"`
	Bid      CandleStickData `json:"bid"`
	Ask      CandleStickData `json:"ask"`
	Mid      CandleStickData `json:"mid"`
	Volume   int             `json:"volume"`
	Complete bool            `json:"complete"`
}

// CandleStickData is a definition for candle stick data
type CandleStickData struct {
	O float64 `json:"o,string"`
	H float64 `json:"h,string"`
	L float64 `json:"l,string"`
	C float64 `json:"c,string"`
}

// Prices is a definition for prices
type Prices struct {
	Prices []ClientPrice `json:"prices"`
	Time   float64       `json:"time,string"`
}

// ClientPrice is a definition for client price
type ClientPrice struct {
	Instrument string  `json:"instrument"`
	Time       float64 `json:"time,string"`
	Tradable   bool    `json:"tradable"`
	Bids       []Price `json:"bids"`
	Asks       []Price `json:"asks"`
}

// Price is a definition for price
type Price struct {
	Price     float64 `json:"price,string"`
	Liquidity int     `json:"liquidity"`
}

// LatestCandles is a definition for latest candles
type LatestCandles struct {
	LatestCandles []Candles `json:"latestCandles"`
}

// Trades is a definition for trades
type Trades struct {
	Trades            []Trade `json:"trades"`
	LastTransactionID string  `json:"lastTransactionID"`
}

// Trade is a definition for trade
type Trade struct {
	ID           string  `json:"id"`
	Instrument   string  `json:"instrument"`
	Price        float64 `json:"price,string"`
	OpenTime     float64 `json:"openTime,string"`
	State        string  `json:"state"`
	InitialUnits float64 `json:"initialUnits,string"`
	CurrentUnits float64 `json:"currentUnits,string"`
	RealizedPL   float64 `json:"realizedPL,string"`
	UnrealizedPL float64 `json:"unrealizedPL,string"`
}

// OrderCreated is a definition for order created
type OrderCreated struct {
	LastTransactionID string `json:"lastTransactionID"`
}

// TradeClosed is a definition for trade closed
type TradeClosed struct {
	LastTransactionID string `json:"lastTransactionID"`
}
