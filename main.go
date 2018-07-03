package main

import (
	"fmt"
	
	"github.com/tyler-smith/go-bip39"
	"github.com/davecgh/go-spew/spew"
	
	"github.com/thomasbeukema/dargent/transaction"
	"github.com/thomasbeukema/dargent/address"
	"github.com/thomasbeukema/dargent/library"
)

func main() {
	w1 := "table half snack crystal push husband awkward walk social educate general report shield asset border hole world dream pencil occur visual spy absorb shell"
	entropy, _ := bip39.NewEntropy(256)
	w2, _ := bip39.NewMnemonic(entropy)
	
	seed1 := bip39.NewSeed(w1, "boterham1234")
	seed2 := bip39.NewSeed(w2, "blabla")
	
	kp1 := address.GenerateKeyPair(seed1)
	kp2 := address.GenerateKeyPair(seed2)
	
	a1 := string(kp1.GetAddress())
	a2 := string(kp2.GetAddress())
	
	fmt.Printf("Address 1: %s\n", a1)
	fmt.Printf("Address 2: %s\n", a2)
	
	tx, _ := transaction.NewSendTransaction(kp1, "", a2, 666, transaction.NativeCurrency())	
	tx.Signature = kp1.SignTx([]byte(tx.Hash))
	
	spew.Dump(tx)
	
	shelf := library.NewShelf(kp1)
	shelf.AddToLibrary()
}

		
		
		
		
		
		