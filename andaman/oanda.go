package andaman

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"yukimaterrace/andaman/config"
)

type oanda struct {
	client *apiClient
}

func newOanda() *oanda {
	host := "api-fxtrade.oanda.com"
	if config.OandaPractice {
		host = "api-fxpractice.oanda.com"
	}

	var header = http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", config.OandaToken))
	header.Add("Connection", "Keep-Alive")
	header.Add("Content-Type", "application/json")
	header.Add("Accept-Datetime-Format", "UNIX")

	return &oanda{
		client: newApiClient(host, header),
	}
}

func (oanda *oanda) getAccounts() (*oandaAccounts, error) {
	var accounts oandaAccounts
	if err := oanda.client.get("/v3/accounts", nil, &accounts); err != nil {
		return nil, err
	}
	return &accounts, nil
}

func (oanda *oanda) getAccount(accountID string) (*oandaAccount, error) {
	path := fmt.Sprintf("/v3/accounts/%s", accountID)

	var account oandaAccount
	if err := oanda.client.get(path, nil, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func (oanda *oanda) getAccountChanges(accountID string, sinceTransactionID string) (*oandaAccountChanges, error) {
	path := fmt.Sprintf("/v3/accounts/%s/changes", accountID)

	query := url.Values{}
	query.Add("sinceTransactionID", sinceTransactionID)

	var accountChanges oandaAccountChanges
	if err := oanda.client.get(path, query, &accountChanges); err != nil {
		return nil, err
	}
	return &accountChanges, nil
}

func (oanda *oanda) getCandles(instrument string, granularity string, count int, from float64, includeFirst bool) (*oandaCandles, error) {
	path := fmt.Sprintf("/v3/instruments/%s/candles", instrument)

	query := url.Values{}
	query.Add("granularity", granularity)
	query.Add("count", oandaInt(count).String())
	if from > 0 {
		query.Add("from", oandaFloat64(from).String())
		query.Add("includeFirst", oandaBool(includeFirst).String())
	}

	var candles oandaCandles
	if err := oanda.client.get(path, query, &candles); err != nil {
		return nil, err
	}
	return &candles, nil
}

func (oanda *oanda) getPricing(accountID string, instruments []string, since float64) (*oandaPrices, error) {
	path := fmt.Sprintf("/v3/accounts/%s/pricing", accountID)

	query := url.Values{}
	query.Add("instruments", strings.Join(instruments, ","))
	query.Add("since", oandaFloat64(since).String())

	var prices oandaPrices
	if err := oanda.client.get(path, query, &prices); err != nil {
		return nil, err
	}
	return &prices, nil
}

func (oanda *oanda) getLatestCandles(accountID string, specs []string) (*oandaLatestCandles, error) {
	path := fmt.Sprintf("/v3/accounts/%s/candles/latest", accountID)

	query := url.Values{}
	query.Add("candleSpecifications", strings.Join(specs, ","))

	var latestCandles oandaLatestCandles
	if err := oanda.client.get(path, query, &latestCandles); err != nil {
		return nil, err
	}
	return &latestCandles, nil
}

func (oanda *oanda) getOpenTrades(accountID string) (*oandaTrades, error) {
	path := fmt.Sprintf("/v3/accounts/%s/openTrades", accountID)

	var trades oandaTrades
	if err := oanda.client.get(path, nil, &trades); err != nil {
		return nil, err
	}
	return &trades, nil
}

func (oanda *oanda) postOrder(accountID string, orderType string, instrument string, units float64) (*oandaOrderCreated, error) {
	path := fmt.Sprintf("/v3/accounts/%s/orders", accountID)

	requestBody := oandaOrder{
		Order: oandaOrderRequest{
			OrderType:  orderType,
			Instrument: instrument,
			Units:      units,
		},
	}

	var orderCreated oandaOrderCreated
	if err := oanda.client.post(path, &requestBody, &orderCreated); err != nil {
		return nil, err
	}
	return &orderCreated, nil
}

func (oanda *oanda) putCloseTrade(accountID string, tradeID string) (*oandaTradeClosed, error) {
	path := fmt.Sprintf("/v3/accounts/%s/trades/%s/close", accountID, tradeID)

	requestBody := oandaCloseTrade{
		Units: "ALL",
	}

	var tradeClosed oandaTradeClosed
	if err := oanda.client.put(path, &requestBody, &tradeClosed); err != nil {
		return nil, err
	}
	return &tradeClosed, nil
}

func (oanda *oanda) makeCandleSpecs(granularity string, instruments ...string) []string {
	specs := make([]string, 0)
	for _, instrument := range instruments {
		spec := fmt.Sprintf("%s:%s:M", instrument, granularity)
		specs = append(specs, spec)
	}
	return specs
}
