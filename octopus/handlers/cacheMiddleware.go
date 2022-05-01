package handlers

import (
	"context"
	"net/http"
	"strings"

	d "github.com/fercevik729/STLKER/eagle/data"
	"github.com/fercevik729/STLKER/octopus/data"
)

// Cache is a middleware that checks if a user cached their portfolio's profits in the past
// 15 minutes
func (c *ControlHandler) Cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		// Check the URIs to match the structs
		if strings.Contains(r.RequestURI, "/stocks/more/") {
			var co d.MoreStock
			err = c.getFromCache(r.RequestURI, &co, w)
		} else if strings.Contains(r.RequestURI, "/stocks/") {
			var s d.Stock
			err = c.getFromCache(r.RequestURI, &s, w)
		} else {
			// A single slash corresponds to a call made to GetAll
			// Two slashes corresponds to a call made to GetPortfolio
			// Three slashes corresponds to a call made to ReadSecurity
			// Get the key
			key := retrieveUsername(r) + r.RequestURI
			slashCount := strings.Count(r.RequestURI, "/")
			switch slashCount {
			// All portfolios for a user
			case 1:
				if retrieveAdmin(r) {
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
		//c.l.Println("[INFO] Using cache for key", key)
		data.ToJSON(content, w)
	}
	return err
}
