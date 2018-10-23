package jankpersisted

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"

	"github.com/jmoiron/sqlx"
	"github.com/lambdagrid/queues/auth"
)

type JankPersistedAuthStore struct {
	db *sqlx.DB
}

type accountRecord struct {
	ID         int    `db:"id"`
	AuthKey    string `db:"auth_key"` // auth key is account key, super jank
	Name       string `db:"account_name"`
	AuthSecret string `db:"auth_secret"`
}

func (j *JankPersistedAuthStore) Check(key, secret string) (bool, error) {
	valid := false

	var rec accountRecord

	err := j.db.Get(&rec, "SELECT * FROM accounts WHERE auth_key = $1", key)
	if err == sql.ErrNoRows {
		return valid, nil
	} else if err != nil {
		return valid, err
	}
	valid = rec.AuthSecret == secret
	return valid, err
}

func (j *JankPersistedAuthStore) CreateAccount(accountName string) (key, secret string, err error) {
	key, secret, err = generateAPIKeyAndSecret()
	if err != nil {
		return
	}

	stmt := `INSERT INTO accounts (account_name, auth_key, auth_secret) VALUES ($1, $2, $3)`

	rows, err := j.db.Query(stmt, accountName, key, secret)
	if err != nil {
		return
	}

	defer rows.Close()
	if rows.Err() != nil {
		err = rows.Err()
	}

	return
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

func New(db *sqlx.DB) auth.AuthProvider {
	return &JankPersistedAuthStore{db}
}
