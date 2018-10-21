package main

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	authstore "github.com/lambdagrid/queues/auth/jankpersisted"
	"github.com/lambdagrid/queues/server"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=lambdagrid dbname=lambdagrid sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	ap := authstore.New(db)
	s := server.New(ap, db)

	log.Fatal(http.ListenAndServe(":8080", s.GetRouter()))
}
