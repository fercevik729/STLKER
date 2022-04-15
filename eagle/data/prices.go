package data

import (
	"log"
	"net/http"
	"time"
)

// StockPrices contains a map of tickers and their related prices
type StockPrices struct {
	l      *log.Logger
	Prices map[string]float64
}

// NewStockPrices constructs a new StockPrices struct
func NewStockPrices(l *log.Logger) *StockPrices {
	sp := &StockPrices{
		l:      l,
		Prices: make(map[string]float64),
	}
	return sp
}

// GetInfo sends HTTP requests to the Alpha Vantage API to get stock info for the specified ticker
func (sp *StockPrices) GetInfo(ticker string) *Stock {
	sp.l.Println("[INFO] Handle GetInfo for ticker:", ticker)
	// Check if US markets are closed
	if MarketsClosed(time.Now()) {
		sp.l.Println("[WARNING] Markets are closed")
	}
	// Load the api key
	keyfile := "../key.txt"
	key, err := LoadKey(keyfile)
	if err != nil {
		sp.l.Println("[ERROR] Couldn't open key file at", keyfile)
	}
	// Get the new stock price from Alpha Vantage
	url := "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + ticker + "&apikey=" + key
	resp, err := http.Get(url)
	if err != nil {
		sp.l.Println("[ERROR] Could not reach the url")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		sp.l.Printf("[ERROR] expected http status code 200 got %d", resp.StatusCode)
	}
	// Deconstruct JSON response body to a GlobalQuote then compare
	gq := &GlobalQuote{}
	err = FromJSON(gq, resp.Body)
	if err != nil {
		sp.l.Println("[ERROR] Could not decode the response body")
		return nil
	}

	return &gq.StockData
}

// MoreInfo sends HTTP requests to the Alpha Vantage API to get the company overview for a specified ticker
func (sp *StockPrices) MoreInfo(ticker string) *MoreStock {
	sp.l.Println("[INFO] Handle MoreInfo for ticker:", ticker)

	// Warn if markets are closed
	if MarketsClosed(time.Now()) {
		sp.l.Println("[WARNING] Markets are closed")
	}
	// Load the api key
	keyfile := "../key.txt"
	key, err := LoadKey(keyfile)
	if err != nil {
		sp.l.Println("[ERROR] Couldn't open key file at", keyfile)
	}
	// Get the company overview from Alpha Vantage
	url := "https://www.alphavantage.co/query?function=OVERVIEW&symbol=" + ticker + "&apikey=" + key
	resp, err := http.Get(url)
	if err != nil {
		sp.l.Println("[ERROR] Could not reach the url")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		sp.l.Printf("[ERROR] expected http status code 200 got %d", resp.StatusCode)
	}
	// Deconstruct JSON response body to a GlobalQuote then compare
	ms := &MoreStock{}
	err = FromJSON(ms, resp.Body)
	if err != nil {
		sp.l.Println("[ERROR] Could not decode the response body")
		return nil
	}

	return ms

}

// MarketsClosed is a helper method that returns true if the markets are closed
// Currently only supported for US markets
func MarketsClosed(t time.Time) bool {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	// Get the eastern datetime
	eastDT := t.In(loc)
	// Check if day is Saturday or Sunday
	switch eastDT.Weekday() {
	case time.Saturday:
		return true
	case time.Sunday:
		return true
	}
	// Get current timestamp
	format := "15:04:05"
	currTime, _ := time.Parse(format, eastDT.Format(format))

	// Get opening and closing times as time values
	openString := "09:30:00"
	closeString := "16:00:00"

	// Parse opening and closing hour strings
	openTime, err := time.Parse(format, openString)
	if err != nil {
		panic(err)
	}
	closeTime, err := time.Parse(format, closeString)
	if err != nil {
		panic(err)
	}

	// Check if the current time is before opening hours, after closing hours,
	// or neither
	if currTime.Before(openTime) {
		return true
	} else if currTime.After(closeTime) {
		return true
	}
	return false

}
