package handlers

import (
	"fmt"
	"net/http"
)

func (c *ControlHandler) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		method := r.Method
		username := retrieveUsername(r)
		admin := retrieveAdmin(r)
		logMsg := fmt.Sprintf("[INFO] Handle %s request to %s", method, uri)

		if admin {
			logMsg += " for admin"
		} else if username != "" {
			logMsg += fmt.Sprint(" for user:", username)
		}
		c.l.Info(logMsg)
		next.ServeHTTP(w, r)
	})
}
