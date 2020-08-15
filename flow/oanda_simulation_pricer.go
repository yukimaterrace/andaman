package flow

import (
	"log"
	"time"
	"yukimaterrace/andaman/broker"
)

type oandaSimulationPricerSeed struct {
	candlesMap map[broker.TradePair]*broker.OandaCandles
}

func fetchOandaSimulationPricerSeed(tradePairs []broker.TradePair, granularity string, start int, end int) *oandaSimulationPricerSeed {
	client := broker.NewOandaClient()

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

	return &oandaSimulationPricerSeed{candlesMap: candlesMap}
}
