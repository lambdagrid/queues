package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mauth "github.com/lambdagrid/queues/auth/mock"
)

func TestWrite(t *testing.T) {
	a := mauth.New()
	s := New(a)

	b := struct {
		QueueID int64  `json:"queue_id"`
		Body    string `json:"body"`
	}{
		QueueID: 10,
		Body:    "blah",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(b)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/v1/jobs", &buf)
	req.Header.Set("X-API-Key", "test")
	req.Header.Set("X-API-Secret", "test")
	w := httptest.NewRecorder()
	s.GetRouter().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}

	var resp struct {
		SequenceNumber *int64 `json:"seq"`
	}

	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	if resp.SequenceNumber == nil {
		t.Fatal()
	}
	if *resp.SequenceNumber != 0 {
		t.Fail()
	}
}
