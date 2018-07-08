package library

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/thomasbeukema/dargent/transaction"
)

const (
	maxTxs = 17
)

type book struct {
	Name string
	Txs  []transaction.Transaction
}

func newBook(name string) book {
	return book{name, make([]transaction.Transaction, 0)}
}

func (b *book) shelve(owner string) {
	bookJson, _ := json.Marshal(b)

	var gz bytes.Buffer
	zipper := gzip.NewWriter(&gz)
	zipper.Write(bookJson)
	zipper.Close()

	name := fmt.Sprintf("%x.json.gz", b.Name)

	p := filepath.Join(getPath(), owner, name)
	ioutil.WriteFile(p, []byte(gz.String()), 0644)
}

func retrieveBook(path string) book {
	f, _ := os.Open(path)
	defer f.Close()

	gr, _ := gzip.NewReader(f)
	defer gr.Close()

	var contents bytes.Buffer
	contents.ReadFrom(gr)

	var b book
	_ = json.Unmarshal(contents.Bytes(), &b)

	return b
}

func (b *book) addTx(tx transaction.Transaction) {
	b.Txs = append(b.Txs, tx)
}

func (b *book) getTx(hash string) transaction.Transaction {
	for _, tx := range b.Txs {
		if tx.Hash == hash {
			return tx
		}
	}
	return transaction.Transaction{}
}
