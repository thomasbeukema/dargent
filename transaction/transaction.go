package transaction

import (
	"fmt"
	"encoding/json"
	"encoding/base64"
	"crypto/sha256"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"time"
	"strconv"
	
	"github.com/thomasbeukema/dargent/address"
)

type transactionType int
const (
	Send transactionType = iota
	Claim
	Create
	Trust
)

type SendTransaction struct {
	Hash			string				`json:"h"`
	PreviousHash	string				`json:"p"`
	Action			transactionType		`json:"a"`
	Balance			uint64				`json:"b"`
	Currency		Currency			`json:"c,omitempty"`
	Origin			string				`json:"o"`
	Destination		string				`json:"d,omitempty"`
	Signature		string				`json:"s"`
}

type ClaimTransaction struct {
	Hash			string				`json:"h"`
	PreviousHash	string				`json:"p"`
	Action			transactionType		`json:"a"`
	Origin			string				`json:"o"`
	Destination		string				`json:"d"`
	Signature		string				`json:"s"`
}

type CreateTransaction struct {
	Hash			string				`json:"h"`
	Action			transactionType		`json:"a"`
	Balance			uint64				`json:"b,omitempty"`
	Currency		Currency			`json:"c,omitempty"`
	Origin			string				`json:"o"`
	Signature		string				`json:"s"`
}

type TrustTransaction struct {
	Hash			string				`json:"h"`
	PreviousHash	string				`json:"p"`
	Action			transactionType		`json:"a"`
	Origin			string				`json:"o"`
	Destination		string				`json:"d"`
	Expiration		string				`json:"e"`
	Signature		string				`json:"s"`
}

func NewSendTransaction(account address.KeyPair, ph string, destination string, amount uint64, c Currency) (SendTransaction, error) {
	tx := SendTransaction{
		Hash: "",
		PreviousHash: ph,
		Action: Send,
		Balance: amount,
		Currency: c,
		Origin: account.GetAddress(),
		Destination: destination,
		Signature: "",
	}
	
	tx.Hash,_ = tx.GenerateHash()
	
	return tx, nil
}

func (tx *SendTransaction) GenerateHash() (string, error) {
	minTx := SendTransaction{
		Hash: "",
		PreviousHash: tx.PreviousHash,
		Action: Send,
		Balance: tx.Balance,
		Currency: tx.Currency,
		Origin: tx.Origin,
		Destination: tx.Destination,
		Signature: "",
	}
	
	txJson, err := json.Marshal(minTx)
	if err != nil {
		return "", err	
	}
	
	hash := sha256.Sum256(txJson)
	
	return fmt.Sprintf("%x", hash[:]), nil
}
	
func (tx *SendTransaction) Verify() bool {
	// TODO: Check previous hash
	if tx.Balance < 0 {
		return false
	}	
	if address.ValidateAddress(tx.Origin) != true {
		return false
	}
	if address.ValidateAddress(tx.Destination) != true {
		return false
	}
	if h, _ := tx.GenerateHash(); tx.Hash != h {
		return false
	}
	
	pubkey := address.AddressToPubKey(tx.Origin)
	curve := elliptic.P256()
	
	signatureBytes, _ := base64.StdEncoding.DecodeString(tx.Signature)
	signatureLength := len(signatureBytes)
	
	r := big.Int{}
	s := big.Int{}
	
	r.SetBytes(signatureBytes[:(signatureLength/2)])
	s.SetBytes(signatureBytes[(signatureLength/2):])
	
	x := big.Int{}
	y := big.Int{}
	
	keyLength := len(pubkey)
	x.SetBytes(pubkey[:(keyLength/2)])
	y.SetBytes(pubkey[(keyLength/2):])
	
	rawPubKey := ecdsa.PublicKey{curve, &x, &y}
	if ecdsa.Verify(&rawPubKey, []byte(tx.Hash), &r, &s) == false {
		return false
	}
	
	return true
}

func NewClaimTransaction(account address.KeyPair, ph string, txId string) (ClaimTransaction, error) {
	tx := ClaimTransaction{
		Hash: "",
		PreviousHash: ph,
		Action: Claim,
		Origin: txId,
		Destination: account.GetAddress(),
		Signature: "",
	}
	
	tx.Hash,_ = tx.GenerateHash()
	
	return tx, nil
}

func (tx *ClaimTransaction) GenerateHash() (string, error) {
	minTx := ClaimTransaction{
		Hash: "",
		PreviousHash: tx.PreviousHash,
		Action: Send,
		Origin: tx.Origin,
		Destination: tx.Destination,
		Signature: "",
	}
	
	txJson, err := json.Marshal(minTx)
	if err != nil {
		return "", err	
	}
	
	hash := sha256.Sum256(txJson)
	
	return fmt.Sprintf("%x", hash[:]), nil
}
	
func (tx *ClaimTransaction) Verify() bool {
	// TODO: Check send hash
	if address.ValidateAddress(tx.Destination) != true {
		return false
	}
	if h, _ := tx.GenerateHash(); tx.Hash != h {
		return false
	}
	
	pubkey := address.AddressToPubKey(tx.Origin)
	curve := elliptic.P256()
	
	signatureBytes, _ := base64.StdEncoding.DecodeString(tx.Signature)
	signatureLength := len(signatureBytes)
	
	r := big.Int{}
	s := big.Int{}
	
	r.SetBytes(signatureBytes[:(signatureLength/2)])
	s.SetBytes(signatureBytes[(signatureLength/2):])
	
	x := big.Int{}
	y := big.Int{}
	
	keyLength := len(pubkey)
	x.SetBytes(pubkey[:(keyLength/2)])
	y.SetBytes(pubkey[(keyLength/2):])
	
	rawPubKey := ecdsa.PublicKey{curve, &x, &y}
	if ecdsa.Verify(&rawPubKey, []byte(tx.Hash), &r, &s) == false {
		return false
	}
	
	return true
}

