package main

import (
	"log"
	"net/http"

	mauth "github.com/lambdagrid/queues/auth/mock"
	"github.com/lambdagrid/queues/server"
)

func main() {
	ap := mauth.New()
	s := server.New(ap)

	log.Fatal(http.ListenAndServe(":8080", s.GetRouter()))
}
