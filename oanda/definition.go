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
