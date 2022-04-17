package handlers

import (
	"log"

	pb "github.com/fercevik729/STLKER/watcher-api/protos"
)

// ControlHandler is a http.Handler
type ControlHandler struct {
	l      *log.Logger
	client pb.WatcherClient
}

// NewControlHandler is a constructor
func NewControlHandler(log *log.Logger, wc pb.WatcherClient) *ControlHandler {
	return &ControlHandler{
		l:      log,
		client: wc,
	}
}
