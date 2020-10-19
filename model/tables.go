package model

// TradeSet is a struct for trade set table
type TradeSet struct {
	TradeSetID int
	Name       string
	Type       TradeSetType
	State      TradeSetState
	CreatedAt  int
	UpdatedAt  int
}

// TradeAlgorithm is a struct for trade algorithm table
type TradeAlgorithm struct {
	TradeAlgorithmID int
	Type             TradeAlgorithmType
	ParamHash        string
	Param            string
	TradeDirection   TradeDirection
}

// TradeConfiguration is a struct for trade configuration table
type TradeConfiguration struct {
	TradeConfigurationID int
	TradePair            TradePair
	Timezone             Timezone
	TradeAlgorithmID     int
}

// TradeSetConfigurationRel is a struct for trade set configurationrel table
type TradeSetConfigurationRel struct {
	TradeSetID           int
	TradeConfigurationID int
}

// Order is a struct for order table
type Order struct {
	OrderID              int
	BrokerOrderID        int
	TradeConfigurationID int
	Units                float64
	State                OrderState
	Profit               float64
	TimeAtOpen           int
	PriceAtOpen          float64
	TimeAtClose          int
	PriceAtClose         float64
}
