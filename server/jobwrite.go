package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lambdagrid/queues/record"

	"github.com/julienschmidt/httprouter"
)

func (s Server) handleJobWrite() httprouter.Handle {
	type writeRequest struct {
		QueueID *int64  `json:"queue_id"`
		JobBody *string `json:"body"`
	}

	type writeResponse struct {
		SequenceNumber int64 `json:"seq"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read body", http.StatusBadRequest)
			return
		}

		var req writeRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Could not decode body", http.StatusBadRequest)
			return
		}

		if req.JobBody == nil || req.QueueID == nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
		}

		// perform write operation here
		rec := record.Record{
			SequenceNumber: 0, // todo FIX
			Body:           []byte(*req.JobBody),
		}

		log.Println(rec)

		// TODO: Logger write

		resp := writeResponse{}

		respData, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(respData)
	}
}
