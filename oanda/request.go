package oanda

// Order is a definition for order
type Order struct {
	Order OrderRequest `json:"order"`
}

// OrderRequest is a definition for order request
type OrderRequest struct {
	Type       string  `json:"type"`
	Instrument string  `json:"instrument"`
	Units      float64 `json:"units,string"`
}

// CloseTrade is a definition for close trade
type CloseTrade struct {
	Units string `json:"units"`
}
