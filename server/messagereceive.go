package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s Server) receiveMessage() httprouter.Handle {
	type receiveMessageResponse struct {
		ErrorMessage *string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	}
}
