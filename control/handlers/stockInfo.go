package handlers

import (
	"fmt"
	"net/http"

	"github.com/fercevik729/STLKER/control/data"
	"github.com/gorilla/mux"
)

func (c *ControlHandler) GetInfo(w http.ResponseWriter, r *http.Request) {

	// Retrieve URI variables
	vars := mux.Vars(r)
	ticker := vars["ticker"]
	destCurr := vars["currency"]

	// Get the stock information
	stock, err := Info(ticker, destCurr, c.client)
	if err != nil {
		c.logHTTPError(w, "couldn't get ticker information, ensure ticker and destination currency are valid", http.StatusBadRequest)
		return
	}
	// Write the data to the client
	w.Header().Set("Content-Type", "application/json")
	if c.cache != nil {
		c.setStockCache(r, &stock)
	}
	data.ToJSON(stock, w)

}
func (c *ControlHandler) MoreInfo(w http.ResponseWriter, r *http.Request) {
	// Retrieve URI variable
	vars := mux.Vars(r)
	ticker := vars["ticker"]

	// Get the company overview
	co, err := CompanyOverview(ticker, c.client)
	if err != nil {
		c.logHTTPError(w, fmt.Sprint("couldn't get company overview information for ticker:", ticker, "err:", err), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if c.cache != nil {
		c.setStockCache(r, &co)
	}
	data.ToJSON(co, w)
}
