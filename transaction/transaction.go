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
	"errors"
	
	"github.com/thomasbeukema/dargent/address"
)

// Define possible transaction types
type transactionType int
const (
	Send transactionType = iota // 0
	Claim // 1
	Create // 2
	Trust // 3
)

// Define struct for transaction structure
type Transaction struct {
	Hash			string				`json:"h"`				// Hash for authenticity and txId
	PreviousHash	string				`json:"p,omitempty"`	// Hash of the previous tx
	Action			transactionType		`json:"a"`				// Type of transaction [Send, Claim, Create, Trust]
	Balance			uint64				`json:"b,omitempty"`	// Balance of the address; balance, not the tx amount
	Currency		Currency			`json:"c,omitempty"`	// Currency of the tx
	Origin			string				`json:"o"`				// Sender / src tx for claim tx
	Destination		string				`json:"d,omitempty"`	// Receiver
	Expiration		string				`json:"e,omitempty"`	// Expiration for trust certificates
	Signature		string				`json:"s"`				// Signature to ensure the owner produced this
}

// Generate hash for a transaction to ensure the authenticity of the contents of the transaction
func (tx *Transaction) GenerateHash() (string, error) {
	switch tx.Action { // Every type of transaction has a slightly different way of calculating the hash
		case Send:
			minTx := Transaction{ // Fill in important parts for the hash
				Hash: "",
				PreviousHash: tx.PreviousHash,
				Action: Send,
				Balance: tx.Balance,
				Currency: tx.Currency,
				Origin: tx.Origin,
				Destination: tx.Destination,
				Signature: "",
			}

			txJson, err := json.Marshal(minTx) // Encode tx in JSON
			if err != nil {
				return "", err	
			}

			hash := sha256.Sum256(txJson) // Generate SHA-256 hash of json content

		return fmt.Sprintf("%v%x", int(tx.Action), hash[:]), nil // Return hash with the txtype in front for convenience later on
		
		case Claim:
			minTx := Transaction{
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

			return fmt.Sprintf("%v%x", int(tx.Action), hash[:]) + tx.PreviousHash[:1], nil // For claim transaction, also put type of connected tx after hash for convenience
		
		case Create:
			minTx := Transaction{
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

			return fmt.Sprintf("%v%x", int(tx.Action), hash[:]), nil
		
		case Trust:
			minTx := Transaction{
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

			return fmt.Sprintf("%v%x", int(tx.Action), hash[:]), nil
		
		default:
			return "", errors.New("Invalid Transaction Type")
	}
}

// Verify if tx is valid
func (tx *Transaction) Verify() bool {
	switch tx.Action { // Each txtype has other factors to determine if tx is valid
		case Send:
			// TODO: Check previous hash
			// TODO: Check if address has sufficient balance
			if tx.Balance < 0 { // Balance can't be negative
				return false
			}	
			if address.ValidateAddress(tx.Origin) != true {
				return false
			}
			if address.ValidateAddress(tx.Destination) != true {
				return false
			}
		case Claim:
			// TODO: Check origin hash
			if address.ValidateAddress(tx.Destination) != true {
				return false
			}
		case Create:
			// TODO: Check previous hash
			if tx.Balance != 0 && tx.Currency == NativeCurrency() { // For account creation balance can't be != 0; token creation balance must be != 0
				return false
			}
			if tx.Balance == 0 && tx.Currency != NativeCurrency() {
				return false
			}
			if address.ValidateAddress(tx.Origin) != true {
				return false
			}
		case Trust:
			// TODO: Check previous hash
			currentTime := time.Now().UnixNano()
			expiring, _ := strconv.ParseInt(tx.Expiration, 10, 64)
			// TODO: Check for error
			if currentTime > expiring && expiring != 0 { // Check if trust certificate isn't expired; token trust expiring is always 0
				return false
			}
			if address.ValidateAddress(tx.Origin) != true {
				return false
			}
			if address.ValidateAddress(tx.Destination) != true {
				return false
			}
	}
	
	if h, _ := tx.GenerateHash(); tx.Hash != h { // Check the authenticity of the content
		return false
	}
	 
	if tx.Action != Create { // Check previous hash
			
	}
		
	pubkey := address.AddressToPubKey(tx.Origin) // Extract public key from address
	curve := elliptic.P256() // Init curve
	
	signatureBytes, _ := base64.StdEncoding.DecodeString(tx.Signature) // Extract signature
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
	if ecdsa.Verify(&rawPubKey, []byte(tx.Hash), &r, &s) == false {
		return false
	}
	
	return true
}

func NewSendTransaction(account address.KeyPair, ph string, destination string, amount uint64, c Currency) (Transaction, error) {
	tx := Transaction{
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

func NewClaimTransaction(account address.KeyPair, ph string, txId string) (Transaction, error) {
	tx := Transaction{
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

func NewCreateTransaction(account string) (Transaction, error) {
	tx := Transaction{
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

func NewCreateTokenTransaction(account address.KeyPair, c Currency, amount uint64) (Transaction, error) {
	tx := Transaction{
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

func NewTrustTransaction(account address.KeyPair, destination string, expiration string) (Transaction, error) {
	tx := Transaction{
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



























