package model

import (
	"log"
	"testing"
)

const (
	tradeSetName    = "trade_set_name_1"
	algorithmParam1 = "{\"a\": 1, \"b\": 2}"
	algorithmParam2 = "{\"c\": 3, \"d\": 4}"
)

func TestTradeSet(t *testing.T) {
	name := tradeSetName
	ts, err := getTradeSetByName(name)
	if err != nil {
		log.Println(err)
	} else {
		if err := deleteTradeSet(ts.TradeSetID); err != nil {
			log.Fatal(err)
		}
	}

	err = addTradeSet(name, Simulation)
	if err != nil {
		log.Fatal(err)
	}

	ts, err = getTradeSetByName(name)
	log.Printf("%v", ts)

	tradeSets, err := getTradeSetsByType(Simulation, 20, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", tradeSets)
}

func TestTradeAlgorithm(t *testing.T) {
	_type := Frame
	param := algorithmParam1
	tradeDirection := Short

	ta, err := getTradeAlgorithmByTypeAndParam(_type, param)
	if err != nil {
		log.Println(err)
	} else {
		if err := deleteTradeAlgorithm(ta.TradeAlgorithmID); err != nil {
			log.Fatal(err)
		}
	}

	err = addTradeAlgorithm(_type, param, tradeDirection)
	if err != nil {
		log.Fatal(err)
	}

	ta, err = getTradeAlgorithmByTypeAndParam(_type, param)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", ta)
}

func TestTradeSetConfiguration(t *testing.T) {
	if err := addTradeAlgorithm(Frame, algorithmParam1, Short); err != nil {
		log.Println(err)
	}
	if err := addTradeAlgorithm(Frame, algorithmParam2, Long); err != nil {
		log.Println(err)
	}

	ta1, err := getTradeAlgorithmByTypeAndParam(Frame, algorithmParam1)
	if err != nil {
		log.Fatal(err)
	}

	ta2, err := getTradeAlgorithmByTypeAndParam(Frame, algorithmParam2)
	if err != nil {
		log.Fatal(err)
	}

	if err := addTradeConfiguration(GbpUsd, TokyoAM, ta1.TradeAlgorithmID); err != nil {
		log.Println(err)
	}

	if err := addTradeConfiguration(UsdJpy, LondonAM, ta2.TradeAlgorithmID); err != nil {
		log.Println(err)
	}

	tc1, err := getTradeConfigurationByFields(GbpUsd, TokyoAM, ta1.TradeAlgorithmID)
	if err != nil {
		log.Fatal(err)
	}

	tc2, err := getTradeConfigurationByFields(UsdJpy, LondonAM, ta2.TradeAlgorithmID)
	if err != nil {
		log.Fatal(err)
	}

	ts, err := getTradeSetByName(tradeSetName)
	if err != nil {
		log.Fatal(err)
	}

	addTradeSetConfigurationRel(ts.TradeSetID, tc1.TradeConfigurationID)
	addTradeSetConfigurationRel(ts.TradeSetID, tc2.TradeConfigurationID)
}
