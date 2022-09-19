package handlers

import (
	"fmt"
	"net/http"

	"github.com/fercevik729/STLKER/control/data"
	"github.com/gorilla/mux"
)

// Stock is the struct equivalent to the body returned by the gRPC API
// swagger:model
type Stock struct {
	Symbol        string
	Open          string
	High          string
	Low           string
	Price         string
	Volume        string
	LTD           string
	PrevClose     string
	Change        string
	PercentChange string
}

// MoreStock contains important financial metrics
// swagger:model
type MoreStock struct {
	Ticker            string
	Name              string
	Exchange          string
	Sector            string
	MarketCap         string
	PERatio           string
	PEGRatio          string
	DivPerShare       string
	EPS               string
	RevPerShare       string
	ProfitMargin      string
	YearHigh          string
	YearLow           string
	SharesOutstanding string
	PriceToBookRatio  string
	Beta              string
}

// swagger:route GET /stocks/{ticker}/{currency} stocks getInfo
// Outputs a stock's financial details to the client in the requested currency
// responses:
//  200: stockResponse
//  400: errorResponse
//  500: errorResponse
func (c *ControlHandler) GetInfo(w http.ResponseWriter, r *http.Request) {

	// Retrieve URI variables
	vars := mux.Vars(r)
	ticker := vars["ticker"]
	destCurr := vars["currency"]

	// Get the stock information
	s, err := Info(ticker, destCurr, c.client)
	if err != nil {
		c.logHTTPError(w, "couldn't get ticker information, ensure ticker and destination currency are valid", http.StatusBadRequest)
		return
	}
	// Write the data to the client
	w.Header().Set("Content-Type", "application/json")
	if c.cache != nil {
		c.setStockCache(r, &s)
	}
	data.ToJSON(&Stock{
		Symbol:        s.Symbol,
		Open:          s.Open,
		High:          s.High,
		Low:           s.Low,
		Price:         s.Price,
		Volume:        s.Volume,
		LTD:           s.LTD,
		PrevClose:     s.PrevClose,
		Change:        s.Change,
		PercentChange: s.PercentChange,
	}, w)

}

// swagger:route GET /stocks/more/{ticker} stocks moreInfo
// Outputs more sophisticated stock informations
// responses:
//  200: moreStockResponse
//  400: errorResponse
//  500: errorResponse
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
	data.ToJSON(&MoreStock{
		Ticker:            co.Ticker,
		Name:              co.Name,
		Exchange:          co.Exchange,
		Sector:            co.Sector,
		MarketCap:         co.MarketCap,
		PERatio:           co.PERatio,
		PEGRatio:          co.PEGRatio,
		DivPerShare:       co.DivPerShare,
		EPS:               co.EPS,
		RevPerShare:       co.RevPerShare,
		ProfitMargin:      co.ProfitMargin,
		YearHigh:          co.YearHigh,
		YearLow:           co.YearLow,
		SharesOutstanding: co.SharesOutstanding,
		PriceToBookRatio:  co.PriceToBookRatio,
		Beta:              co.Beta,
	}, w)
}
