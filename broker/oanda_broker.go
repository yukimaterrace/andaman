package broker

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"yukimaterrace/andaman/util"
)

// OandaBroker is OANDA broker
type OandaBroker struct {
	client    *apiClient
	accountID string
}

// NewOandaBroker is a constructor for Oanda
func NewOandaBroker() *OandaBroker {
	oandaHost := util.GetEnv("OANDA_HOST")
	oandaToken := util.GetEnv("OANDA_TOKEN")

	var header = http.Header{}

	header.Add("Authorization", fmt.Sprintf("Bearer %s", oandaToken))
	header.Add("Connection", "Keep-Alive")
	header.Add("Content-Type", "application/json")
	header.Add("Accept-Datetime-Format", "UNIX")

	oanda := &OandaBroker{
		client: newAPIClient(oandaHost, header),
	}

	resp, err := oanda.Accounts()
	if err != nil {
		panic(err)
	}
	oanda.accountID = resp.Accounts[0].ID

	return oanda
}

// Accounts is a method to get accounts
func (oanda *OandaBroker) Accounts() (*OandaAccounts, error) {
	var accounts OandaAccounts
	if err := oanda.client.get("/v3/accounts", nil, &accounts); err != nil {
		return nil, err
	}
	return &accounts, nil
}

// Account is a method to get account
func (oanda *OandaBroker) Account() (*OandaAccount, error) {
	path := fmt.Sprintf("/v3/accounts/%s", oanda.accountID)

	var account OandaAccount
	if err := oanda.client.get(path, nil, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// AccountChanges is a method to get account changes
func (oanda *OandaBroker) AccountChanges(sinceTransactionID int) (*OandaAccountChanges, error) {
	path := fmt.Sprintf("/v3/accounts/%s/changes", oanda.accountID)

	query := url.Values{}
	query.Add("sinceTransactionID", oandaInt(sinceTransactionID).String())

	var accountChanges OandaAccountChanges
	if err := oanda.client.get(path, query, &accountChanges); err != nil {
		return nil, err
	}
	return &accountChanges, nil
}

// Candles is a method to get candles
func (oanda *OandaBroker) Candles(instrument string, granularity string, count int, from int, to int, includeFirst bool) (*OandaCandles, error) {
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
func (oanda *OandaBroker) Pricing(instruments []string, since int) (*OandaPrices, error) {
	path := fmt.Sprintf("/v3/accounts/%s/pricing", oanda.accountID)

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
func (oanda *OandaBroker) LatestCandles(specs []string) (*OandaLatestCandles, error) {
	path := fmt.Sprintf("/v3/accounts/%s/candles/latest", oanda.accountID)

	query := url.Values{}
	query.Add("candleSpecifications", strings.Join(specs, ","))

	var latestCandles OandaLatestCandles
	if err := oanda.client.get(path, query, &latestCandles); err != nil {
		return nil, err
	}
	return &latestCandles, nil
}

// OpenTrades is a method to get open trades
func (oanda *OandaBroker) OpenTrades() (*OandaTrades, error) {
	path := fmt.Sprintf("/v3/accounts/%s/openTrades", oanda.accountID)

	var trades OandaTrades
	if err := oanda.client.get(path, nil, &trades); err != nil {
		return nil, err
	}
	return &trades, nil
}

// CreateOrder is a method to post order
func (oanda *OandaBroker) CreateOrder(orderType string, instrument string, units float64) (*OandaOrderCreated, error) {
	path := fmt.Sprintf("/v3/accounts/%s/orders", oanda.accountID)

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
func (oanda *OandaBroker) CloseTrade(tradeID int) (*OandaTradeClosed, error) {
	path := fmt.Sprintf("/v3/accounts/%s/trades/%d/close", oanda.accountID, tradeID)

	requestBody := oandaCloseTrade{
		Units: "ALL",
	}

	var tradeClosed OandaTradeClosed
	if err := oanda.client.put(path, &requestBody, &tradeClosed); err != nil {
		return nil, err
	}
	return &tradeClosed, nil
}

func (oanda *OandaBroker) makeCandleSpecs(granularity string, instruments ...string) []string {
	specs := make([]string, 0)
	for _, instrument := range instruments {
		spec := fmt.Sprintf("%s:%s:MBA", instrument, granularity)
		specs = append(specs, spec)
	}
	return specs
}
