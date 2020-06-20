package api

import "strconv"

type (
	// OandaAccounts is a definition for oanda acount
	OandaAccounts struct {
		Accounts []OandaAccountProperties `json:"accounts"`
	}

	// OandaAccountProperties is a definition for oanda account properties
	OandaAccountProperties struct {
		ID string `json:"id"`
	}

	// OandaAccount is a definition for oanda account
	OandaAccount struct {
		LastTransactionID string `json:"lastTransactionID"`
	}

	// OandaAccountChanges is a definition for oanda account changes
	OandaAccountChanges struct {
		State             OandaAccountState `json:"state"`
		LastTransactionID string            `json:"lastTransactionID"`
	}

	// OandaAccountState is a definition for oanda account state
	OandaAccountState struct {
		UnrealizedPL float64 `json:"unrealizedPL,string"`
		NAV          float64 `json:"NAV,string"`
	}

	// OandaCandles is a definition for oanda candles
	OandaCandles struct {
		Instrument  string             `json:"instrument"`
		Granularity string             `json:"granularity"`
		Candles     []OandaCandleStick `json:"candles"`
	}

	// OandaCandleStick is a definition for oanda candle stick
	OandaCandleStick struct {
		Time     float64              `json:"time,string"`
		Bid      OandaCandleStickData `json:"bid"`
		Ask      OandaCandleStickData `json:"ask"`
		Mid      OandaCandleStickData `json:"mid"`
		Volume   int                  `json:"volume"`
		Complete bool                 `json:"complete"`
	}

	// OandaCandleStickData is a definition for oanda candle stick data
	OandaCandleStickData struct {
		O float64 `json:"o,string"`
		H float64 `json:"h,string"`
		L float64 `json:"l,string"`
		C float64 `json:"c,string"`
	}

	// OandaPrices is a definition for oanda prices
	OandaPrices struct {
		Prices []OandaClientPrice `json:"prices"`
		Time   float64            `json:"time,string"`
	}

	// OandaClientPrice is a definition for oanda client price
	OandaClientPrice struct {
		Instrument string       `json:"instrument"`
		Time       float64      `json:"time,string"`
		Tradable   bool         `json:"tradable"`
		Bids       []OandaPrice `json:"bids"`
		Asks       []OandaPrice `json:"asks"`
	}

	// OandaPrice is a definition for oanda price
	OandaPrice struct {
		Price     float64 `json:"price,string"`
		Liquidity int     `json:"liquidity"`
	}

	// OandaLatestCandles is a definition for oanda latest candles
	OandaLatestCandles struct {
		LatestCandles []OandaCandles `json:"latestCandles"`
	}

	// OandaTrades is a definition for oanda trades
	OandaTrades struct {
		Trades            []OandaTrade `json:"trades"`
		LastTransactionID string       `json:"lastTransactionID"`
	}

	// OandaTrade is a definition for oanda trade
	OandaTrade struct {
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

	// OandaOrderCreated is a definition for oanda order created
	OandaOrderCreated struct {
		LastTransactionID string `json:"lastTransactionID"`
	}

	// OandaTradeClosed is a definition for oanda trade closed
	OandaTradeClosed struct {
		LastTransactionID string `json:"lastTransactionID"`
	}
)

type (
	oandaOrder struct {
		Order oandaOrderRequest `json:"order"`
	}

	oandaOrderRequest struct {
		OrderType  string  `json:"type"`
		Instrument string  `json:"instrument"`
		Units      float64 `json:"units,string"`
	}

	oandaCloseTrade struct {
		Units string `json:"units"`
	}
)

type (
	oandaFloat64 float64
	oandaInt     int
	oandaBool    bool
)

func (f oandaFloat64) String() string {
	return strconv.FormatFloat(float64(f), 'f', 10, 64)
}

func (i oandaInt) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (b oandaBool) String() string {
	return strconv.FormatBool(bool(b))
}
