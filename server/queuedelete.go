package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/julienschmidt/httprouter"
)

func (s Server) deleteQueue() httprouter.Handle {
	type deleteQueueRequest struct {
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

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read body", http.StatusBadRequest)
			return
		}

		var req deleteQueueRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Could not decode body", http.StatusBadRequest)
			return
		}

		if req.QueueName == nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
		}

		// TODO: Add validation for alphanum, hyphens, underscores, fifo suffix,
		// and length

		var queue Queue
		stmt := `SELECT * FROM queues WHERE name = $1`
		err = s.DB.Get(&queue, stmt, *req.QueueName)

		_, err = s.sqs.DeleteQueue(&sqs.DeleteQueueInput{
			QueueUrl: queue.QueueURL,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stmt = `DELETE FROM queues WHERE name = $1`
		_, err = s.DB.Queryx(stmt, *req.QueueName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
