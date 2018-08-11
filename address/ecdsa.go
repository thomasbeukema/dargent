package address

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
)

// Define custom base32 encoding
var waspEncoding = base32.NewEncoding(charset)

// Struct to hold ECC Keys
type ECCKeyPair struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// Generate new KeyPair (new address)
func GenerateECCKeyPair(seed []byte) ECCKeyPair {
	curve := elliptic.P256() // Init the curve we're using
	if seed != nil {
		private, err := ecdsa.GenerateKey(curve, bytes.NewReader(seed)) // Generate the new private key from bip39 seed
		if err != nil {
			panic(1)
		}

		public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...) // Derive public key from private key
		return ECCKeyPair{*private, public}
	} else {
		private, err := ecdsa.GenerateKey(curve, rand.Reader) // Generate address from own seed, no mnemonic phrase tho
		if err != nil {
			panic(1)
		}
		public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...) // Derive public key from private key
		return ECCKeyPair{*private, public}
	}
}

func ECCPubKeyToAddress(pubkey []byte) string {
	pubkey = append(pubkey, byte(0x00)) // Append padding byte to avoid the Base32 '=' padding

	b32Pubkey := waspEncoding.EncodeToString(pubkey)                                 // Generate Base32 of public key
	b32Checksum := waspEncoding.EncodeToString(generateChecksum(HashPubKey(pubkey))) // Get checksum for the public key and generate Base32 of it

	address := "666" + b32Pubkey + b32Checksum + "999" // Pre-/append hardcoded strings to complete address
	return address
}

// Get address from public key
func (kp ECCKeyPair) GetAddress() string {
	return ECCPubKeyToAddress(kp.PublicKey)
}

// Check if an address is actually valid
func ValidateECCAddress(address string) bool {
	address = address[3 : len(address)-3]                                                  // Remove the '666' prefix and '999' appendix
	checksum := address[len(address)-8:]                                                       // Extract already generated checksum
	decodedPubkey, _ := waspEncoding.DecodeString(address[:len(address)-8])                    // Extract en decode public key from address
	targetChecksum := waspEncoding.EncodeToString(generateChecksum(HashPubKey(decodedPubkey))) // Generate new checksum to see if there hasn't been tampered with
	return checksum == targetChecksum
}

// Reverse PubKeyToAddress function
func ECCAddressToPubKey(address string) []byte {
	address = address[3 : len(address)-3]                            // Remove the '666' prefix and '999' appendix
	pubkey, _ := waspEncoding.DecodeString(address[:len(address)-8]) // Extract en decode public key from address
	return pubkey[:len(pubkey)-1]
}

// Sign a transaction with private key
func (kp ECCKeyPair) Sign(hash []byte) string {
	// TODO: Error checking
	r, s, _ := ecdsa.Sign(rand.Reader, &kp.PrivateKey, hash) // Sign the hash
	signature := append(r.Bytes(), s.Bytes()...)             // Append both parts to get 1 signature
	return base64.StdEncoding.EncodeToString(signature)      // Return the base64 encoded signature
}
