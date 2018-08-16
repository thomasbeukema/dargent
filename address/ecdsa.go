package address

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"math/big"

	"github.com/mr-tron/base58/base58"
	"github.com/thomasbeukema/fastrand"
)

// Struct to hold ECC Keys
type ECCKeyPair struct {
	PrivateKey	ecdsa.PrivateKey
	PublicKey	[]byte
	Entropy		*[32]byte
}

// Generate new KeyPair (new address)
func GenerateECCKeyPair(ent []byte) ECCKeyPair {

	fastrand.New()
	var entropy [32]byte

	if ent != nil {
		for i := range ent {
			entropy[i] = ent[i]
		}
	} else {
		entropy = fastrand.GetEntropy()
	}

	curve := elliptic.P256() // Init the curve we're using
	private, err := ecdsa.GenerateKey(curve, fastrand.Reader) // Generate address from own seed
	if err != nil {
		panic(1)
	}
	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...) // Derive public key from private key
	return ECCKeyPair{*private, public, &entropy}
}

func ECCPubKeyToAddress(pubkey []byte) string {
    pubkey = append(HashPubKey(pubkey), ecdsaPadding...)

    b32Pubkey := base58.Encode(pubkey)                                 // Generate Base58 of public key
    b32Checksum := base58.Encode(generateChecksum(HashPubKey(pubkey))) // Get checksum for the public key and generate Base32 of it

    return "666" + b32Pubkey + b32Checksum + "999"
}

func validateECDSAAddress(address string) bool {

    if address[:3] != "666" || address[len(address)-3:] != "999" {
        return false
    }

    address = address[3:len(address)-3] // Strip '666' & '999'

    checksum := address[len(address)-8:]
    generatedChecksum := base58.Encode(generateChecksum([]byte(address[:len(address)-8])))

    return checksum == generatedChecksum
}

// Get address from public key
func (kp ECCKeyPair) GetAddress() string {
	return ECCPubKeyToAddress(kp.PublicKey)
}

func (kp *ECCKeyPair) Mnemonic() string {
	return getMnemonic(kp.Entropy[:])
}

// Check if an address is actually valid
func ValidateECCAddress(address string) bool {
	address = address[3 : len(address)-3]                                                  // Remove the '666' prefix and '999' appendix
	checksum := address[len(address)-8:]                                                       // Extract already generated checksum
	decodedPubkey, _ := base58.Decode(address[:len(address)-8])                    // Extract en decode public key from address
	targetChecksum := base58.Encode(generateChecksum(HashPubKey(decodedPubkey))) // Generate new checksum to see if there hasn't been tampered with
	return checksum == targetChecksum
}

// Sign with private key
func (kp ECCKeyPair) Sign(hash []byte) string {
	// TODO: Error checking
	r, s, _ := ecdsa.Sign(rand.Reader, &kp.PrivateKey, hash) // Sign the hash
	signature := append(r.Bytes(), s.Bytes()...)             // Append both parts to get 1 signature
	return base64.StdEncoding.EncodeToString(signature)      // Return the base64 encoded signature
}

func ValidateECCSignature(sig string, hash []byte, pubkey []byte) bool {
	curve := elliptic.P256() // Init curve

	signatureBytes, _ := base64.StdEncoding.DecodeString(sig) // Extract signature
	signatureLength := len(signatureBytes)

	r := big.Int{} // Parse signature in 2 parts
	s := big.Int{}

	r.SetBytes(signatureBytes[:(signatureLength/2)])
	s.SetBytes(signatureBytes[(signatureLength/2):])

	x := big.Int{} // Parse public key in 2 parts
	y := big.Int{}

	keyLength := len(pubkey)
	x.SetBytes(pubkey[:(keyLength/2)])
	y.SetBytes(pubkey[(keyLength/2):])

	rawPubKey := ecdsa.PublicKey{curve, &x, &y} // Init public key for verification
	if ecdsa.Verify(&rawPubKey, hash, &r, &s) == false {
		return false
	}

	return true
}
