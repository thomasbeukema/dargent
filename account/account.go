package account

import (
    "os"
    "encoding/base64"
    "path/filepath"
    "encoding/json"
    "compress/gzip"
    "io/ioutil"
    "bytes"

    "github.com/thomasbeukema/dargent/address"
)

type AccountType int
const (
    ECC AccountType = iota // 0
    SPHINCS // 1
)

type Account struct {
    Type        AccountType
    PublicKey   []byte
    Address     string
    Currencies  []string
}

// Simple function whichc checks whether file exists
func pathExists(path string) bool {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return false
    } else if _, err := os.Stat(path); err == nil {
        return true
    } else {
        // TODO: Handle unexpected err
        return false
    }
}

func getPathFromPubKey(key string) string {
    wd, err := os.Getwd()
    if err != nil {
        // TODO: Proper err handling
        panic(err)
    }

    return filepath.Join(wd, "data", key)
}

// TODO: Determine t by PublicKey automatically
func OpenAccount(publicKey []byte, t AccountType) Account {
    b64PubKey := base64.StdEncoding.EncodeToString(publicKey)
    path := getPathFromPubKey(b64PubKey)

    if pathExists(path) { // Already opened this account, return saved data
        f, err := os.Open(filepath.Join(path, "index.json.gz"))
        defer f.Close()
        if err != nil {
            // TODO: Proper err handling
            panic(err)
        }

        gzipReader, err := gzip.NewReader(f)
        defer gzipReader.Close()
        if err != nil {
            // TODO: You guessed it... Proper err handling
            panic(err)
        }

        var content bytes.Buffer
        content.ReadFrom(gzipReader)

        var acc Account
        err = json.Unmarshal(content.Bytes(), &acc)
        if err != nil {
            // TODO: Proper err handling
            panic(err)
        }

        return acc

    } else { // Return new account and save it to disk
        var ad string
        switch t {
        case ECC: // ECC account
            ad = address.ECCPubKeyToAddress(publicKey)
        case SPHINCS: // SPHINCS account
            ad = address.SPHINCSPubKeyToAddress(publicKey)
        default:
            // TODO
            panic(1)
        }
        acc := Account{
            Type: t,
            PublicKey: publicKey,
            Address: ad,
            Currencies: make([]string, 0),
        }

        accountJson, err := json.Marshal(acc)
        if err != nil {
            // TODO
            panic(err)
        }

        var gzipped bytes.Buffer
        gzipper := gzip.NewWriter(&gzipped)
        gzipper.Write(accountJson)
        gzipper.Close()

        os.MkdirAll(path, os.ModePerm)
        ioutil.WriteFile(filepath.Join(path, "index.json.gz"), []byte(gzipped.String()), 0644)

        return acc
    }
}

func (acc *Account) getLedgerPath(currency string) string {
    b64PubKey := base64.StdEncoding.EncodeToString(acc.PublicKey)
    path := getPathFromPubKey(b64PubKey)

    return filepath.Join(path, currency)
}

func (acc *Account) OpenLedger(currency string) Ledger {
    ledgerPath := acc.getLedgerPath(currency)

    if pathExists(ledgerPath) { // Already created ledger
        f, err := os.Open(filepath.Join(ledgerPath, "index.json.gzip"))
        defer f.Close()
        if err != nil {
            // TODO: Proper err handling
            panic(err)
        }

        gzipReader, err := gzip.NewReader(f)
        defer gzipReader.Close()
        if err != nil {
            // TODO: You guessed it... Proper err handling
            panic(err)
        }

        var content bytes.Buffer
        content.ReadFrom(gzipReader)

        var led Ledger
        err = json.Unmarshal(content.Bytes(), &led)
        if err != nil {
            // TODO: Proper err handling
            panic(err)
        }

        return led
    } else { // Create ledger
        led := Ledger{
            Path: ledgerPath,
            Currency: currency,
            TxList: transactionList{
                Headers: make([]string, 0),
                Txs: make([]Transaction, 0),
            },
            Signature: "",
        }

        ledgerJson, err := json.Marshal(led)
        if err != nil {
            // TODO
            panic(err)
        }

        var gzipped bytes.Buffer
        gzipper := gzip.NewWriter(&gzipped)
        gzipper.Write(ledgerJson)
        gzipper.Close()

        os.MkdirAll(ledgerPath, os.ModePerm)

        ioutil.WriteFile(filepath.Join(ledgerPath, "index.json.gz"), []byte(gzipped.String()), 0644)

        return led
    }
}
