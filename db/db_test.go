package db

import (
	"log"
	"testing"
	"yukimaterrace/andaman/model"
)

const (
	tradeSetName    = "trade_set_name_1"
	algorithmParam1 = "{\"a\": 1, \"b\": 2}"
	algorithmParam2 = "{\"c\": 3, \"d\": 4}"
)

func TestTradeSet(t *testing.T) {
	name := tradeSetName
	ts, err := GetTradeSetByName(name)
	if err != nil {
		log.Println(err)
	} else {
		if err := DeleteTradeSet(ts.TradeSetID); err != nil {
			log.Fatal(err)
		}
	}

	err = AddTradeSet(name, model.Simulation)
	if err != nil {
		log.Fatal(err)
	}

	ts, err = GetTradeSetByName(name)
	log.Printf("%v", ts)

	tradeSets, err := GetTradeSetsByType(model.Simulation, 20, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", tradeSets)
}

func TestTradeAlgorithm(t *testing.T) {
	_type := model.Frame
	param := algorithmParam1
	tradeDirection := model.Short

	ta, err := GetTradeAlgorithmByTypeAndParam(_type, param)
	if err != nil {
		log.Println(err)
	} else {
		if err := DeleteTradeAlgorithm(ta.TradeAlgorithmID); err != nil {
			log.Fatal(err)
		}
	}

	err = AddTradeAlgorithm(_type, param, tradeDirection)
	if err != nil {
		log.Fatal(err)
	}

	ta, err = GetTradeAlgorithmByTypeAndParam(_type, param)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", ta)
}

func TestTradeSetConfiguration(t *testing.T) {
	if err := AddTradeAlgorithm(model.Frame, algorithmParam1, model.Short); err != nil {
		log.Println(err)
	}
	if err := AddTradeAlgorithm(model.Frame, algorithmParam2, model.Long); err != nil {
		log.Println(err)
	}

	ta1, err := GetTradeAlgorithmByTypeAndParam(model.Frame, algorithmParam1)
	if err != nil {
		log.Fatal(err)
	}

	ta2, err := GetTradeAlgorithmByTypeAndParam(model.Frame, algorithmParam2)
	if err != nil {
		log.Fatal(err)
	}

	if err := AddTradeConfiguration(model.GbpUsd, model.TokyoAM, ta1.TradeAlgorithmID); err != nil {
		log.Println(err)
	}

	if err := AddTradeConfiguration(model.UsdJpy, model.LondonAM, ta2.TradeAlgorithmID); err != nil {
		log.Println(err)
	}

	tc1, err := GetTradeConfigurationByFields(model.GbpUsd, model.TokyoAM, ta1.TradeAlgorithmID)
	if err != nil {
		log.Fatal(err)
	}

	tc2, err := GetTradeConfigurationByFields(model.UsdJpy, model.LondonAM, ta2.TradeAlgorithmID)
	if err != nil {
		log.Fatal(err)
	}

	ts, err := GetTradeSetByName(tradeSetName)
	if err != nil {
		log.Fatal(err)
	}

	AddTradeSetConfigurationRel(ts.TradeSetID, tc1.TradeConfigurationID)
	AddTradeSetConfigurationRel(ts.TradeSetID, tc2.TradeConfigurationID)
}
