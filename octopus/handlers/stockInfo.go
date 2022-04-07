package handlers

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/fercevik729/STLKER/octopus/data"
	pb "github.com/fercevik729/STLKER/watcher-api/protos"
)

// ControlHandler is a http.Handler
type ControlHandler struct {
	l      *log.Logger
	client pb.WatcherClient
}

type StockRequest struct {
	Ticker      string `json:"ticker"`
	Destination string `json:"dest"`
}

// NewControlHandler is a constructor
func NewControlHandler(log *log.Logger, wc pb.WatcherClient) *ControlHandler {
	return &ControlHandler{
		l:      log,
		client: wc,
	}
}

func (c *ControlHandler) GetInfo(w http.ResponseWriter, r *http.Request) {

	// Retrieve and set parameters
	sr, err := getParams(r)
	if err != nil {
		c.l.Println("[ERROR]", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ticker, destCurr := sr.Ticker, sr.Destination
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
	// Get and set parameters
	sr, err := getParams(r)
	if err != nil {
		c.l.Println("[ERROR]", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	ticker := sr.Ticker
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

// getParams is a helper function to retrieve the parameters for the API's endpoints
func getParams(r *http.Request) (*StockRequest, error) {
	// Get the parameters from the request body
	params := &StockRequest{}
	data.FromJSON(params, r.Body)
	// Check if parameters were empty
	if reflect.DeepEqual(*params, StockRequest{}) {
		return nil, fmt.Errorf("must provide a ticker")
	}
	return params, nil

}
