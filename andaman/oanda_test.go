package andaman

import (
	"testing"
	"time"
	"yukimaterrace/andaman/config"
)

var oandaInstance = newOanda()

var accountID string
var lastTransactionID string

func TestGetAccounts(t *testing.T) {
	accounts, err := oandaInstance.getAccounts()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", accounts)

	accountID = accounts.Accounts[0].ID
}

func TestGetAccount(t *testing.T) {
	account, err := oandaInstance.getAccount(accountID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", account)

	lastTransactionID = account.LastTransactionID
}

func TestGetAccountChanges(t *testing.T) {
	accountChanges, err := oandaInstance.getAccountChanges(accountID, lastTransactionID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", accountChanges)
}

func TestGetCandlesLatest(t *testing.T) {
	candles, err := oandaInstance.getCandles("GBP_USD", "S5", 5, 0, true)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", candles)
}

func TestGetCandlesFrom(t *testing.T) {
	from := float64(time.Date(2020, 4, 1, 8, 0, 0, 0, time.UTC).Unix())
	candles, err := oandaInstance.getCandles("GBP_USD", "S5", 5, from, true)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", candles)
}

func TestGetPricing(t *testing.T) {
	since := float64(time.Date(2020, 4, 1, 8, 0, 0, 0, time.UTC).Unix())
	instruments := []string{"GBP_USD", "EUR_AUD"}

	prices, err := oandaInstance.getPricing(accountID, instruments, since)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", prices)
}

func TestGetLatestCandles(t *testing.T) {
	specs := oandaInstance.makeCandleSpecs("S5", "GBP_USD", "EUR_USD")

	latestCandles, err := oandaInstance.getLatestCandles(accountID, specs)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", latestCandles)
}

func TestGetOpenTrades(t *testing.T) {
	trades, err := oandaInstance.getOpenTrades(accountID)
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

	orderCreated, err := oandaInstance.postOrder(accountID, "MARKET", "GBP_USD", units)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", orderCreated)

	openTrades, _ := oandaInstance.getOpenTrades(accountID)
	t.Logf("%+v", openTrades)

	time.Sleep(time.Second)

	tradeID := openTrades.Trades[0].ID
	tradeClosed, err := oandaInstance.putCloseTrade(accountID, tradeID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", tradeClosed)
}
