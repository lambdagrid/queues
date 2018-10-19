package main

import (
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	authstore "github.com/lambdagrid/queues/auth/jankpersisted"
	"github.com/lambdagrid/queues/server"
)

func main() {
	db, err := setupDatabase("mytest.db")
	if err != nil {
		panic(err)
	}
	ap := authstore.New(db)
	s := server.New(ap)

	log.Fatal(http.ListenAndServe(":8080", s.GetRouter()))
}

// this is really fucking jank but YOLO
func setupDatabase(dbpath string) (*bolt.DB, error) {
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		return db, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("accounts"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("queues"))
		if err != nil {
			return err
		}

		return nil
	})

	return db, err
}
