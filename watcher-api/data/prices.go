package data

import (
	"log"
	"net/http"
	"time"
)

type StockPrices struct {
	l      *log.Logger
	prices map[string]float64
}

func NewPrices(l *log.Logger) *StockPrices {
	sp := &StockPrices{
		l:      l,
		prices: make(map[string]float64),
	}
	return sp
}

func (s *StockPrices) MonitorStocks(dur time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	// Load the api key
	keyfile := "../key.txt"
	key, err := LoadKey(keyfile)
	if err != nil {
		s.l.Println("[ERROR] Couldn't open key file at:", keyfile)
		return nil
	}

	go func() {
		ticker := time.NewTicker(dur)

		for range ticker.C {
			for k, v := range s.prices {
				// Send a request to get the new stock price
				resp, err := http.Get("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + k + "&apikey=" + key)
				if err != nil {
					s.l.Println("[ERROR] Couldn't update stock price for ticker:", k)
					continue
				}
				// Deconstruct JSON response body to a BasicStock then compare
				bs := &BasicStock{}
				err = FromJSON(bs, resp.Body)
				if err != nil {
					s.l.Println("[ERROR] Couldn't read JSON body for requested ticker:", k)
				}
				// Update the price if it is different
				if bs.Price != v {
					s.prices[k] = bs.Price
				}

			}
			ch <- struct{}{}
		}
	}()

	return ch
}
