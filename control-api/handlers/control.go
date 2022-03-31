package handlers

import (
	"log"
	"net/http"
)

// ControlHandler is a http.Handler
type ControlHandler struct {
	l *log.Logger
}

// NewControlHandler is a constructor
func NewControlHandler(log *log.Logger) *ControlHandler {
	return &ControlHandler{
		l: log,
	}
}

// MoreInfo
func (c *ControlHandler) MoreInfo(w http.ResponseWriter, r *http.Request) {

	// Get the stock ticker info from the URI
	
}
