package oanda

import (
	"testing"
	"time"
)

var accountID string
var lastTransactionID string

func TestGetAccounts(t *testing.T) {
	accounts, err := GetAccounts()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", accounts)

	accountID = accounts.Accounts[0].ID
}

func TestGetAccount(t *testing.T) {
	account, err := GetAccount(accountID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", account)

	lastTransactionID = account.LastTransactionID
}

func TestGetAccountChanges(t *testing.T) {
	accountChanges, err := GetAccountChanges(accountID, lastTransactionID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", accountChanges)
}

func TestGetCandlesLatest(t *testing.T) {
	candles, err := GetCandles("GBP_USD", "S5", 5, 0, true)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", candles)
}

func TestGetCandlesFrom(t *testing.T) {
	from := float64(time.Date(2020, 4, 1, 8, 0, 0, 0, time.UTC).Unix())
	candles, err := GetCandles("GBP_USD", "S5", 5, from, true)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", candles)
}
