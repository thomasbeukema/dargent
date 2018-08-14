package main

import (
	"fmt"

	"github.com/thomasbeukema/dargent/account"
	"github.com/thomasbeukema/dargent/address"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	kp1 := address.GenerateECCKeyPair(nil)
	kp2 := address.GenerateSPHINCSKeyPair(nil)

	pk2 := kp2.PublicKey[:]

	a1 := string(kp1.GetAddress())
	a2 := string(kp2.GetAddress())

	fmt.Printf("Address 1: %s\n", a1)
	fmt.Printf("Address 2: %s\n", a2)

	tx1, _ := account.NewCreateTransaction(kp1.PublicKey)
	tx2, _ := account.NewCreateTransaction(pk2)

	spew.Dump(tx1)

	acc1 := account.OpenAccount(a1, kp1.PublicKey)
	acc1.AddTransaction(tx1)

	acc2 := account.OpenAccount(a2, pk2)
	acc2.AddTransaction(tx2)

	a3 := "66623i6TYjH7YM1jnLKmTG9emhsEmy1bf63aw6kR95pUR1FYRUPJRtuqnH999"

	pub := account.GetPublicKeyFromAddress(a3)
	spew.Dump(pub)
}
