package trade

import (
	"yukimaterrace/andaman/indicate"
	"yukimaterrace/andaman/market"
)

type zoneTradeStrategy struct{}

func (strategy *zoneTradeStrategy) requireIndicators() []indicate.Indicator {
	return []indicate.Indicator{
		indicate.CreateSMA25(),
		indicate.CreateSMA75(),
		indicate.CreateSMA150(),
		indicate.CreateSMA600(),
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
	indicatorValues []*indicate.Value,
) (*market.MadeOrder, []*market.ClosedOrder) {

	return nil, nil
}
