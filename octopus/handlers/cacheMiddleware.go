package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/fercevik729/STLKER/octopus/data"
)

// Cache is a middleware that checks if a user cached their portfolio's profits in the past
// 15 minutes
func (c *ControlHandler) Cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		// Get the key
		key := c.retrieveUsername(r) + r.RequestURI
		slashCount := strings.Count(r.RequestURI, "/")
		switch slashCount {
		// TODO: fix this
		// All portfolios for a user
		case 1:
			if c.retrieveAdmin(r) {
				var content map[string][]string
				err = c.getFromCache(key, &content, w)
			} else {
				var content []*Profits
				err = c.getFromCache(key, &content, w)
			}
		// Single portfolio
		case 2:
			var content Profits
			err = c.getFromCache(key, &content, w)
		// Single security
		case 3:
			var content *Security
			err = c.getFromCache(key, &content, w)
		}

		// If there was an error retrieving, serve the next handler
		if err != nil {
			next.ServeHTTP(w, r)
		}

	})
}

func (c *ControlHandler) getFromCache(key string, content interface{}, w http.ResponseWriter) error {
	err := c.cache.Get(context.Background(), key, &content)
	if err == nil {
		c.l.Println("[INFO] Using cache for key", key)
		data.ToJSON(content, w)
	}
	return err
}
