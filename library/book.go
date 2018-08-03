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

// Define some constants used
const (
	maxTxs = 50 // Max txs in one book
)

// Define default structure for a book
type book struct {
	Name string
	Txs  []transaction.Transaction
}

// Return a new, empty txbook
func newBook(name string) book {
	return book{name, make([]transaction.Transaction, 0)}
}

// Save a book to the library of specified user
func (b *book) shelve(owner string) {
	bookJson, _ := json.Marshal(b)

	var gz bytes.Buffer // Create a buffer for the gzip to write to
	zipper := gzip.NewWriter(&gz)
	zipper.Write(bookJson) // write contents
	zipper.Close() // close the gzipWriter

	name := fmt.Sprintf("%x.json.gz", b.Name) // Concatenate file name

	p := filepath.Join(getPath(), owner, name) // Get filepath
	ioutil.WriteFile(p, []byte(gz.String()), 0644) // Save the file
}

// Opens a book to write in it
func retrieveBook(path string) book {
	f, _ := os.Open(path)
	defer f.Close() // Make sure it gets closed

	gr, _ := gzip.NewReader(f) // Because everything is saved in gzip we need gzipReader
	defer gr.Close() // Make sure it gets closed

	var contents bytes.Buffer
	contents.ReadFrom(gr)

	var b book
	_ = json.Unmarshal(contents.Bytes(), &b) // Extract data in a book struct

	return b
}

// Add a tx to the tx list in book
func (b *book) addTx(tx transaction.Transaction) {
	b.Txs = append(b.Txs, tx)
}

// Find a tx in the book
func (b *book) getTx(hash string) transaction.Transaction {
	for _, tx := range b.Txs {	// Loop over every tx
		if tx.Hash == hash { // Check hash
			return tx
		}
	}
	return transaction.Transaction{} // Return empty tx on failure
}
