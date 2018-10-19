package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s Server) createQueue() httprouter.Handle {
	type createQueueRequest struct {
		QueueName *string `json:"name"`
		QueueType *string `json:"type"`
	}

	type createQueueResponse struct {
		ErrorMessage *string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	}
}
