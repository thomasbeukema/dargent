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
    TxList      transactionList
    Signature   string
}

type transactionList struct {
    Headers []string
    Txs     []Transaction
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
    if tx.Verify() && tx.Origin == led.TxList.Txs[0].Origin {
        led.TxList.Headers = append(led.TxList.Headers, tx.Hash)
        led.TxList.Txs = append(led.TxList.Txs, tx)
    }
}

func (led *Ledger) VerifySignature(acc Account) {

}
