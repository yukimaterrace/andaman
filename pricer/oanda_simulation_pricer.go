package pricer

import (
	"log"
	"time"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/model"
)

type oandaSimulationPricer struct {
	start           int
	end             int
	seed            oandaSimulationPriceSeed
	currentIndexMap map[model.TradePair]int
	currentTime     int64
	granularitySec  int64
	unitSize        int
	initialized     bool
}

func newOandaSimulationPricer(tradePairs []model.TradePair, start int, end int) *oandaSimulationPricer {
	unitSize := 250

	currentTime := int64(0)
	currentIndexMap := make(map[model.TradePair]int)

	for _, tradePair := range tradePairs {
		currentIndexMap[tradePair] = unitSize
	}

	return &oandaSimulationPricer{
		start:           start,
		end:             end,
		seed:            nil,
		currentIndexMap: currentIndexMap,
		currentTime:     currentTime,
		granularitySec:  60,
		unitSize:        unitSize,
		initialized:     false,
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
	feedCandlesMap := make(map[model.TradePair]*broker.OandaCandles)

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

func (pricer *oandaSimulationPricer) Initialize() {
	var tradePairs []model.TradePair
	for tradePair := range pricer.currentIndexMap {
		tradePairs = append(tradePairs, tradePair)
	}

	pricer.seed = fetchOandaSimulationPriceSeed(tradePairs, "M1", pricer.start, pricer.end)

	for _, candles := range pricer.seed {
		time := int64(candles.Candles[pricer.unitSize-1].Time)
		if pricer.currentTime == 0 || pricer.currentTime > time {
			pricer.currentTime = time
		}
	}
	pricer.initialized = true
}

func (pricer *oandaSimulationPricer) CreatePrice(done chan<- *flow.CreatePriceResult) {
	if !pricer.initialized {
		panic(model.ErrInconsistentLogic)
	}

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
func (factory *OandaSimulationPricerFactory) Create(broker broker.Broker, tradePairs []model.TradePair) flow.Pricer {
	return newOandaSimulationPricer(tradePairs, factory.start, factory.end)
}

type oandaSimulationPrice struct {
	*oandaPrice
	candlesMap map[model.TradePair]*broker.OandaCandles
}

func newOandaSimulationPrice(candlesMap map[model.TradePair]*broker.OandaCandles, priceTime int64) *oandaSimulationPrice {
	return &oandaSimulationPrice{
		oandaPrice: newOandaPrice(candlesMap, nil, priceTime),
		candlesMap: candlesMap,
	}
}

func (oandaSimulationPrice *oandaSimulationPrice) Price(tradePair model.TradePair) broker.Price {
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

type oandaSimulationPriceSeed map[model.TradePair]*broker.OandaCandles

func fetchOandaSimulationPriceSeed(tradePairs []model.TradePair, granularity string, start int, end int) oandaSimulationPriceSeed {
	client := broker.NewOandaBroker()

	log.Println("start fetch candles...")

	startTime := time.Now().Unix()
	requestCount := 0

	candlesMap := map[model.TradePair]*broker.OandaCandles{}

	remainedTradePairs := map[model.TradePair]bool{}
	fromMap := map[model.TradePair]int{}
	includeFirstMap := map[model.TradePair]bool{}
	resultMap := map[model.TradePair]chan *broker.OandaCandles{}

	for _, pair := range tradePairs {
		remainedTradePairs[pair] = true
		fromMap[pair] = start
		includeFirstMap[pair] = true
		resultMap[pair] = make(chan *broker.OandaCandles, 1)
	}

	for len(remainedTradePairs) > 0 {
		pairs := map[model.TradePair]bool{}
		for pair := range remainedTradePairs {
			pairs[pair] = true
		}

		requestedPairs := map[model.TradePair]bool{}
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

			go func(result chan<- *broker.OandaCandles, tradePair model.TradePair) {
				candles, err := client.Candles(tradePair.OandaInstrument(), granularity, count, from, to, includeFirstMap[tradePair])
				if err != nil {
					log.Println(err.Error())
					result <- nil
					return
				}
				result <- candles
			}(resultMap[pair], pair)

			requestedPairs[pair] = true
		}

		requestCount += len(requestedPairs)

		for pair := range requestedPairs {
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

		for pair, t := range fromMap {
			tm := time.Unix(int64(t), 0)
			log.Printf("%s %s %d:%d:%d\n", pair.String(), tm.Month().String(), tm.Day(), tm.Hour(), tm.Minute())
		}
	}

	for pair, candles := range candlesMap {
		log.Printf("fetch %s: %d sticks\n", pair.String(), len(candles.Candles))
	}

	endTime := time.Now().Unix()
	log.Printf("finish fetch candles / %d requets for %d s\n", requestCount, endTime-startTime)

	return candlesMap
}