func NewCreateTransaction(account string) (CreateTransaction, error) {
	tx := CreateTransaction{
		Hash: "",
		Action: Create,
		Currency: NativeCurrency(),
		Balance: 0,
		Origin: account,
		Signature: "",
	}
	
	tx.Hash,_ = tx.GenerateHash()
	
	return tx, nil
}

func NewCreateTokenTransaction(account address.KeyPair, c Currency, amount uint64) (CreateTransaction, error) {
	tx := CreateTransaction{
		Hash: "",
		Action: Create,
		Currency: c,
		Balance: amount,
		Origin: account.GetAddress(),
		Signature: "",
	}
	
	tx.Hash,_ = tx.GenerateHash()
	
	return tx, nil
}

func (tx *CreateTransaction) GenerateHash() (string, error) {
	minTx := SendTransaction{
		Hash: "",
		Action: Create,
		Balance: tx.Balance,
		Currency: tx.Currency,
		Origin: tx.Origin,
		Signature: "",
	}
	
	txJson, err := json.Marshal(minTx)
	if err != nil {
		return "", err	
	}
	
	hash := sha256.Sum256(txJson)
	
	return fmt.Sprintf("%x", hash[:]), nil
}
	
func (tx *CreateTransaction) Verify() bool {
	// TODO: Check previous hash
	if tx.Balance != 0 && tx.Currency == NativeCurrency() {
		return false
	}	
	if address.ValidateAddress(tx.Origin) != true {
		return false
	}
	if h, _ := tx.GenerateHash(); tx.Hash != h {
		return false
	}
	
	pubkey := address.AddressToPubKey(tx.Origin)
	curve := elliptic.P256()
	
	signatureBytes, _ := base64.StdEncoding.DecodeString(tx.Signature)
	signatureLength := len(signatureBytes)
	
	r := big.Int{}
	s := big.Int{}
	
	r.SetBytes(signatureBytes[:(signatureLength/2)])
	s.SetBytes(signatureBytes[(signatureLength/2):])
	
	x := big.Int{}
	y := big.Int{}
	
	keyLength := len(pubkey)
	x.SetBytes(pubkey[:(keyLength/2)])
	y.SetBytes(pubkey[(keyLength/2):])
	
	rawPubKey := ecdsa.PublicKey{curve, &x, &y}
	if ecdsa.Verify(&rawPubKey, []byte(tx.Hash), &r, &s) == false {
		return false
	}
	
	return true
}

func NewTrustTransaction(account address.KeyPair, destination string, expiration string) (TrustTransaction, error) {
	tx := TrustTransaction{
		Hash: "",
		PreviousHash: "",
		Action: Trust,
		Origin: account.GetAddress(),
		Destination: destination,
		Expiration: expiration,
		Signature: "",
	}
	
	tx.Hash,_ = tx.GenerateHash()
	tx.PreviousHash,_ = tx.GenerateHash()
	
	return tx, nil
}

func (tx *TrustTransaction) GenerateHash() (string, error) {
	minTx := TrustTransaction{
		Hash: "",
		PreviousHash: tx.PreviousHash,
		Action: Trust,
		Origin: tx.Origin,
		Destination: tx.Destination,
		Expiration: tx.Origin,
		Signature: "",
	}
	
	txJson, err := json.Marshal(minTx)
	if err != nil {
		return "", err	
	}
	
	hash := sha256.Sum256(txJson)
	
	return fmt.Sprintf("%x", hash[:]), nil
}
	
func (tx *TrustTransaction) Verify() bool {
	// TODO: Check previous hash
	if address.ValidateAddress(tx.Origin) != true {
		return false
	}
	if address.ValidateAddress(tx.Destination) != true {
		return false
	}
	
	currentTime := time.Now().UnixNano()
	expiring, _ := strconv.ParseInt(tx.Expiration, 10, 64)
	// TODO: Check for error
	
	if currentTime > expiring {
		return false
	}
	
	if h, _ := tx.GenerateHash(); tx.Hash != h {
		return false
	}
	
	pubkey := address.AddressToPubKey(tx.Origin)
	curve := elliptic.P256()
	
	signatureBytes, _ := base64.StdEncoding.DecodeString(tx.Signature)
	signatureLength := len(signatureBytes)
	
	r := big.Int{}
	s := big.Int{}
	
	r.SetBytes(signatureBytes[:(signatureLength/2)])
	s.SetBytes(signatureBytes[(signatureLength/2):])
	
	x := big.Int{}
	y := big.Int{}
	
	keyLength := len(pubkey)
	x.SetBytes(pubkey[:(keyLength/2)])
	y.SetBytes(pubkey[(keyLength/2):])
	
	rawPubKey := ecdsa.PublicKey{curve, &x, &y}
	if ecdsa.Verify(&rawPubKey, []byte(tx.Hash), &r, &s) == false {
		return false
	}
	
	return true
}



























