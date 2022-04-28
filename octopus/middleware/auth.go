package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fercevik729/STLKER/octopus/data"
	"github.com/fercevik729/STLKER/octopus/handlers"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate token
		status, claims := handlers.ValidateJWT(r, "token")
		if status != http.StatusOK {
			w.WriteHeader(status)
			data.ToJSON(handlers.ResponseMessage{
				Msg: fmt.Sprintf("Error %d", status),
			}, w)
		} else {
			// Use contexts to pass username and isAdmin to subsequent handlers
			username := claims.Name
			isAdmin := claims.Admin
			ctx := context.WithValue(r.Context(), handlers.Username{}, string(username))
			ctx = context.WithValue(ctx, handlers.IsAdmin{}, bool(isAdmin))
			r = r.WithContext(ctx)
			// Call next handler if the token was valid
			next.ServeHTTP(w, r)
		}

	})

}
