package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s Server) getStatus() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("Hello"))
	}
}
