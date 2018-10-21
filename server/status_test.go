package server

/*
import (
	"net/http"
	"net/http/httptest"
	"testing"

	mauth "github.com/lambdagrid/queues/auth/mock"
)

func TestStatus(t *testing.T) {
	a := mauth.New()
	s := New(a)
	req, err := http.NewRequest(http.MethodGet, "/v1/status", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	s.GetRouter().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
	if w.Body.String() != "Hello" {
		t.Fail()
	}
}*/
