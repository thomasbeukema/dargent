package library

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/thomasbeukema/dargent/address"
	"github.com/thomasbeukema/dargent/transaction"
)

// Define struct for storage
type Shelf struct {
	CreateTx transaction.Transaction
	Txs      map[string][]string
}

func getPath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "data")
}

func NewShelf(account address.KeyPair) Shelf {
	tx, _ := transaction.NewCreateTransaction(account.GetAddress())
	tx.Signature = account.SignTx([]byte(tx.Hash))

	s := Shelf{tx, make(map[string][]string)}

	return s
}

func OpenShelf(account string) Shelf {
	path := filepath.Join(getPath(), account, "shelf.json.gz")

	f, _ := os.Open(path)
	defer f.Close()

	gr, _ := gzip.NewReader(f)
	defer gr.Close()

	var contents bytes.Buffer
	contents.ReadFrom(gr)

	var s Shelf
	_ = json.Unmarshal(contents.Bytes(), &s)

	return s
}

func AddShelf(account string, tx transaction.Transaction) {

}

func (s *Shelf) UpdateLibrary() {
	owner := s.CreateTx.Origin
	shelfJson, _ := json.Marshal(s)

	var gz bytes.Buffer
	zipper := gzip.NewWriter(&gz)
	zipper.Write(shelfJson)
	zipper.Close()

	p := filepath.Join(getPath(), owner)

	os.MkdirAll(p, os.ModePerm)
	ioutil.WriteFile(filepath.Join(p, "shelf.json.gz"), []byte(gz.String()), 0644)
}

func (s *Shelf) ShelveTx(tx transaction.Transaction) bool {
	if s.CreateTx.Origin != tx.Origin {
		return false
	}
	if tx.Verify() == true {
		b := s.latestBook()
		b.addTx(tx)
		s.shelveBook(b)
		s.Txs[b.Name] = append(s.Txs[b.Name], tx.Hash)
		s.UpdateLibrary()
		return true
	}
	return false
}

func (s *Shelf) shelveBook(b book) {
	b.shelve(s.CreateTx.Origin)
}

func (s *Shelf) latestBook() book {
	if len(s.Txs) == 0 {
		return newBook(string(len(s.Txs) + 1))
	} else {
		if len(s.Txs[string(len(s.Txs)-1)]) == maxTxs {
			return newBook(string(len(s.Txs)))
		} else {
			name := fmt.Sprintf("%x.json.gz", string(len(s.Txs)))
			path := filepath.Join(getPath(), s.CreateTx.Origin, name)
			return retrieveBook(path)
		}
	}
}

func (s *Shelf) FindTx(hash string) transaction.Transaction {
	for k, v := range s.Txs {
		for _, b := range v {
			if b == hash {
				bk := retrieveBook(filepath.Join(getPath(), s.CreateTx.Origin, k))
				return bk.getTx(hash)
			}
		}
	}
	return transaction.Transaction{}
}

func (s *Shelf) Newest() transaction.Transaction {
	if len(s.Txs) == 0 {
		return s.CreateTx
	} else {
		latest := fmt.Sprintf("%x", string(len(s.Txs)))
		path := filepath.Join(getPath(), s.CreateTx.Origin, string(latest)+".json.gz")
		latestBook := retrieveBook(path)
		latestInBook := len(latestBook.Txs) - 1

		return latestBook.Txs[latestInBook]
	}
	return transaction.Transaction{}
}
