package server

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/julienschmidt/httprouter"
)

func (s Server) receiveMessage() httprouter.Handle {

	type receiveMessageResponse struct {
		VisilibityTimeout int    `json:"visibility_timeout"`
		Payload           string `json:"payload"`
		ReceiptHandle     string `json:"receipt_handle"`
	}

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		apikey := r.Header.Get("X-API-Key")
		var ownerid int
		err := s.DB.Get(&ownerid, `SELECT id FROM accounts WHERE auth_key = $1`, apikey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		queueName := ps.ByName("queuename")

		var queue Queue
		// look up the queue
		stmt := `SELECT * FROM queues WHERE name = $1`
		err = s.DB.Get(&queue, stmt, queueName)
		if err != nil || queue.OwnerID != ownerid {
			http.Error(w, "Error looking up queue", http.StatusInternalServerError)
			return
		}

		result, err := s.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            queue.QueueURL,
			MaxNumberOfMessages: aws.Int64(1),
			VisibilityTimeout:   aws.Int64(60), // TODO: Allow configuring visibility timeout
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respMessages := make([]receiveMessageResponse, len(result.Messages))

		for i, msg := range result.Messages {
			respMessages[i].Payload = *msg.Body
			respMessages[i].ReceiptHandle = *msg.ReceiptHandle
			respMessages[i].VisilibityTimeout = 60
		}

		respData, err := json.Marshal(&respMessages)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(respData)
	}
}
