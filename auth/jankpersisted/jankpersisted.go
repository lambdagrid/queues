package jankpersisted

import (
	"fmt"

	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/lambdagrid/queues/auth"
)

type JankPersistedAuthStore struct {
	db *bolt.DB
}

type accountRecord struct {
	AuthKey    string // auth key is account key, super jank
	Name       string
	AuthSecret string
}

func (j *JankPersistedAuthStore) Check(key, secret string) (bool, error) {
	valid := false
	err := j.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("accounts"))
		accountBytes := b.Get([]byte(key))
		if accountBytes == nil {
			return fmt.Errorf("That account doesn't exist")
		}

		var acc accountRecord
		err := json.Unmarshal(accountBytes, &acc)
		if err != nil {
			return err
		}

		valid = secret != acc.AuthSecret
		return err
	})

	return valid, err
}

func (j *JankPersistedAuthStore) CreateAccount(accountName string) (key, secret string, err error) {
	err = j.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("accounts"))
		// get the account name, fail if it doesn't exist
		existing := b.Get([]byte(accountName))
		if existing != nil {
			return fmt.Errorf("An account exists with that name already")
		}

		key, secret, err = generateAPIKeyAndSecret()
		if err != nil {
			return fmt.Errorf("There was an error generating the key")
		}

		rec := accountRecord{}
		rec.Name = accountName
		rec.AuthKey = key
		rec.AuthSecret = secret

		encoded, err := json.Marshal(rec)
		if err != nil {
			return fmt.Errorf("There was an error serializing the account record")
		}

		b.Put([]byte(key), encoded)
		return nil
	})

	return key, secret, err
}

func generateAPIKeyAndSecret() (key, secret string, err error) {
	b := make([]byte, 64) // first 32 will be key, last 32 secret
	_, err = rand.Read(b)
	if err != nil {
		return "", "", err
	}

	key = base64.StdEncoding.EncodeToString(b[0:32])
	secret = base64.StdEncoding.EncodeToString(b[32:64])

	return key, secret, nil
}

func New(db *bolt.DB) auth.AuthProvider {
	return &JankPersistedAuthStore{db}
}
