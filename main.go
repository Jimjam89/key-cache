package main

import (
	"log"
	"net/http"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	keyServer := keyServer{
		db: db,
	}

	http.HandleFunc("/", keyServer.getKeyHandler)
	http.HandleFunc("/basic", keyServer.basicGetKeyHandler)

	http.ListenAndServe(":8000", nil)
}

type keyServer struct {
	db *badger.DB
}

func (k keyServer) fetchValueByKey(key string) *badger.Item {
	var value *badger.Item
	var err error

	err = k.db.View(func(txn *badger.Txn) error {
		value, err = txn.Get([]byte(key))

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil
	}

	return value
}

func (k keyServer) setKey(key, value string) error {

	err := k.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})

	if err != nil {
		return err
	}

	return nil
}

func (k keyServer) getKeyHandler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("key") {
		w.Write([]byte("key missing"))
	}

	key := r.URL.Query().Get("key")

	value := k.fetchValueByKey(key)

	if value != nil && value.String() != "" {
		w.Write([]byte(value.String()))

		return
	}

	vaultValue, err := loadFromVault(key)

	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	err = k.setKey(key, vaultValue)

	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	w.Write([]byte(vaultValue))
}

func (k keyServer) basicGetKeyHandler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("key") {
		w.Write([]byte("key missing"))
	}

	key := r.URL.Query().Get("key")

	values := make(map[string]string)

	values["test"] = "aaa"

	w.Write([]byte(values[key]))
}

func loadFromVault(key string) (string, error) {
	time.Sleep(time.Millisecond * 500)

	return "aaa", nil
}
