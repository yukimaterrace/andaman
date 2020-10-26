package model

// GetTradeSets is a method to get trade sets
func GetTradeSets(_type TradeSetType, count int, offset int) (*TradeSetsResponse, error) {
	tradeSets, err := getTradeSetsByType(_type, count, offset)
	if err != nil {
		return nil, handleError(err)
	}

	return &TradeSetsResponse{TradeSets: tradeSets}, nil
}

// GetTradeSetDetail is a method to get trade set detail
func GetTradeSetDetail(name string, paramObjectCreator TradeParamObjectCreator) (*TradeSetDetail, error) {
	tradeSet, err := getTradeSetByName(name)
	if err != nil {
		return nil, err
	}

	tradeConfigurationDetails, err := getTradeConfigurationDetailsByTradeSetID(tradeSet.TradeSetID)
	if err != nil {
		return nil, err
	}

	for _, detail := range tradeConfigurationDetails {
		paramObject, err := paramObjectCreator(detail.Algorithm.Type, detail.Algorithm.Param)
		if err != nil {
			return nil, err
		}
		detail.Algorithm.ParamObject = paramObject
	}

	tradeSetDetail := &TradeSetDetail{
		TradeSet:       tradeSet,
		Configurations: tradeConfigurationDetails,
	}

	return tradeSetDetail, nil
}
