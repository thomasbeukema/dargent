package address

import (
    "github.com/NebulousLabs/entropy-mnemonics"
)

// Generate new seed and mnemonic
func getMnemonic(key []byte) string {
    phrase, _ := mnemonics.ToPhrase(key, mnemonics.English)
    return phrase.String()
}

func MnemonicToEntropy(phrase string) [32]byte {
    ent, _ := mnemonics.FromString(phrase, mnemonics.English)

    var ent32 [32]byte

    for i := range ent {
        ent32[i] = ent[i]
    }

    return ent32
}
