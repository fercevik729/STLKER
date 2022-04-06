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

type Request struct {
	Ticker      string `json:"ticker"`
	Destination string `json:"dest"`
}

// NewControlHandler is a constructor
func NewControlHandler(log *log.Logger, s *data.StockClientDB) *ControlHandler {
	return &ControlHandler{
		l:   log,
		sdb: s,
	}
}

// setParams is a helper function to set the paramters for the API's endpoints
func setParams(r *http.Request) (string, string) {
	// Get the parameters from the request body
	params := &Request{}
	data.FromJSON(params, r.Body)

	return params.Ticker, params.Destination
}
func (c *ControlHandler) GetInfo(w http.ResponseWriter, r *http.Request) {

	ticker, destCurr := setParams(r)
	c.l.Println("[DEBUG] Handle GetInfo for", ticker, "in", destCurr)

	stock, err := c.sdb.GetInfo(ticker, destCurr)
	if err != nil {
		c.l.Println("[ERROR] Couldn't get ticker information, ensure ticker and destination currency are valid")
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	data.ToJSON(stock, w)

}
func (c *ControlHandler) MoreInfo(w http.ResponseWriter, r *http.Request) {
	ticker, _ := setParams(r)
	c.l.Println("[DEBUG] Handle MoreInfo for", ticker)

	moreStock, err := c.sdb.MoreInfo(ticker)
	if err != nil {
		c.l.Println("[ERROR] Couldn't get company overview information for ticker:", ticker, "err:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	data.ToJSON(moreStock, w)
}
