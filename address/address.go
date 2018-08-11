package address

import (
    "crypto/sha256"
)

const (
	// Define custom charset for address
	charset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

	// Length of the checksum attached
	checksumLength = 5
)

// padding for SPHINCS address
var sphincsPadding []byte = []byte{0x00, 0x11, 0x22}

// TODO: Complete function
func ValidateAddress(_ string) bool {
    return true
}

// Generate hash for validation
func generateChecksum(payload []byte) []byte {
	hash := sha256.Sum256(payload)
	finalHash := sha256.Sum256(hash[:])

	return finalHash[:checksumLength]
}

// Just generate the SHA-256 of a pub key
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	return publicSHA256[:]
}
