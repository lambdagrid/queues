package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s Server) createQueue() httprouter.Handle {
	type createQueueRequest struct {
		QueueName *string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		apikey := r.Header.Get("X-API-Key")
		var ownerid int
		err := s.DB.Get(&ownerid, `SELECT id FROM accounts WHERE auth_key = $1`, apikey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// insert new row
	}
}
