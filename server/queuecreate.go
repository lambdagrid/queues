package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/julienschmidt/httprouter"
)

func (s Server) createQueue() httprouter.Handle {
	type createQueueRequest struct {
		QueueName *string `json:"name"`
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	if err != nil {
		panic(err)
	}

	svc := sqs.New(sess)

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

		var req createQueueRequest
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

		stmt := `INSERT INTO queues (name, owner_id) VALUES ($1, $2)`
		rows, err := s.DB.Query(stmt, *req.QueueName, ownerid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		rows.Close()

		// now create the queue on the backend, TODO: make safer

		result, err := svc.CreateQueue(&sqs.CreateQueueInput{
			QueueName: aws.String(fmt.Sprintf("%d_%s.fifo", ownerid, *req.QueueName)),
			Attributes: map[string]*string{
				"FifoQueue": aws.String("true"),
			},
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stmt = `UPDATE queues SET queue_url = $1 WHERE name = $2`
		_, err = s.DB.Queryx(stmt, result.QueueUrl, *req.QueueName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
