package broker

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// OandaClient is OANDA API client
type OandaClient struct {
	client *apiClient
}

// NewOandaClient is a constructor for Oanda
func NewOandaClient() *OandaClient {
	oandaHost := os.Getenv("OANDA_HOST")
	if oandaHost == "" {
		log.Panicln("OANDA_HOST has not been set")
	}

	oandaToken := os.Getenv("OANDA_TOKEN")
	if oandaToken == "" {
		log.Panicln("OANDA_TOKEN has not been set")
	}

	var header = http.Header{}

	header.Add("Authorization", fmt.Sprintf("Bearer %s", oandaToken))
	header.Add("Connection", "Keep-Alive")
	header.Add("Content-Type", "application/json")
	header.Add("Accept-Datetime-Format", "UNIX")

	return &OandaClient{client: newAPIClient(oandaHost, header)}
}

// Accounts is a method to get accounts
func (oanda *OandaClient) Accounts() (*OandaAccounts, error) {
	var accounts OandaAccounts
	if err := oanda.client.get("/v3/accounts", nil, &accounts); err != nil {
		return nil, err
	}
	return &accounts, nil
}

// Account is a method to get account
func (oanda *OandaClient) Account(accountID string) (*OandaAccount, error) {
	path := fmt.Sprintf("/v3/accounts/%s", accountID)

	var account OandaAccount
	if err := oanda.client.get(path, nil, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// AccountChanges is a method to get account changes
func (oanda *OandaClient) AccountChanges(accountID string, sinceTransactionID int) (*OandaAccountChanges, error) {
	path := fmt.Sprintf("/v3/accounts/%s/changes", accountID)

	query := url.Values{}
	query.Add("sinceTransactionID", oandaInt(sinceTransactionID).String())

	var accountChanges OandaAccountChanges
	if err := oanda.client.get(path, query, &accountChanges); err != nil {
		return nil, err
	}
	return &accountChanges, nil
}

// Candles is a method to get candles
func (oanda *OandaClient) Candles(instrument string, granularity string, count int, from int, to int, includeFirst bool) (*OandaCandles, error) {
	path := fmt.Sprintf("/v3/instruments/%s/candles", instrument)

	query := url.Values{}

	query.Add("price", "MBA")
	query.Add("granularity", granularity)

	if count > 0 {
		query.Add("count", oandaInt(count).String())
	}

	if from > 0 {
		query.Add("from", oandaInt(from).String())
		query.Add("includeFirst", oandaBool(includeFirst).String())
	}

	if to > 0 {
		query.Add("to", oandaInt(to).String())
	}

	var candles OandaCandles
	if err := oanda.client.get(path, query, &candles); err != nil {
		return nil, err
	}
	return &candles, nil
}

// Pricing is a method to get pricing
func (oanda *OandaClient) Pricing(accountID string, instruments []string, since int) (*OandaPrices, error) {
	path := fmt.Sprintf("/v3/accounts/%s/pricing", accountID)

	query := url.Values{}
	query.Add("instruments", strings.Join(instruments, ","))
	query.Add("since", oandaInt(since).String())

	var prices OandaPrices
	if err := oanda.client.get(path, query, &prices); err != nil {
		return nil, err
	}
	return &prices, nil
}

// LatestCandles is a method to get latest candles
func (oanda *OandaClient) LatestCandles(accountID string, specs []string) (*OandaLatestCandles, error) {
	path := fmt.Sprintf("/v3/accounts/%s/candles/latest", accountID)

	query := url.Values{}
	query.Add("candleSpecifications", strings.Join(specs, ","))

	var latestCandles OandaLatestCandles
	if err := oanda.client.get(path, query, &latestCandles); err != nil {
		return nil, err
	}
	return &latestCandles, nil
}

// OpenTrades is a method to get open trades
func (oanda *OandaClient) OpenTrades(accountID string) (*OandaTrades, error) {
	path := fmt.Sprintf("/v3/accounts/%s/openTrades", accountID)

	var trades OandaTrades
	if err := oanda.client.get(path, nil, &trades); err != nil {
		return nil, err
	}
	return &trades, nil
}

// CreateOrder is a method to post order
func (oanda *OandaClient) CreateOrder(accountID string, orderType string, instrument string, units float64) (*OandaOrderCreated, error) {
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

// CloseTrade is a method to put close trade
func (oanda *OandaClient) CloseTrade(accountID string, tradeID int) (*OandaTradeClosed, error) {
	path := fmt.Sprintf("/v3/accounts/%s/trades/%d/close", accountID, tradeID)

	requestBody := oandaCloseTrade{
		Units: "ALL",
	}

	var tradeClosed OandaTradeClosed
	if err := oanda.client.put(path, &requestBody, &tradeClosed); err != nil {
		return nil, err
	}
	return &tradeClosed, nil
}

func (oanda *OandaClient) makeCandleSpecs(granularity string, instruments ...string) []string {
	specs := make([]string, 0)
	for _, instrument := range instruments {
		spec := fmt.Sprintf("%s:%s:MBA", instrument, granularity)
		specs = append(specs, spec)
	}
	return specs
}
