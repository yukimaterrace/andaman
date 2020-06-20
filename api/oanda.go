package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"yukimaterrace/andaman/config"
)

// Oanda is OANDA API client
type Oanda struct {
	*client
}

// NewOanda is a constructor for Oanda
func NewOanda() *Oanda {
	host := "api-fxtrade.oanda.com"
	if config.OandaPractice {
		host = "api-fxpractice.oanda.com"
	}

	var header = http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", config.OandaToken))
	header.Add("Connection", "Keep-Alive")
	header.Add("Content-Type", "application/json")
	header.Add("Accept-Datetime-Format", "UNIX")

	return &Oanda{newClient(host, header)}
}

// GetAccounts is a method to get accounts
func (oanda *Oanda) GetAccounts() (*OandaAccounts, error) {
	var accounts OandaAccounts
	if err := oanda.client.get("/v3/accounts", nil, &accounts); err != nil {
		return nil, err
	}
	return &accounts, nil
}

// GetAccount is a method to get account
func (oanda *Oanda) GetAccount(accountID string) (*OandaAccount, error) {
	path := fmt.Sprintf("/v3/accounts/%s", accountID)

	var account OandaAccount
	if err := oanda.client.get(path, nil, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountChanges is a method to get account changes
func (oanda *Oanda) GetAccountChanges(accountID string, sinceTransactionID string) (*OandaAccountChanges, error) {
	path := fmt.Sprintf("/v3/accounts/%s/changes", accountID)

	query := url.Values{}
	query.Add("sinceTransactionID", sinceTransactionID)

	var accountChanges OandaAccountChanges
	if err := oanda.client.get(path, query, &accountChanges); err != nil {
		return nil, err
	}
	return &accountChanges, nil
}

// GetCandles is a method to get candles
func (oanda *Oanda) GetCandles(instrument string, granularity string, count int, from float64, includeFirst bool) (*OandaCandles, error) {
	path := fmt.Sprintf("/v3/instruments/%s/candles", instrument)

	query := url.Values{}
	query.Add("granularity", granularity)
	query.Add("count", oandaInt(count).String())
	if from > 0 {
		query.Add("from", oandaFloat64(from).String())
		query.Add("includeFirst", oandaBool(includeFirst).String())
	}

	var candles OandaCandles
	if err := oanda.client.get(path, query, &candles); err != nil {
		return nil, err
	}
	return &candles, nil
}

// GetPricing is a method to get pricing
func (oanda *Oanda) GetPricing(accountID string, instruments []string, since float64) (*OandaPrices, error) {
	path := fmt.Sprintf("/v3/accounts/%s/pricing", accountID)

	query := url.Values{}
	query.Add("instruments", strings.Join(instruments, ","))
	query.Add("since", oandaFloat64(since).String())

	var prices OandaPrices
	if err := oanda.client.get(path, query, &prices); err != nil {
		return nil, err
	}
	return &prices, nil
}

// GetLatestCandles is a method to get latest candles
func (oanda *Oanda) GetLatestCandles(accountID string, specs []string) (*OandaLatestCandles, error) {
	path := fmt.Sprintf("/v3/accounts/%s/candles/latest", accountID)

	query := url.Values{}
	query.Add("candleSpecifications", strings.Join(specs, ","))

	var latestCandles OandaLatestCandles
	if err := oanda.client.get(path, query, &latestCandles); err != nil {
		return nil, err
	}
	return &latestCandles, nil
}

// GetOpenTrades is a method to get open trades
func (oanda *Oanda) GetOpenTrades(accountID string) (*OandaTrades, error) {
	path := fmt.Sprintf("/v3/accounts/%s/openTrades", accountID)

	var trades OandaTrades
	if err := oanda.client.get(path, nil, &trades); err != nil {
		return nil, err
	}
	return &trades, nil
}

// PostOrder is a method to post order
func (oanda *Oanda) PostOrder(accountID string, orderType string, instrument string, units float64) (*OandaOrderCreated, error) {
	path := fmt.Sprintf("/v3/accounts/%s/orders", accountID)

	requestBody := oandaOrder{
		Order: oandaOrderRequest{
			OrderType:  orderType,
			Instrument: instrument,
			Units:      units,
		},
	}

	var orderCreated OandaOrderCreated
	if err := oanda.client.post(path, &requestBody, &orderCreated); err != nil {
		return nil, err
	}
	return &orderCreated, nil
}

// PutCloseTrade is a method to put close trade
func (oanda *Oanda) PutCloseTrade(accountID string, tradeID string) (*OandaTradeClosed, error) {
	path := fmt.Sprintf("/v3/accounts/%s/trades/%s/close", accountID, tradeID)

	requestBody := oandaCloseTrade{
		Units: "ALL",
	}

	var tradeClosed OandaTradeClosed
	if err := oanda.client.put(path, &requestBody, &tradeClosed); err != nil {
		return nil, err
	}
	return &tradeClosed, nil
}

func (oanda *Oanda) makeCandleSpecs(granularity string, instruments ...string) []string {
	specs := make([]string, 0)
	for _, instrument := range instruments {
		spec := fmt.Sprintf("%s:%s:M", instrument, granularity)
		specs = append(specs, spec)
	}
	return specs
}
