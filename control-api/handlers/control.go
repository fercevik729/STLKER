package handlers

import (
	"log"
	"net/http"

	"github.com/fercevik729/STLKER/control-api/data"
)

// ControlHandler is a http.Handler
type ControlHandler struct {
	l   *log.Logger
	sdb *data.StockClientDB
}

// NewControlHandler is a constructor
func NewControlHandler(log *log.Logger, s *data.StockClientDB) *ControlHandler {
	return &ControlHandler{
		l:   log,
		sdb: s,
	}
}

func (c *ControlHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	destCurr := r.URL.Query().Get("dest")
	c.l.Println("[DEBUG] Handle GetInfo for", ticker, "in USD", destCurr)

	stock, err := c.sdb.GetInfo(ticker, destCurr)
	if err != nil {
		c.l.Println("[ERROR] Couldn't get ticker information err:", err)
	}
	data.ToJSON(stock, w)

}
func (c *ControlHandler) MoreInfo(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	c.l.Println("[DEBUG] Handle MoreInfo for", ticker)

	moreStock, err := c.sdb.MoreInfo(ticker)
	if err != nil {
		c.l.Println("[ERROR] Couldn't get company overview information for ticker:", ticker, "err:", err)
	}
	data.ToJSON(moreStock, w)
}
func (c *ControlHandler) SubscribeTicker(w http.ResponseWriter, r *http.Request) {

	// Get query parameters
	ticker := r.URL.Query().Get("ticker")
	destCurr := r.URL.Query().Get("dest")

	ch := make(chan *data.StockPrice)
	go c.sdb.SubscribeTicker(ticker, destCurr, ch)

	for stock := range ch {
		data.ToJSON(stock, w)
	}
}
