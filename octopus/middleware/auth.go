package middleware

import (
	"net/http"

	"github.com/fercevik729/STLKER/octopus/handlers"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Validate token
		status, _ := handlers.ValidateJWT(r)
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		// Call next handler if the token was valid
		next.ServeHTTP(w, r)

	}

}
