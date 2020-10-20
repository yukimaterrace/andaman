package model

import (
	"log"
	"testing"
)

func TestTradeSet(t *testing.T) {
	name := "trade_set_name_1"
	err := deleteTradeSetByName(name)
	if err != nil {
		log.Println(err)
	}

	err = addTradeSet(name, Simulation, Stopped)
	if err != nil {
		log.Fatal(err)
	}

	ts, err := getTradeSetByName(name)
	log.Printf("%v", ts)

	err = updateTradeSetByName(name, Running)
	if err != nil {
		log.Fatal(err)
	}

	tradeSets, err := getTradeSetsByType(Simulation, 20, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", tradeSets)
}

func TestTradeAlgorithm(t *testing.T) {
	_type := Frame
	param := "{\"a\": 1, \"b\": 2}"
	tradeDirection := Short

	err := deleteTradeAlgorithmByTypeAndParam(_type, param)
	if err != nil {
		log.Println(err)
	}

	err = addTradeAlgorithm(_type, param, tradeDirection)
	if err != nil {
		log.Fatal(err)
	}

	ta, err := getTradeAlgorithmByTypeAndParam(_type, param)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", ta)
}
