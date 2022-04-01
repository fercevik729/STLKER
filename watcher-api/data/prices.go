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
	prices map[string]float64
}

// NewStockPrices constructs a new StockPrices struct
func NewStockPrices(l *log.Logger) *StockPrices {
	sp := &StockPrices{
		l:      l,
		prices: make(map[string]float64),
	}
	return sp
}

// MonitorStocks checks for updates to the specified stocks' prices
func (s *StockPrices) MonitorStocks(dur time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		ticker := time.NewTicker(dur)

		for range ticker.C {
			for k, v := range s.prices {

				// Get the stock
				stock := s.GetInfo(k)
				p, err := strconv.ParseFloat(stock.Price, 64)
				if err != nil {
					// Shouldn't happen
					s.l.Println("[ERROR] Couldn't parse the stock price")
				}

				// Update the price if it is different
				if p != v {
					s.prices[k] = p
				}
			}
			ch <- struct{}{}
		}
	}()

	return ch
}

// GetInfo sends HTTP requests to the Alpha Vantage API to get stock info for the specified ticker
func (sp *StockPrices) GetInfo(ticker string) *Stock {
	// Load the api key
	keyfile := "../key.txt"
	key, err := LoadKey(keyfile)
	if err != nil {
		sp.l.Println("[ERROR] Couldn't open key file at", keyfile)
	}
	// Get the new stock price from Alpha Vantage
	resp, err := http.Get("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + ticker + "&apikey=" + key)
	if err != nil {
		sp.l.Println("[ERROR] Could not reach the url")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		sp.l.Printf("[ERROR] expected http status code 200 got %d", resp.StatusCode)
	}
	// Deconstruct JSON response body to a BasicStock then compare
	bs := &Stock{}
	err = FromJSON(bs, resp.Body)
	if err != nil {
		sp.l.Println("[ERROR] Could not decode the response body")
		return nil
	}

	return bs
}
