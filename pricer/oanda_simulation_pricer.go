package pricer

import (
	"log"
	"time"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
)

type oandaSimulationPricer struct {
	seed            oandaSimulationPriceSeed
	currentIndexMap map[broker.TradePair]int
	currentTime     int64
	granularitySec  int64
	unitSize        int
}

func newOandaSimulationPricer(seed oandaSimulationPriceSeed) *oandaSimulationPricer {
	unitSize := 250

	currentTime := int64(0)
	currentIndexMap := make(map[broker.TradePair]int)

	for pair, candles := range seed {
		currentIndexMap[pair] = unitSize

		time := int64(candles.Candles[unitSize-1].Time)
		if currentTime == 0 || currentTime > time {
			currentTime = time
		}
	}

	return &oandaSimulationPricer{
		seed:            seed,
		currentIndexMap: currentIndexMap,
		currentTime:     currentTime,
		granularitySec:  60,
		unitSize:        unitSize,
	}
}

func (pricer *oandaSimulationPricer) hasNext() bool {
	for pair, candles := range pricer.seed {
		if pricer.currentIndexMap[pair] < len(candles.Candles) {
			return true
		}
	}
	return false
}

func (pricer *oandaSimulationPricer) next() *oandaSimulationPrice {
	feedCandlesMap := make(map[broker.TradePair]*broker.OandaCandles)

	for pair, candles := range pricer.seed {
		currentIndex := pricer.currentIndexMap[pair]

		feedCandles := &broker.OandaCandles{
			Instrument:  candles.Instrument,
			Granularity: candles.Granularity,
			Candles:     candles.Candles[currentIndex-pricer.unitSize : currentIndex],
		}
		feedCandlesMap[pair] = feedCandles

		if currentIndex+1 <= len(candles.Candles) && pricer.currentTime+pricer.granularitySec >= int64(candles.Candles[currentIndex].Time) {
			pricer.currentIndexMap[pair] = currentIndex + 1
		}
	}

	createPrice := newOandaSimulationPrice(feedCandlesMap, pricer.currentTime)
	pricer.currentTime += pricer.granularitySec

	return createPrice
}

func (pricer *oandaSimulationPricer) CreatePrice(done chan<- *flow.CreatePriceResult) {
	if !pricer.hasNext() {
		done <- &flow.CreatePriceResult{
			TradeMaterial: nil,
			Err:           flow.ErrNoMorePrice,
		}
		return
	}

	done <- &flow.CreatePriceResult{
		TradeMaterial: pricer.next(),
		Err:           nil,
	}
}

// OandaSimulationPricerFactory is a pricer factory for oanda simulation pricer
type OandaSimulationPricerFactory struct {
	start int
	end   int
}

// NewOandaSimulationPricerFactory is a constructor for OandaSimulationPricerFactory
func NewOandaSimulationPricerFactory(startTime time.Time, endTime time.Time) *OandaSimulationPricerFactory {
	return &OandaSimulationPricerFactory{
		start: int(startTime.Unix()),
		end:   int(endTime.Unix()),
	}
}

// Create is a factory method to create oanda simulation pricer factory
func (factory *OandaSimulationPricerFactory) Create(broker broker.Broker, tradePairs []broker.TradePair) flow.Pricer {
	seed := fetchOandaSimulationPriceSeed(tradePairs, "M1", factory.start, factory.end)
	return newOandaSimulationPricer(seed)
}

type oandaSimulationPrice struct {
	*oandaPrice
	candlesMap map[broker.TradePair]*broker.OandaCandles
}

func newOandaSimulationPrice(candlesMap map[broker.TradePair]*broker.OandaCandles, priceTime int64) *oandaSimulationPrice {
	return &oandaSimulationPrice{
		oandaPrice: newOandaPrice(candlesMap, nil, priceTime),
		candlesMap: candlesMap,
	}
}

func (oandaSimulationPrice *oandaSimulationPrice) Price(tradePair broker.TradePair) broker.Price {
	candles, ok := oandaSimulationPrice.candlesMap[tradePair]
	if !ok {
		log.Panicf("no candles exist for %v\n", tradePair)
	}

	candleStick := candles.Candles[len(candles.Candles)-1]

	return &price{
		bid: candleStick.Bid.C,
		ask: candleStick.Ask.C,
	}
}

type price struct {
	bid float64
	ask float64
}

func (price *price) Bid() float64 {
	return price.bid
}

func (price *price) Ask() float64 {
	return price.ask
}

type oandaSimulationPriceSeed map[broker.TradePair]*broker.OandaCandles

func fetchOandaSimulationPriceSeed(tradePairs []broker.TradePair, granularity string, start int, end int) oandaSimulationPriceSeed {
	client := broker.NewOandaBroker()

	log.Println("start fetch candles...")

	startTime := time.Now().Unix()
	requestCount := 0

	candlesMap := map[broker.TradePair]*broker.OandaCandles{}

	remainedTradePairs := map[broker.TradePair]bool{}
	fromMap := make(map[broker.TradePair]int)
	includeFirstMap := map[broker.TradePair]bool{}
	resultMap := map[broker.TradePair]chan *broker.OandaCandles{}

	for _, pair := range tradePairs {
		remainedTradePairs[pair] = true
		fromMap[pair] = start
		includeFirstMap[pair] = true
		resultMap[pair] = make(chan *broker.OandaCandles, 1)
	}

	for len(remainedTradePairs) > 0 {
		pairs := map[broker.TradePair]bool{}

		for pair := range remainedTradePairs {
			pairs[pair] = true
		}

		for pair := range pairs {
			count := 5000
			from := fromMap[pair]
			to := end

			if to-from > count*60 {
				to = 0
			} else if to <= from {
				delete(remainedTradePairs, pair)
				continue
			} else {
				count = 0
			}

			go func(result chan<- *broker.OandaCandles, tradePair broker.TradePair) {
				candles, err := client.Candles(string(tradePair), granularity, count, from, to, includeFirstMap[tradePair])
				if err != nil {
					log.Println(err.Error())
					result <- nil
					return
				}
				result <- candles
			}(resultMap[pair], pair)
		}

		requestCount += len(pairs)

		for pair := range pairs {
			resp := <-resultMap[pair]
			if resp == nil {
				delete(remainedTradePairs, pair)
				continue
			}

			if candles, ok := candlesMap[pair]; !ok {
				candlesMap[pair] = resp
			} else {
				candles.Candles = append(candles.Candles, resp.Candles...)
			}

			if len(resp.Candles) == 0 {
				delete(remainedTradePairs, pair)
				continue
			}

			fromMap[pair] = int(resp.Candles[len(resp.Candles)-1].Time)
			includeFirstMap[pair] = false
		}
	}

	for pair, candles := range candlesMap {
		log.Printf("fetch %s: %d sticks\n", string(pair), len(candles.Candles))
	}

	endTime := time.Now().Unix()
	log.Printf("finish fetch candles / %d requets for %d s\n", requestCount, endTime-startTime)

	return candlesMap
}
