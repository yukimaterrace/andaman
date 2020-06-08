package oanda

import "testing"

var accountID string

func TestGetAccounts(t *testing.T) {
	accounts, err := GetAccounts()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", accounts)

	accountID = accounts.Accounts[0].ID
}
