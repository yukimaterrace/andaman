package api

import (
	"testing"
	"time"
	"yukimaterrace/andaman/config"
)

var oanda = NewOanda()

var accountID string
var lastTransactionID string

func TestGetAccounts(t *testing.T) {
	accounts, err := oanda.GetAccounts()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", accounts)

	accountID = accounts.Accounts[0].ID
}

func TestGetAccount(t *testing.T) {
	account, err := oanda.GetAccount(accountID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", account)

	lastTransactionID = account.LastTransactionID
}

func TestGetAccountChanges(t *testing.T) {
	accountChanges, err := oanda.GetAccountChanges(accountID, lastTransactionID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", accountChanges)
}

func TestGetCandlesLatest(t *testing.T) {
	candles, err := oanda.GetCandles("GBP_USD", "S5", 5, 0, true)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", candles)
}

func TestGetCandlesFrom(t *testing.T) {
	from := float64(time.Date(2020, 4, 1, 8, 0, 0, 0, time.UTC).Unix())
	candles, err := oanda.GetCandles("GBP_USD", "S5", 5, from, true)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", candles)
}

func TestGetPricing(t *testing.T) {
	since := float64(time.Date(2020, 4, 1, 8, 0, 0, 0, time.UTC).Unix())
	instruments := []string{"GBP_USD", "EUR_AUD"}

	prices, err := oanda.GetPricing(accountID, instruments, since)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", prices)
}

func TestGetLatestCandles(t *testing.T) {
	specs := oanda.makeCandleSpecs("S5", "GBP_USD", "EUR_USD")

	latestCandles, err := oanda.GetLatestCandles(accountID, specs)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", latestCandles)
}

func TestGetOpenTrades(t *testing.T) {
	trades, err := oanda.GetOpenTrades(accountID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", trades)
}

func TestOrder(t *testing.T) {
	if !config.OandaPractice {
		t.Skip("not practice mode")
	}

	units := 1000.0

	orderCreated, err := oanda.PostOrder(accountID, "MARKET", "GBP_USD", units)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", orderCreated)

	openTrades, _ := oanda.GetOpenTrades(accountID)
	t.Logf("%+v", openTrades)

	time.Sleep(time.Second)

	tradeID := openTrades.Trades[0].ID
	tradeClosed, err := oanda.PutCloseTrade(accountID, tradeID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", tradeClosed)
}
