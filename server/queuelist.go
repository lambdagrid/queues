package server

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Queue struct {
	Name     string  `db:"name"`
	OwnerID  string  `db:"owner_id"`
	QueueURL *string `db:"queue_url"`
}

func (s Server) listQueues() httprouter.Handle {
	type listQueueResponse struct {
		Queues []string `json:"queues"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		apikey := r.Header.Get("X-API-Key")

		var ownerid int
		err := s.DB.Get(&ownerid, `SELECT id FROM accounts WHERE auth_key = $1`, apikey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// now search the queues

		queues := []Queue{}

		err = s.DB.Select(&queues, `SELECT * FROM queues WHERE owner_id = $1`, ownerid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		names := make([]string, len(queues))

		for i, q := range queues {
			names[i] = q.Name
		}

		resp := listQueueResponse{
			Queues: names,
		}
		respData, err := json.Marshal(&resp)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(respData)

		return
	}
}
