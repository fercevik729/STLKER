package handlers

import (
	"log"
	"net/http"

	"github.com/fercevik729/STLKER/octopus/data"
	pb "github.com/fercevik729/STLKER/watcher-api/protos"
	"github.com/gorilla/mux"
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

func (c *ControlHandler) GetInfo(w http.ResponseWriter, r *http.Request) {

	// Retrieve URI variables
	vars := mux.Vars(r)
	ticker := vars["ticker"]
	destCurr := vars["currency"]

	c.l.Println("[DEBUG] Handle GetInfo for", ticker, "in", destCurr)

	// Get the stock information
	stock, err := Info(ticker, destCurr, c.client)
	if err != nil {
		c.l.Println("[ERROR] Couldn't get ticker information, ensure ticker and destination currency are valid")
		w.WriteHeader(http.StatusBadRequest)
	}
	// Write the data to the client
	w.Header().Set("Content-Type", "application/json")
	data.ToJSON(stock, w)

}
func (c *ControlHandler) MoreInfo(w http.ResponseWriter, r *http.Request) {
	// Retrieve URI variable
	vars := mux.Vars(r)
	ticker := vars["ticker"]
	c.l.Println("[DEBUG] Handle MoreInfo for", ticker)

	// Get the company overview
	co, err := CompanyOverview(ticker, c.client)
	if err != nil {
		c.l.Println("[ERROR] Couldn't get company overview information for ticker:", ticker, "err:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data.ToJSON(co, w)
}
