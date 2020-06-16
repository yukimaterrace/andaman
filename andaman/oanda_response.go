package andaman

type (
	oandaAccounts struct {
		Accounts []oandaAccountProperties `json:"accounts"`
	}

	oandaAccountProperties struct {
		ID string `json:"id"`
	}

	oandaAccount struct {
		LastTransactionID string `json:"lastTransactionID"`
	}

	oandaAccountChanges struct {
		State             oandaAccountState `json:"state"`
		LastTransactionID string            `json:"lastTransactionID"`
	}

	oandaAccountState struct {
		UnrealizedPL float64 `json:"unrealizedPL,string"`
		NAV          float64 `json:"NAV,string"`
	}

	oandaCandles struct {
		Instrument  string             `json:"instrument"`
		Granularity string             `json:"granularity"`
		Candles     []oandaCandleStick `json:"candles"`
	}

	oandaCandleStick struct {
		Time     float64              `json:"time,string"`
		Bid      oandaCandleStickData `json:"bid"`
		Ask      oandaCandleStickData `json:"ask"`
		Mid      oandaCandleStickData `json:"mid"`
		Volume   int                  `json:"volume"`
		Complete bool                 `json:"complete"`
	}

	oandaCandleStickData struct {
		O float64 `json:"o,string"`
		H float64 `json:"h,string"`
		L float64 `json:"l,string"`
		C float64 `json:"c,string"`
	}

	oandaPrices struct {
		Prices []oandaClientPrice `json:"prices"`
		Time   float64            `json:"time,string"`
	}

	oandaClientPrice struct {
		Instrument string       `json:"instrument"`
		Time       float64      `json:"time,string"`
		Tradable   bool         `json:"tradable"`
		Bids       []oandaPrice `json:"bids"`
		Asks       []oandaPrice `json:"asks"`
	}

	oandaPrice struct {
		Price     float64 `json:"price,string"`
		Liquidity int     `json:"liquidity"`
	}

	oandaLatestCandles struct {
		LatestCandles []oandaCandles `json:"latestCandles"`
	}

	oandaTrades struct {
		Trades            []oandaTrade `json:"trades"`
		LastTransactionID string       `json:"lastTransactionID"`
	}

	oandaTrade struct {
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

	oandaOrderCreated struct {
		LastTransactionID string `json:"lastTransactionID"`
	}

	oandaTradeClosed struct {
		LastTransactionID string `json:"lastTransactionID"`
	}
)
