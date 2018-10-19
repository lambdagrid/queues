package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s Server) failMessage() httprouter.Handle {
	type failMessageRequest struct {
		QueueName *string `json:"name"`
		QueueType *string `json:"type"`
	}

	type failMessageResponse struct {
		ErrorMessage *string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	}
}
