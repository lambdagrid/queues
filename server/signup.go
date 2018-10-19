package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s Server) signup() httprouter.Handle {
	type signupRequest struct {
		Name string `json:"account_name"`
	}

	type signupResponse struct {
		ErrorMessage *string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	}
}
