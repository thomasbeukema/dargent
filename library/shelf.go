package library

import (
	"io/ioutil"
	"bytes"
    "compress/gzip"
    "encoding/json"
	"os"
	"path/filepath"
	
	"github.com/thomasbeukema/dargent/transaction"
	"github.com/thomasbeukema/dargent/address"
)

type Shelf struct {
	CreateTx	transaction.CreateTransaction
	Txs			[]string
}

func NewShelf(account address.KeyPair) Shelf {
	tx, _ := transaction.NewCreateTransaction(account.GetAddress())
	tx.Signature = account.SignTx([]byte(tx.Hash))
	
	s := Shelf{tx, []string{}}
	
	return s
}

func AddShelf(account string) {
}

func (s *Shelf) AddToLibrary() {
	owner := s.CreateTx.Origin
	shelfJson,_ := json.Marshal(s)
	
	var gz bytes.Buffer
	zipper := gzip.NewWriter(&gz)
	zipper.Write(shelfJson)
	zipper.Close()
	
	wd, _ := os.Getwd()
	
	p := filepath.Join(wd, "data", owner)
	
	os.MkdirAll(p, os.ModePerm)
	ioutil.WriteFile(filepath.Join(p, "shelf.json.gz"), []byte(gz.String()), 0644)
	ioutil.WriteFile(filepath.Join(p, "shelf.json"), shelfJson, 0644)
}

func (s *Shelf) ShelveTransaction() {
	
}