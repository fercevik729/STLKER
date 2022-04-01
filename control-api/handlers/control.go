package handlers

import (
	"log"
	"net/http"

	"github.com/fercevik729/STLKER/control-api/data"
)

// ControlHandler is a http.Handler
type ControlHandler struct {
	l   *log.Logger
	sdb *data.StockPricesDB
}

// NewControlHandler is a constructor
func NewControlHandler(log *log.Logger, s *data.StockPricesDB) *ControlHandler {
	return &ControlHandler{
		l:   log,
		sdb: s,
	}
}

func (c *ControlHandler) MoreInfo(w http.ResponseWriter, r *http.Request) {
}
