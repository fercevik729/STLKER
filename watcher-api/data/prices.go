package data

import (
	"log"
	"net/http"
	"strconv"
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

// MonitorStocks checks for updates to the specified stocks' prices
func (s *StockPrices) MonitorStocks(dur time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		ticker := time.NewTicker(dur)

		for range ticker.C {
			for k, v := range s.Prices {
				// Get the stock
				stock := s.GetInfo(k)
				p, err := strconv.ParseFloat(stock.Price, 64)
				if err != nil {
					// Only occurs if ticker is invalid, delete the faulty ticker
					s.l.Println("[ERROR] Couldn't parse the stock price for ticker:", k)
					delete(s.Prices, k)
					continue
				}

				// Update the price if it is different
				if p != v {
					s.Prices[k] = p
				}
			}
			ch <- struct{}{}
		}
	}()

	return ch
}

// GetInfo sends HTTP requests to the Alpha Vantage API to get stock info for the specified ticker
func (sp *StockPrices) GetInfo(ticker string) *Stock {
	sp.l.Println("[INFO] Handle GetInfo for ticker:", ticker)
	// Check if US markets are closed
	if MarketsClosed() {
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

// MarketsClosed is a helper method that returns true if the markets are closed
func MarketsClosed() bool {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	// Get the eastern datetime
	eastDatetime := time.Now().In(loc)
	// Check if day is Saturday or Sunday
	switch eastDatetime.Weekday() {
	case time.Saturday:
		return true
	case time.Sunday:
		return true
	}
	// Check if time is after closing hours or before opening hours
	cY, cM, cD := eastDatetime.Date()
	format := "15:04:05 MST"
	openString := "09:30:00 EDT"
	closeString := "16:00:00 EDT"

	// Parse opening and closing hour strings
	openTime, err := time.Parse(format, openString)
	if err != nil {
		panic(err)
	}
	closeTime, err := time.Parse(format, closeString)
	if err != nil {
		panic(err)
	}

	// Add current date to open and closing times
	closeTime = closeTime.Add(time.Duration(cY)).Add(time.Duration(cM)).Add(time.Duration(cD))
	openTime = openTime.Add(time.Duration(cY)).Add(time.Duration(cM)).Add(time.Duration(cD))

	// Check if the EDT time is before opening hours, after closing hours,
	// or neither
	if eastDatetime.Before(openTime) {
		return true
	} else if eastDatetime.After(openTime) {
		return true
	}
	return false

}
