package model

// TradeSet is a struct for trade set table
type TradeSet struct {
	TradeSetID int          `json:"-"`
	Name       string       `json:"name"`
	Version    int          `json:"version"`
	Type       TradeSetType `json:"type"`
	CreatedAt  int          `json:"created_at"`
}

// TradeAlgorithm is a struct for trade algorithm table
type TradeAlgorithm struct {
	TradeAlgorithmID int                `json:"-"`
	Type             TradeAlgorithmType `json:"type"`
	ParamHash        string             `json:"-"`
	Param            string             `json:"-"`
	TradeDirection   TradeDirection     `json:"trade_direction"`
}

// TradeConfiguration is a struct for trade configuration table
type TradeConfiguration struct {
	TradeConfigurationID int       `json:"-"`
	TradePair            TradePair `json:"trade_pair"`
	Timezone             Timezone  `json:"timezone"`
	TradeAlgorithmID     int       `json:"-"`
}

// TradeSetConfigurationRel is a struct for trade set configurationrel table
type TradeSetConfigurationRel struct {
	TradeSetID           int
	TradeConfigurationID int
}

// Order is a struct for order table
type Order struct {
	OrderID              int            `json:"-"`
	TradeRunID           int            `json:"trade_run_id"`
	BrokerOrderID        int            `json:"order_id"`
	TradeConfigurationID int            `json:"-"`
	Units                float64        `json:"units"`
	TradeDirection       TradeDirection `json:"trade_direction"`
	State                OrderState     `json:"state"`
	Profit               float64        `json:"profit"`
	TimeAtOpen           int            `json:"time_at_open"`
	PriceAtOpen          float64        `json:"price_at_open"`
	TimeAtClose          int            `json:"time_at_close"`
	PriceAtClose         float64        `json:"price_at_close"`
}

// TradeRun is a struct for trade run table
type TradeRun struct {
	TradeRunID int           `json:"trade_run_id"`
	TradeSetID int           `json:"-"`
	Type       TradeRunType  `json:"type"`
	State      TradeRunState `json:"state"`
	CreatedAt  int           `json:"created_at"`
	StartAt    int           `json:"start_at"`
	FinishAt   int           `json:"finish_at"`
}
