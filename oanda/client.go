package oanda

import (
	"fmt"
	"net/http"
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
