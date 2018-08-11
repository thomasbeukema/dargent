package main

import (
	"fmt"

	_ "github.com/thomasbeukema/dargent/account"
	"github.com/thomasbeukema/dargent/address"

	_ "github.com/davecgh/go-spew/spew"
)

func main() {
	kp1 := address.GenerateECCKeyPair(nil)
	kp2 := address.GenerateSPHINCSKeyPair(nil)

	a1 := string(kp1.GetAddress())
	a2 := string(kp2.GetAddress())

	fmt.Printf("Address 1: %s\n", a1)
	fmt.Printf("Valid: %v\n", address.ValidateAddress(a1))
	fmt.Printf("Address 2: %s\n", a2)
	fmt.Printf("Valid: %v\n", address.ValidateAddress(a2))

	/*//shelf := transaction.NewShelf(kp1)
	shelf := transaction.OpenShelf(kp1.GetAddress())
	shelf.UpdateLibrary()

	newCurrency := transaction.NewCurrency("SATAN", "STN", a1)

	tx, _ := transaction.NewCreateTokenTransaction(kp1, newCurrency, 666)
	tx.Signature = kp1.SignTx([]byte(tx.Hash))
	shelf.ShelveTx(tx)

	//spew.Dump(tx)

	shelf2 := transaction.NewShelf(kp2)
	//shelf2 := transaction.OpenShelf(kp2)
	shelf2.UpdateLibrary()

	tx2, _ := transaction.NewSendTransaction(kp2, shelf2.Newest().Hash, a1, 100, transaction.NativeCurrency())
	tx2.Signature = kp2.SignTx([]byte(tx2.Hash))
	spew.Dump(shelf2)
	shelf2.ShelveTx(tx2)*/
}
