package andaman

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
