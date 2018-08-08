package transaction

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/thomasbeukema/dargent/address"
)

// Define struct for storage
type Shelf struct {
	CreateTx Transaction
	Txs      map[string][]string
}

// Utility func
func getPath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "data")
}

// Return new empty shell with a signed CREATE tx
func NewShelf(account address.KeyPair) Shelf {
	tx, _ := NewCreateTransaction(account.GetAddress())
	tx.Signature = account.SignTx([]byte(tx.Hash))

	s := Shelf{tx, make(map[string][]string)}

	return s
}

// Open an already created shelf
func OpenShelf(account string) Shelf {
	path := filepath.Join(getPath(), account, "shelf.json.gz")

	f, _ := os.Open(path)
	defer f.Close()

	gr, _ := gzip.NewReader(f) // Everything is gzipped
	defer gr.Close()

	var contents bytes.Buffer
	contents.ReadFrom(gr)

	var s Shelf
	_ = json.Unmarshal(contents.Bytes(), &s)

	return s
}

// Update a shelf in the library
func (s *Shelf) UpdateLibrary() {
	owner := s.CreateTx.Origin
	shelfJson, _ := json.Marshal(s)

	var gz bytes.Buffer
	zipper := gzip.NewWriter(&gz)
	zipper.Write(shelfJson)
	zipper.Close()

	p := filepath.Join(getPath(), owner)

	os.MkdirAll(p, os.ModePerm)

	if s.CreateTx.Currency.Ticker != NativeCurrency().Ticker { // Token
		shelfName := "shelf_" + s.CreateTx.Currency.Ticker + ".json.gz"
		ioutil.WriteFile(filepath.Join(p, shelfName), []byte(gz.String()), 0644)
	} else { // Native Currency
		ioutil.WriteFile(filepath.Join(p, "shelf.json.gz"), []byte(gz.String()), 0644)
	}
}

// Add a new tx to a shelf
func (s *Shelf) ShelveTx(tx Transaction) bool {
	if s.CreateTx.Origin != tx.Origin { // Check if txOwner is also owner of the shelf, only txs of owner can be written to a shelf
		return false
	}
	if tx.Verify() == true { // Check if tx is valid
		if tx.Action == 2 {
			tokenShelf := Shelf{tx, make(map[string][]string)} // Create new shelf with tx as opening tx
			tokenShelf.UpdateLibrary()
		} else {
			b := s.latestBook() // Retrieve latest book
			b.addTx(tx)
			s.shelveBook(b)
			s.Txs[b.Name] = append(s.Txs[b.Name], tx.Hash) // Add txhash to list of txs in shelf
			s.UpdateLibrary()
			return true
		}
	}
	return false
}

// Utility func
func (s *Shelf) shelveBook(b book) {
	b.shelve(s.CreateTx.Origin)
}

// Get the latest book in shelf
func (s *Shelf) latestBook() book {
	if len(s.Txs) == 0 { // First tx
		return newBook(string(len(s.Txs) + 1))
	} else {
		if len(s.Txs[string(len(s.Txs)-1)]) == maxTxs {
			return newBook(string(len(s.Txs))) // Book has reached max capacity: return empty one
		} else {
			name := fmt.Sprintf("%x.json.gz", string(len(s.Txs)))
			path := filepath.Join(getPath(), s.CreateTx.Origin, name)
			return retrieveBook(path) // Return latest book
		}
	}
}

// Find a tx in shelf
func (s *Shelf) FindTx(hash string) *Transaction {
	fmt.Println(hash)

	if s.CreateTx.Hash == hash {
		return &s.CreateTx
	}

	for k, v := range s.Txs { // Loop over all books in shelf
		for _, b := range v { // Loop over every tx in book
			if b == hash {
				bk := retrieveBook(filepath.Join(getPath(), s.CreateTx.Origin, k))
				return bk.getTx(hash)
			}
		}
	}
	return &Transaction{} // Return empty at failure
}

// Get the latest tx
func (s *Shelf) Newest() Transaction {
	if len(s.Txs) == 0 {
		return s.CreateTx // If no txs return CREATE tx
	} else {
		latest := fmt.Sprintf("%x", string(len(s.Txs))) // Get newest book name
		path := filepath.Join(getPath(), s.CreateTx.Origin, string(latest)+".json.gz")
		latestBook := retrieveBook(path) // Get latest book
		latestInBook := len(latestBook.Txs) - 1 // Newest tx appended, so last is newest

		return latestBook.Txs[latestInBook] // Return latest tx
	}
	return Transaction{} // Return empty at failure
}
