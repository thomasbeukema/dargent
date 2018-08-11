package account

type Currency struct {
	Name	string			`json:"n"`
	Ticker	string			`json:"t"`
	Owner	string			`json:"o"`
}

// Create new currency
func NewCurrency(name, ticker, owner string) Currency {
	return Currency{name, ticker, owner}
}

// Get name of currency
func (c *Currency) GetName() string {
	return c.Name
}

// Get owner address of currency
func (c *Currency) GetOwnerAddress() string {
	return c.Owner
}

// Return NativeCurrency 'ART'
func NativeCurrency() Currency {
	return Currency{"Argent", "ART", ""}
}
