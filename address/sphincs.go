package address

import (
    "bytes"
    "crypto/rand"

    "github.com/Yawning/sphincs256"
)

type SPHINCSKeyPair struct {
	PrivateKey	*[sphincs256.PrivateKeySize]byte
	PublicKey	*[sphincs256.PublicKeySize]byte
}

func GenerateSPHINCSKeyPair(seed []byte) SPHINCSKeyPair {
	if seed != nil { // Seed provided
		pub, priv, err := sphincs256.GenerateKey(bytes.NewBuffer(seed))
		if err != nil {
			panic(err)
		}

		return SPHINCSKeyPair{priv, pub}
	} else { // Seed not provided, generate new one
		pub, priv, err := sphincs256.GenerateKey(rand.Reader)
		if err != nil {
			panic(1)
		}

		return SPHINCSKeyPair{priv, pub}
	}
}

func SPHINCSPubKeyToAddress(pubkey []byte) string {
    pubkey = append(HashPubKey(pubkey), sphincsPadding...)

    b32Pubkey := waspEncoding.EncodeToString(pubkey)                                 // Generate Base32 of public key
    b32Checksum := waspEncoding.EncodeToString(generateChecksum(HashPubKey(pubkey))) // Get checksum for the public key and generate Base32 of it

    return b32Pubkey + b32Checksum
}

func (kp SPHINCSKeyPair) GetAddress() string {
    return SPHINCSPubKeyToAddress((*kp.PublicKey)[:])
}
