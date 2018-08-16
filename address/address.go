package address

import (
    "crypto/sha256"
)

const (
	// Length of the checksum attached
	checksumLength = 5
)

// padding for addresses
var ecdsaPadding []byte = []byte{0xAA, 0xBB, 0xCC}
var sphincsPadding []byte = []byte{0x00, 0x11, 0x22}

type AccountType int
const (
    ECC AccountType = iota // 0
    SPHINCS // 1
    UNKNOWN
)

// TODO: Complete function
func ValidateAddress(address string) bool {

    if len(address) == 61 && address[:3] == "666" && address[len(address)-3:] == "999" { // ECDSA
        validateECDSAAddress(address)
    } else if len(address) == 61 && address[:3] == "999" && address[len(address)-3:] == "666" { // SPHINCS
        validateSPHINCSAddress(address)
    } else {
        return false
    }

    return true
}

func TypeOfAddress(address string) AccountType {
    if len(address) == 61 && address[:3] == "666" && address[len(address)-3:] == "999" { // ECDSA
        return ECC
    } else if  len(address) == 61 && address[:3] == "999" && address[len(address)-3:] == "666" { // SPHINCS
        return SPHINCS
    } else {
        return UNKNOWN
    }
}

func TypeOfPublicKey(publicKey []byte) AccountType {
    if len(publicKey) == 64 { // ECDSA
        return ECC
    } else if len(publicKey) == 1056 { // SPHINCS
        return SPHINCS
    } else { // unknown or invalid
        return UNKNOWN
    }
}

func PubKeyToAddress(pubkey []byte) string {
    t := TypeOfPublicKey(pubkey)

    if t == ECC {
        return ECCPubKeyToAddress(pubkey)
    } else if t == SPHINCS {
        return SPHINCSPubKeyToAddress(pubkey)
    } else {
        return ""
    }
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
