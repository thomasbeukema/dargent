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

type Account struct {
    Type        address.AccountType
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

func getPathByAddress(key string) string {
    wd, err := os.Getwd()
    if err != nil {
        // TODO: Proper err handling
        panic(err)
    }

    return filepath.Join(wd, "data", key)
}

// TODO: Determine t by PublicKey automatically
func OpenAccount(addr string, publicKey []byte) Account {
    path := getPathByAddress(addr)
    var t address.AccountType = address.TypeOfAddress(addr)

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
        acc := Account{
            Type: t,
            PublicKey: publicKey,
            Address: addr,
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
    path := getPathByAddress(b64PubKey)

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
            TxList: make([]Transaction, 0),
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

func GetPublicKeyFromAddress(addr string) string {
    acc := OpenAccount(addr, nil)
    led := acc.OpenLedger(NativeCurrency().Ticker)

    return led.TxList[0].Origin
}
