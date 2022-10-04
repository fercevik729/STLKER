package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/fercevik729/STLKER/control/data"
	d "github.com/fercevik729/STLKER/grpc/data"
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
			c.l.Println("[INFO] Getting from cache...")
		} else if strings.Contains(r.RequestURI, "/stocks/") {
			var s d.Stock
			err = c.getFromCache(r.RequestURI, &s, w)
			c.l.Println("[INFO] Getting from cache...")
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
