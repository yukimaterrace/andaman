package model

// GetTradeSets is a method to get trade sets
func GetTradeSets(_type TradeSetType, count int, offset int) (*TradeSetsResponse, error) {
	tradeSets, err := getTradeSetsByType(_type, count, offset)
	if err != nil {
		return nil, handleError(err)
	}

	return &TradeSetsResponse{TradeSets: tradeSets}, nil
}
