package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s Server) signup() httprouter.Handle {
	type signupRequest struct {
		Name *string `json:"account_name"`
	}

	type signupResponse struct {
		Key    *string `json:"authKey"`
		Secret *string `json:"authSecret"`
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read body", http.StatusBadRequest)
			return
		}

		var req signupRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Could not decode body", http.StatusBadRequest)
			return
		}

		if req.Name == nil {
			http.Error(w, "Need to provide name", http.StatusBadRequest)
			return
		}

		key, secret, err := s.authProvider.CreateAccount(*req.Name)
		if err != nil {
			http.Error(w, fmt.Sprintf("There was an issue creating the account: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		resp := signupResponse{
			Key:    &key,
			Secret: &secret,
		}

		log.Println("Signed up account ", *req.Name, " with key ", key, "/", secret)

		respData, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(respData)
	}
}
