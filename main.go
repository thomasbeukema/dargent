package main

import (
	"fmt"

	"github.com/thomasbeukema/dargent/account"
	"github.com/thomasbeukema/dargent/address"

	_ "github.com/davecgh/go-spew/spew"
)

func main() {

	m := "mundane atom sack seventh goldfish cottage vacation lemon pram eclipse syndrome return firm after arises bobsled tadpoles shipped tugs second sipped uphill afraid ardent"
	e := address.MnemonicToEntropy(m)

	kp1 := address.GenerateECCKeyPair(e[:])
	kp2 := address.GenerateSPHINCSKeyPair(e[:])

	pk2 := kp2.PublicKey[:]

	a1 := string(kp1.GetAddress())
	a2 := string(kp2.GetAddress())

	fmt.Printf("Address 1: %s\n", a1)
	fmt.Printf("Mnemonic 1: %s\n", kp1.Mnemonic())
	fmt.Printf("Address 2: %s\n", a2)
	fmt.Printf("Mnemonic 2: %s\n", kp2.Mnemonic())

	tx1, _ := account.NewCreateTransaction(kp1.PublicKey)
	tx2, _ := account.NewCreateTransaction(pk2)

	acc1 := account.OpenAccount(a1, kp1.PublicKey)
	acc1.AddTransaction(tx1)

	acc2 := account.OpenAccount(a2, pk2)
	acc2.AddTransaction(tx2)

	led1 := acc1.OpenLedger(account.NativeCurrency().Ticker)
	led2 := acc2.OpenLedger(account.NativeCurrency().Ticker)

	s1 := kp1.Sign([]byte(led1.Hash))
	s2 := kp2.Sign([]byte(led2.Hash))

	fmt.Println(led1.UpdateSignature(s1, &acc1))
	fmt.Println(led2.UpdateSignature(s2, &acc2))
}
