package middleware

import (
	"fmt"
	"net/http"

	"github.com/lambdagrid/queues/auth"

	"github.com/julienschmidt/httprouter"
)

func HeaderAuth(ap auth.AuthProvider, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		APIKey := r.Header.Get("X-API-Key")
		APISecret := r.Header.Get("X-API-Secret")
		if APIKey == "" {
			http.Error(w, "Request requires API key", http.StatusBadRequest)
			return
		}

		if APISecret == "" {
			http.Error(w, "Request requires API secret", http.StatusBadRequest)
			return
		}
		valid, err := ap.Check(APIKey, APISecret)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error checking authentication:%s ", err.Error()), http.StatusInternalServerError)
			return
		}

		if !valid {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		next(w, r, ps)
	}
}
