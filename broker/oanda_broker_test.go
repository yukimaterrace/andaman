package broker

import (
	"testing"
	"time"
)

var oanda = NewOandaBroker()

var lastTransactionID int

func TestAccounts(t *testing.T) {
	accounts, err := oanda.Accounts()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", accounts)
}

func TestAccount(t *testing.T) {
	account, err := oanda.Account()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", account)

	lastTransactionID = account.LastTransactionID
}

func TestAccountChanges(t *testing.T) {
	accountChanges, err := oanda.AccountChanges(lastTransactionID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", accountChanges)
}

func TestCandlesLatest(t *testing.T) {
	candles, err := oanda.Candles("GBP_USD", "M1", 5, 0, 0, true)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", candles)
}

func TestCandlesFrom(t *testing.T) {
	from := time.Date(2020, 4, 1, 8, 0, 0, 0, time.UTC).Unix()
	candles, err := oanda.Candles("GBP_USD", "M1", 5, int(from), 0, true)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", candles)
}

func TestCandlesFromTo(t *testing.T) {
	from := time.Date(2020, 4, 1, 8, 0, 0, 0, time.UTC).Unix()
	to := time.Date(2020, 4, 1, 8, 5, 0, 0, time.UTC).Unix()
	candles, err := oanda.Candles("GBP_USD", "M1", 0, int(from), int(to), true)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", candles)
}

func TestPricing(t *testing.T) {
	since := time.Date(2020, 4, 1, 8, 0, 0, 0, time.UTC).Unix()
	instruments := []string{"GBP_USD", "EUR_AUD"}

	prices, err := oanda.Pricing(instruments, int(since))
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", prices)
}

func TestGetLatestCandles(t *testing.T) {
	specs := oanda.makeCandleSpecs("M1", "GBP_USD", "EUR_USD")

	latestCandles, err := oanda.LatestCandles(specs)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", latestCandles)
}

func TestGetOpenTrades(t *testing.T) {
	trades, err := oanda.OpenTrades()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", trades)
}

func TestOrder(t *testing.T) {
	if oanda.client.host != "api-fxpractice.oanda.com" {
		t.Skip("not practice mode")
	}

	units := 1000.0

	orderCreated, err := oanda.CreateOrder("MARKET", "GBP_USD", units)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", orderCreated)

	openTrades, _ := oanda.OpenTrades()
	t.Logf("%+v", openTrades)

	time.Sleep(time.Second)

	tradeID := openTrades.Trades[0].ID
	tradeClosed, err := oanda.CloseTrade(tradeID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", tradeClosed)
}
