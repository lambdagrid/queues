package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/julienschmidt/httprouter"
)

func (s Server) sendMessage() httprouter.Handle {
	type sendMessageRequest struct {
		QueueName *string `json:"name"`
		Payload   *string `json:"payload"`
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

		var req sendMessageRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Could not decode body", http.StatusBadRequest)
			return
		}

		if req.Payload == nil || req.QueueName == nil {
			http.Error(w, "Must include the required fields", http.StatusBadRequest)
			return
		}

		// TODO: more message validation
		if len(*req.Payload) < 1 || len(*req.Payload) > 262144 {
			http.Error(w, "Invalid payload length (0 < payload <= 262144)", http.StatusBadRequest)
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

		// we've verified the user can see the queue, let's insert the message now
		// content based deduplication IS ENABLED
		result, err := s.sqs.SendMessage(&sqs.SendMessageInput{
			MessageBody:    aws.String(*req.Payload),
			QueueUrl:       aws.String(*queue.QueueURL),
			MessageGroupId: aws.String("defaultgroupid"),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Sent message, seq number", *result.SequenceNumber)
	}
}
