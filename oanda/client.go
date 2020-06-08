package oanda

import (
	"fmt"
	"net/http"
	"net/url"
	"yukimaterrace/andaman/api"
	"yukimaterrace/andaman/config"
)

var client *api.Client

func init() {
	host := "api-fxtrade.oanda.com"
	if config.OandaPractice {
		host = "api-fxpractice.oanda.com"
	}

	var header = http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", config.OandaToken))
	header.Add("Connection", "Keep-Alive")
	header.Add("Content-Type", "application/json")
	header.Add("Accept-Datetime-Format", "UNIX")

	client = api.NewClient(host, header)
}

// GetAccounts is a method to get accounts
func GetAccounts() (*Accounts, error) {
	var accounts Accounts
	if err := client.Get("/v3/accounts", nil, &accounts); err != nil {
		return nil, err
	}
	return &accounts, nil
}

// GetAccount is a method to get account
func GetAccount(accountID string) (*Account, error) {
	path := fmt.Sprintf("/v3/accounts/%s", accountID)

	var account Account
	if err := client.Get(path, nil, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountChanges is a method to get account changes
func GetAccountChanges(accountID string, sinceTransactionID string) (*AccountChanges, error) {
	path := fmt.Sprintf("/v3/accounts/%s/changes", accountID)

	query := url.Values{}
	query.Add("sinceTransactionID", sinceTransactionID)

	var accountChanges AccountChanges
	if err := client.Get(path, query, &accountChanges); err != nil {
		return nil, err
	}
	return &accountChanges, nil
}

// GetCandles is a method to get candles
func GetCandles(instrument string, granularity string, count int, from float64, includeFirst bool) (*Candles, error) {
	path := fmt.Sprintf("/v3/instruments/%s/candles", instrument)

	query := url.Values{}
	query.Add("granularity", granularity)
	query.Add("count", Int(count).String())
	if from > 0 {
		query.Add("from", Float64(from).String())
		query.Add("includeFirst", Bool(includeFirst).String())
	}

	var candles Candles
	if err := client.Get(path, query, &candles); err != nil {
		return nil, err
	}
	return &candles, nil
}
