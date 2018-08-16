package account

import (
	"fmt"
	"encoding/json"
	"encoding/base64"
	"crypto/sha256"
	_ "crypto/ecdsa"
	_ "crypto/elliptic"
	_ "math/big"
	"time"
	"strconv"
	"errors"

	"github.com/thomasbeukema/dargent/address"
)

// Define possible transaction types
type transactionType int
const (
	SEND transactionType = iota // 0
	CLAIM // 1
	CREATE // 2
	TRUST // 3
)

// Define struct for transaction structure
type Transaction struct {
	Hash			string				`json:"h"`				// Hash for authenticity and txId
	PreviousHash	string				`json:"p,omitempty"`	// Hash of the previous tx
	Action			transactionType		`json:"a"`				// Type of transaction [SEND, CLAIM, CREATE, TRUST]
	Balance			uint64				`json:"b,omitempty"`	// Balance of the address; balance, not the tx amount
	Currency		Currency			`json:"c,omitempty"`	// Currency of the tx
	Origin			string				`json:"o"`				// SENDer / src tx for claim tx
	Destination		string				`json:"d,omitempty"`	// Receiver
	Expiration		string				`json:"e,omitempty"`	// Expiration for trust certificates
}

// Generate hash for a transaction to ensure the authenticity of the contents of the transaction
func (tx *Transaction) GenerateHash() (string, error) {
	switch tx.Action { // Every type of transaction has a slightly different way of calculating the hash
		case SEND:
			minTx := Transaction{ // Fill in important parts for the hash
				Hash: "",
				PreviousHash: tx.PreviousHash,
				Action: SEND,
				Balance: tx.Balance,
				Currency: tx.Currency,
				Origin: tx.Origin,
				Destination: tx.Destination,
			}

			txJson, err := json.Marshal(minTx) // Encode tx in JSON
			if err != nil {
				return "", err
			}

			hash := sha256.Sum256(txJson) // Generate SHA-256 hash of json content

		return fmt.Sprintf("%v%x", int(tx.Action), hash[:]), nil // Return hash with the txtype in front for convenience later on

		case CLAIM:
			minTx := Transaction{
				Hash: "",
				PreviousHash: tx.PreviousHash,
				Action: SEND,
				Origin: tx.Origin,
				Destination: tx.Destination,
			}

			txJson, err := json.Marshal(minTx)
			if err != nil {
				return "", err
			}

			hash := sha256.Sum256(txJson)

			return fmt.Sprintf("%v%x", int(tx.Action), hash[:]) + tx.PreviousHash[:1], nil // For claim transaction, also put type of connected tx after hash for convenience

		case CREATE:
			minTx := Transaction{
				Hash: "",
				Action: CREATE,
				Balance: tx.Balance,
				Currency: tx.Currency,
				Origin: tx.Origin,
			}

			txJson, err := json.Marshal(minTx)
			if err != nil {
				return "", err
			}

			hash := sha256.Sum256(txJson)

			return fmt.Sprintf("%v%x", int(tx.Action), hash[:]), nil

		case TRUST:
			minTx := Transaction{
				Hash: "",
				PreviousHash: tx.PreviousHash,
				Action: TRUST,
				Origin: tx.Origin,
				Destination: tx.Destination,
				Expiration: tx.Origin,
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
		case SEND:
			if tx.Balance < 0 { // Balance can't be negative
				return false
			}
			if address.ValidateAddress(tx.Origin) != true {
				return false
			}
			if address.ValidateAddress(tx.Destination) != true {
				return false
			}
		case CLAIM:
			// TODO: Check origin tx
			if address.ValidateAddress(tx.Destination) != true {
				return false
			}
		case CREATE:
			if tx.Balance == 0 && tx.Currency != NativeCurrency() {
				return false
			}
			if tx.Currency != NativeCurrency() { // No origin when creating token, since it's stored in currency struct
				if address.ValidateAddress(tx.Currency.Owner) != true {
					return false
				}
			} else {
				if address.ValidateAddress(tx.Origin) != true {
					return false
				}
			}
		case TRUST:
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

	return true
}

func NewSendTransaction(account string, ph string, destination string, amount uint64, c Currency) (Transaction, error) {
	tx := Transaction{
		Hash: "",
		PreviousHash: ph,
		Action: SEND,
		Balance: amount,
		Currency: c,
		Origin: account,
		Destination: destination,
	}

	tx.Hash,_ = tx.GenerateHash()

	return tx, nil
}

func NewClaimTransaction(account string, ph string, txId string) (Transaction, error) {
	tx := Transaction{
		Hash: "",
		PreviousHash: ph,
		Action: CLAIM,
		Origin: txId,
		Destination: account,
	}

	tx.Hash,_ = tx.GenerateHash()

	return tx, nil
}

func NewCreateTransaction(pubkey []byte) (Transaction, error) {
	b64pubkey := base64.StdEncoding.EncodeToString(pubkey)

	tx := Transaction{
		Hash: "",
		Action: CREATE,
		Currency: NativeCurrency(),
		Balance: 0,
		Origin: b64pubkey,
	}

	tx.Hash,_ = tx.GenerateHash()

	return tx, nil
}

func NewCreateTokenTransaction(pubkey string, c Currency, amount uint64) (Transaction, error) {
	tx := Transaction{
		Hash: "",
		Action: CREATE,
		Currency: c,
		Balance: amount,
	}

	tx.Hash,_ = tx.GenerateHash()

	return tx, nil
}

func NewTrustTransaction(account string, destination string, expiration string) (Transaction, error) {
	tx := Transaction{
		Hash: "",
		PreviousHash: "",
		Action: TRUST,
		Origin: account,
		Destination: destination,
		Expiration: expiration,
	}

	tx.Hash,_ = tx.GenerateHash()
	tx.PreviousHash,_ = tx.GenerateHash()

	return tx, nil
}
