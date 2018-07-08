package address

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"

	"github.com/golang/crypto/blake2b"
	"github.com/tyler-smith/go-bip39"
	_ "golang.org/x/crypto/bcrypt"
)

const (
	// Define custom charset for address
	charset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

	// Length of the checksum attached
	checksumLength = 5
)

// Define custom base32 encoding
var waspEncoding = base32.NewEncoding(charset)

// Struct to hold ECC Keys
type KeyPair struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// Generate new seed and mnemonic
func GenerateSeedAndMnemonic(password string) ([]byte, []byte) {
	// TODO: Check for errors

	entropy, _ := bip39.NewEntropy(256)
	s, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed(s, password)

	return entropy, seed
}

// Generate new KeyPair (new address)
func GenerateKeyPair(seed []byte) KeyPair {
	curve := elliptic.P256() // Init the curve we're using
	if seed != nil {
		private, err := ecdsa.GenerateKey(curve, bytes.NewReader(seed)) // Generate the new private key from bip39 seed
		if err != nil {
			panic(1)
		}

		public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...) // Derive public key from private key
		return KeyPair{*private, public}
	} else {
		private, err := ecdsa.GenerateKey(curve, rand.Reader) // Generate address from own seed, no mnemonic phrase tho
		if err != nil {
			panic(1)
		}
		public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...) // Derive public key from private key
		return KeyPair{*private, public}
	}
}

// Get address from public key
func (kp KeyPair) GetAddress() string {
	pubkey := kp.PublicKey
	pubkey = append(pubkey, byte(0x00)) // Append padding byte to avoid the Base32 '=' padding

	b32Pubkey := waspEncoding.EncodeToString(pubkey)                                 // Generate Base32 of public key
	b32Checksum := waspEncoding.EncodeToString(generateChecksum(HashPubKey(pubkey))) // Get checksum for the public key and generate Base32 of it

	address := "666" + b32Pubkey + b32Checksum + "999" // Pre-/append hardcoded strings to complete address
	return address
}

// Generate hash to slow things down a bit
func generateChecksum(payload []byte) []byte {
	hash, _ := blake2b.New(11, nil)
	hash.Write(payload)
	finalHash := sha256.Sum256(hash.Sum(nil))

	return finalHash[:checksumLength]
}

// Just generate the SHA-256 of a pub key
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	return publicSHA256[:]
}

// Check if an address is actually valid
func ValidateAddress(address string) bool {
	address = address[3 : len(address)-3]                                                      // Remove the '666' prefix and '999' appendix
	checksum := address[len(address)-8:]                                                       // Extract already generated checksum
	decodedPubkey, _ := waspEncoding.DecodeString(address[:len(address)-8])                    // Extract en decode public key from address
	targetChecksum := waspEncoding.EncodeToString(generateChecksum(HashPubKey(decodedPubkey))) // Generate new checksum to see if there hasn't been tampered with
	return checksum == targetChecksum
}

// Reverse PubKeyToAddress function
func AddressToPubKey(address string) []byte {
	address = address[3 : len(address)-3]                            // Remove the '666' prefix and '999' appendix
	pubkey, _ := waspEncoding.DecodeString(address[:len(address)-8]) // Extract en decode public key from address
	return pubkey[:len(pubkey)-1]
}

// Sign a transaction with private key
func (kp KeyPair) SignTx(hash []byte) string {
	// TODO: Error checking
	r, s, _ := ecdsa.Sign(rand.Reader, &kp.PrivateKey, hash) // Sign the hash
	signature := append(r.Bytes(), s.Bytes()...)             // Append both parts to get 1 signature
	return base64.StdEncoding.EncodeToString(signature)      // Return the base64 encoded signature
}
