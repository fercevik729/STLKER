package middleware

import (
	"context"
	"net/http"

	"github.com/fercevik729/STLKER/octopus/handlers"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate token
		status, claims := handlers.ValidateJWT(r)
		if status != http.StatusOK {
			w.WriteHeader(status)
		} else {
			// Use contexts to pass email address to subsequent handlers
			username := claims.Username
			ctx := context.WithValue(r.Context(), handlers.Username{}, string(username))
			r = r.WithContext(ctx)
			// Call next handler if the token was valid
			next.ServeHTTP(w, r)
		}

	})

}
