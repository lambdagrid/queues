package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s Server) completeMessage() httprouter.Handle {
	type completeMessageRequest struct {
		QueueName *string `json:"name"`
		QueueType *string `json:"type"`
	}

	type completeMessageResponse struct {
		ErrorMessage *string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	}
}
