package address

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/rand"
	
	"github.com/golang/crypto/blake2b"
	_ "github.com/tyler-smith/go-bip39"
	_ "golang.org/x/crypto/bcrypt"
)

const (
	charset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	
	version = byte(0x21)
	checksumLength = 5
)

var	waspEncoding = base32.NewEncoding(charset)

type KeyPair struct {
	PrivateKey	ecdsa.PrivateKey
	PublicKey	[]byte
}

func GenerateKeyPair(seed []byte) KeyPair {
	curve := elliptic.P256()
	if seed != nil {
		private, err := ecdsa.GenerateKey(curve, bytes.NewReader(seed))
		if err != nil {
			panic(1)	
		}
		
		public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
		return KeyPair{*private, public}
	} else {
		private, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			panic(1)	
		}
		public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
		return KeyPair{*private, public}
	}
}

func (kp KeyPair) GetAddress() string {
	pubkey := kp.PublicKey
	pubkey = append(pubkey, byte(0x00))
	
	b32Pubkey := waspEncoding.EncodeToString(pubkey)
	b32Checksum := waspEncoding.EncodeToString(generateChecksum(HashPubKey(pubkey)))
	
	address := "666" + b32Pubkey + b32Checksum + "999"
	return address
}

func generateChecksum(payload []byte) []byte {
	hash, _ := blake2b.New(11, nil)
	hash.Write(payload)
	finalHash := sha256.Sum256(hash.Sum(nil))

	return finalHash[:checksumLength]
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	return publicSHA256[:]
}

func ValidateAddress(address string) bool {
	address = address[3:len(address)-3]
	checksum := address[len(address)-8:]
	decodedPubkey, _ := waspEncoding.DecodeString(address[:len(address)-8])
	targetChecksum := waspEncoding.EncodeToString(generateChecksum(HashPubKey(decodedPubkey)))
	return checksum == targetChecksum
}

func AddressToPubKey(address string) []byte {
	address = address[3:len(address)-3]
	pubkey, _ := waspEncoding.DecodeString(address[:len(address)-8])
	return pubkey[:len(pubkey)-1]
}

func (kp KeyPair) SignTx(hash []byte) string {
	// TODO: Error checking
	r,s,_ := ecdsa.Sign(rand.Reader, &kp.PrivateKey, hash)
	signature := append(r.Bytes(), s.Bytes()...)
	return base64.StdEncoding.EncodeToString(signature)
}



















