package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/julienschmidt/httprouter"
)

func (s Server) completeMessage() httprouter.Handle {
	type completeMessageRequest struct {
		QueueName     *string `json:"name"`
		ReceiptHandle *string `json:"receipt_handle"`
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
		var req completeMessageRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Could not decode body", http.StatusBadRequest)
			return
		}

		if req.QueueName == nil || req.ReceiptHandle == nil {
			http.Error(w, "Must include the required fields", http.StatusBadRequest)
			return
		}

		var queue Queue
		// look up the queue
		stmt := `SELECT * FROM queues WHERE name = $1`
		err = s.DB.Get(&queue, stmt, *req.QueueName)
		if err != nil || queue.OwnerID != ownerid {
			http.Error(w, "Error looking up queue", http.StatusInternalServerError)
			return
		}

		// authorized, now complete the message
		_, err = s.sqs.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      queue.QueueURL,
			ReceiptHandle: req.ReceiptHandle,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
