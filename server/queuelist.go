package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s Server) listQueues() httprouter.Handle {
	type listQueueResponse struct {
		Queues []string `json:"queues"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	}
}
