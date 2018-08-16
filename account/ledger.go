package account

import (
    "encoding/json"
    "bytes"
    "compress/gzip"
    "io/ioutil"
    "path/filepath"
    "encoding/base64"
    "crypto/sha256"

    "github.com/thomasbeukema/dargent/address"
)

type Ledger struct {
    Currency    string
    Hash        string
    TxList      []Transaction
    Signature   string
}

func (led *Ledger) Write(acc *Account) {
    ledgerJson, err := json.Marshal(*led)
    if err != nil {
        // TODO
        panic(err)
    }

    var gzipped bytes.Buffer
    gzipper := gzip.NewWriter(&gzipped)
    gzipper.Write(ledgerJson)
    gzipper.Close()

    ledgerPath := acc.getLedgerPath(led.Currency)

    ioutil.WriteFile(filepath.Join(ledgerPath, "index.json.gz"), []byte(gzipped.String()), 0644)
}

func (led *Ledger) addTransaction(tx Transaction, acc *Account) bool {
    if len(led.TxList) == 0 { // 'CREATE' tx

        if tx.Action != CREATE {
            return false
        }

        decodedOrigin, err := base64.StdEncoding.DecodeString(tx.Origin)
        if err != nil {
            // TODO: Proper err handling
            panic(err)
        }

        if bytes.Equal(acc.PublicKey, decodedOrigin) {
            led.TxList = append(led.TxList, tx)
            led.CalculateHash()
            led.Write(acc)
        }
    } else { // SEND, CLAIM, TRUST
        if tx.Verify() {
            led.TxList = append(led.TxList, tx)
            led.CalculateHash()
            led.Write(acc)
        }
    }

    return true
}

func (led *Ledger) CalculateHash() string {
    hashes := ""

    for _, tx := range led.TxList {
        hashes = hashes + ":" + tx.Hash
    }

    firstHash := sha256.Sum256([]byte(hashes))
    finalHash := sha256.Sum256(firstHash[:])

    led.Hash = base64.StdEncoding.EncodeToString(finalHash[:])

    return led.Hash
}

func (led *Ledger) UpdateSignature(signature string, acc *Account) bool {
    b64DecPubKey, err := base64.StdEncoding.DecodeString(led.TxList[0].Origin)
    if err != nil {
        return false
    }
    switch len(signature) {
    case 88: // ECDSA
        if address.ValidateECCSignature(signature, []byte(led.Hash), b64DecPubKey) {
            led.Signature = signature
            led.Write(acc)
        } else {
            return false
        }
    case 54668: // SPHINCS
    var sphincsPubKey [1056]byte

    for i := range b64DecPubKey {
        sphincsPubKey[i] = b64DecPubKey[i]
    }

    if address.ValidateSPHINCSSignature(signature, []byte(led.Hash), &sphincsPubKey) {
        led.Signature = signature
        led.Write(acc)
    } else {
        return false
    }
    }
    return true
}
