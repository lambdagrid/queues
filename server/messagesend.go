package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s Server) sendMessage() httprouter.Handle {
	type sendMessageRequest struct {
		QueueName *string `json:"name"`
		QueueType *string `json:"type"`
	}

	type sendMessageResponse struct {
		ErrorMessage *string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	}
}
