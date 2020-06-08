package oanda

// Accounts is a definition for accounts
type Accounts struct {
	Accounts []AccountProperties `json:"accounts"`
}

// AccountProperties is a definition for account properties
type AccountProperties struct {
	ID string `json:"id"`
}
