package transaction

import (

)

type Currency struct {
	Name	string			`json:"n"`
	Ticker	string			`json:"t"`
	Owner	string			`json:"o"`
}

func NewCurrency(name, ticker, owner string) Currency {
	return Currency{name, ticker, owner}
}

func (c *Currency) GetName() string {
	return c.Name
}

func (c *Currency) GetOwnerAddress() string {
	return c.Owner
}

func NativeCurrency() Currency {
	return Currency{"Argent", "ART", ""}
}