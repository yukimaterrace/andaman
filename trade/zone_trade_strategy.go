package trade

import (
	"yukimaterrace/andaman/indicator"
	"yukimaterrace/andaman/market"
)

type zoneTradeStrategy struct{}

func (strategy *zoneTradeStrategy) requireIndicators() []indicator.Indicator {
	return []indicator.Indicator{
		indicator.CreateSMA25(),
		indicator.CreateSMA75(),
		indicator.CreateSMA150(),
		indicator.CreateSMA600(),
	}
}

func (strategy *zoneTradeStrategy) requireGranularity() market.Granularity {
	return market.M1
}

func (strategy *zoneTradeStrategy) requirePriceCount() int {
	return 500
}

func (strategy *zoneTradeStrategy) processOrder(
	orderer market.Orderer,
	tradePrice *market.TradePrice,
	indicatorValues []*indicator.Value,
) (*market.MadeOrder, []*market.ClosedOrder) {

	return nil, nil
}
