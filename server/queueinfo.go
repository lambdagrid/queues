package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/julienschmidt/httprouter"
)

func (s Server) queueInfo() httprouter.Handle {

	type publicQueueInfo struct {
		Name                   string `json:"name"`
		ApproxNumMessages      int    `json:"approx_num_messages"`
		ApproxMessagesInflight int    `json:"approx_messages_inflight"`
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

		// now search the queues
		var queue Queue
		// look up the queue
		stmt := `SELECT * FROM queues WHERE name = $1`
		err = s.DB.Get(&queue, stmt, queueName)
		if err != nil || queue.OwnerID != ownerid {
			http.Error(w, "Error looking up queue", http.StatusInternalServerError)
			return
		}

		var resp publicQueueInfo
		resp.Name = queue.Name

		result, err := s.sqs.GetQueueAttributes(&sqs.GetQueueAttributesInput{
			AttributeNames: []*string{
				aws.String("ApproximateNumberOfMessages"),
				aws.String("ApproximateNumberOfMessagesNotVisible"),
			},
			QueueUrl: queue.QueueURL,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp.ApproxNumMessages, err = strconv.Atoi(*result.Attributes["ApproximateNumberOfMessages"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.ApproxMessagesInflight, err = strconv.Atoi(*result.Attributes["ApproximateNumberOfMessagesNotVisible"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respData, err := json.Marshal(&resp)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(respData)
	}
}
