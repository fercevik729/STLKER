package handlers

import (
	"context"
	"net/http"

	"github.com/fercevik729/STLKER/octopus/data"
)

// Cache is a middleware that checks if a user cached their portfolio's profits in the past
// 15 minutes
func (c *ControlHandler) Cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement this for various GET routes
		// Get the key
		key := c.retrieveUsername(r) + r.RequestURI
		var profits Profits
		// Get the portfolio from the cache
		if err := c.cache.Get(context.Background(), key, &profits); err == nil {
			c.l.Println("[INFO] Using cache for key", key)
			data.ToJSON(profits, w)
		} else {
			next.ServeHTTP(w, r)
		}

	})
}
