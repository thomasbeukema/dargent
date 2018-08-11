package address

import (
    "github.com/tyler-smith/go-bip39"
)

// Generate new seed and mnemonic
func GenerateSeedAndMnemonic(password string) ([]byte, []byte) {
	// TODO: Check for errors

	entropy, _ := bip39.NewEntropy(256)
	s, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed(s, password)

	return entropy, seed
}
