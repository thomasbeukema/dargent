package account

import (
    "encoding/json"
    "bytes"
    "compress/gzip"
    "io/ioutil"
    "path/filepath"
)

type Ledger struct {
    Path        string
    Currency    string
    TxList      []Transaction
    Signature   string
}

func (led *Ledger) Write() {
    ledgerJson, err := json.Marshal(*led)
    if err != nil {
        // TODO
        panic(err)
    }

    var gzipped bytes.Buffer
    gzipper := gzip.NewWriter(&gzipped)
    gzipper.Write(ledgerJson)
    gzipper.Close()

    ioutil.WriteFile(filepath.Join(led.Path, "index.json.gz"), []byte(gzipped.String()), 0644)
}

func (led *Ledger) AddTransaction(tx Transaction) {
    if tx.Verify() && tx.Origin == led.TxList[0].Origin {
        led.TxList = append(led.TxList, tx)
    }
}
