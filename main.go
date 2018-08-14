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

	kp3 := address.GenerateECCKeyPair(nil)
	kp4 := address.GenerateSPHINCSKeyPair(nil)

	a1 := string(kp1.GetAddress())
	a2 := string(kp2.GetAddress())

	fmt.Println(len(kp1.PublicKey))
	fmt.Println(len(kp2.PublicKey))

	fmt.Printf("Address 1: %s\n", a1)
	fmt.Printf("Valid: %v\n", address.ValidateAddress(a1))
	fmt.Printf("Address 2: %s\n", a2)
	fmt.Printf("Valid: %v\n", address.ValidateAddress(a2))

	msg := "Hey, I'm authentic."

	sig1 := kp1.Sign([]byte(msg))
	sig2 := kp2.Sign([]byte(msg))

	valid1 := address.ValidateECCSignature(sig1, msg, kp3.PublicKey)
	valid2 := address.ValidateSPHINCSSignature(sig2, msg, kp4.PublicKey)

	fmt.Println(valid1)
	fmt.Println(valid2)
}
