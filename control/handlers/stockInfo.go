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
	Open          float64
	High          float64
	Low           float64
	Price         float64
	Volume        float64
	LTD           string
	PrevClose     float64
	Change        float64
	PercentChange string
	Destination   string
}

// MoreStock contains important financial metrics
// swagger:model
type MoreStock struct {
	Ticker            string
	Name              string
	Exchange          string
	Sector            string
	MarketCap         float64
	PERatio           float64
	PEGRatio          float64
	DivPerShare       float64
	EPS               float64
	RevPerShare       float64
	ProfitMargin      float64
	YearHigh          float64
	YearLow           float64
	SharesOutstanding float64
	PriceToBookRatio  float64
	Beta              float64
}

// swagger:route GET /stocks/{ticker}/{currency} stocks getInfo
// Outputs a stock's financial details to the client in the requested currency
// responses:
//
//	200: stockResponse
//	400: errorResponse
//	500: errorResponse
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
		c.l.Println("[INFO] Setting cache...")
		c.setStockCache(r, &s)
	}
	data.ToJSON(&Stock{
		Symbol:        s.Symbol,
		Open:          parseFloat(s.Open),
		High:          parseFloat(s.High),
		Low:           parseFloat(s.Low),
		Price:         parseFloat(s.Price),
		Volume:        parseFloat(s.Volume),
		LTD:           s.LTD,
		PrevClose:     parseFloat(s.PrevClose),
		Change:        parseFloat(s.Change),
		PercentChange: s.PercentChange,
		Destination:   destCurr,
	}, w)

}

// swagger:route GET /stocks/more/{ticker} stocks moreInfo
// Outputs more sophisticated stock informations
// responses:
//
//	200: moreStockResponse
//	400: errorResponse
//	500: errorResponse
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
		c.l.Println("[INFO] Setting cache...")
	}
	data.ToJSON(&MoreStock{
		Ticker:            co.Ticker,
		Name:              co.Name,
		Exchange:          co.Exchange,
		Sector:            co.Sector,
		MarketCap:         parseFloat(co.MarketCap),
		PERatio:           parseFloat(co.PERatio),
		PEGRatio:          parseFloat(co.PEGRatio),
		DivPerShare:       parseFloat(co.DivPerShare),
		EPS:               parseFloat(co.EPS),
		RevPerShare:       parseFloat(co.RevPerShare),
		ProfitMargin:      parseFloat(co.ProfitMargin),
		YearHigh:          parseFloat(co.YearHigh),
		YearLow:           parseFloat(co.YearLow),
		SharesOutstanding: parseFloat(co.SharesOutstanding),
		PriceToBookRatio:  parseFloat(co.PriceToBookRatio),
		Beta:              parseFloat(co.Beta),
	}, w)
}
