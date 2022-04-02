package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fercevik729/STLKER/watcher-api/data"
	"github.com/fercevik729/STLKER/watcher-api/protos"
)

type Watcher struct {
	protos.UnimplementedWatcherServer
	stockPrices *data.StockPrices
	l           *log.Logger
	subs        map[protos.Watcher_SubscribeTickerServer][]*protos.TickerRequest
}

func NewWatcher(sp *data.StockPrices, l *log.Logger) *Watcher {
	w := &Watcher{
		stockPrices: sp,
		l:           l,
		subs:        make(map[protos.Watcher_SubscribeTickerServer][]*protos.TickerRequest),
	}
	// Create a goroutine to handle any updates
	go w.handleUpdates()
	return w
}

// SubscribeTicker awaits for TickerRequests from a client and stores them in a map
func (w *Watcher) SubscribeTicker(src protos.Watcher_SubscribeTickerServer) error {
	// Handles messages from the client
	for {
		tr, err := src.Recv()
		if err == io.EOF {
			w.l.Println("[INFO] Client has closed connection")
			return err
		}
		if err != nil {
			w.l.Println("[ERROR] Unable to read from client, err:", err)
			return err
		}
		w.l.Println("[INFO] Handle SubscribeTicker client request, ticker:", tr.Ticker, "dest currency:", tr.Destination.String())

		// Check to see if the client has already subscribed
		// then append the new ticker request to the slice of ticker requests
		trs, ok := w.subs[src]
		if !ok {
			trs = []*protos.TickerRequest{}
		}
		trs = append(trs, tr)
		w.subs[src] = trs
	}
}

// GetInfo returns a TickerResponse containing the price of the security in USD
func (w *Watcher) GetInfo(ctx context.Context, tr *protos.TickerRequest) (*protos.TickerResponse, error) {
	s := w.stockPrices.GetInfo(tr.Ticker)
	return &protos.TickerResponse{
		Symbol:        s.Symbol,
		Open:          s.Open,
		High:          s.High,
		Low:           s.Low,
		Price:         s.Price,
		Volume:        s.Volume,
		LTD:           s.LTD,
		PrevClose:     s.PrevClose,
		PercentChange: s.PercentChange,
	}, nil
}

// handleUpdates is a helper method that is called to concurrently send the updated prices
func (w *Watcher) handleUpdates() {
	su := w.stockPrices.MonitorStocks(60 * time.Second)

	for range su {
		// Loop over subscribed clients
		for k, stocks := range w.subs {
			// Loop over subbed stocks
			for i, tr := range stocks {
				w.l.Println("[INFO] Got updated stock prices")

				// Get the stock info
				stock := w.stockPrices.GetInfo(tr.Ticker)
				// Get the price in USD
				// If the stock price is nonexistent this means that the ticker was faulty and the program will
				// attempt to remove it from the slice associated with the client
				if stock == nil {
					w.l.Println("[WARNING] no results for ticker:", tr.Ticker)
					w.l.Println("[WARNING] removing faulty ticker:", tr.Ticker)
					w.subs[k] = append(w.subs[k][:i], w.subs[k][i+1:]...)
					continue
				}
				price, err := strconv.ParseFloat(stock.Price, 64)
				if err != nil {
					w.l.Println("[ERROR] couldn't parse the stock price")
					continue
				}
				// Get the destination currency
				destC := tr.Destination.String()

				// Convert the price
				convPrice, err := convert(price, destC)
				if err != nil {
					w.l.Println("[ERROR] Couldn't convert the price to the destination currency, err:", err)
				}
				// Send the PriceResponse back to the correct client
				err = k.Send(&protos.PriceResponse{Ticker: tr.Ticker, StockPrice: convPrice, Currency: destC})
				if err != nil {
					// Client will have closed their connection so their subscriptions should be removed
					w.l.Println("[ERROR] Couldn't send the ticker response to the client")
					delete(w.subs, k)
				}
			}
			if data.MarketsClosed(time.Now()) {
				w.l.Println("[WARNING] Subscriptions will be terminated")
				// Clear stock prices cache
				w.stockPrices.Prices = map[string]float64{}
				// Clear subscriptions
				w.subs = map[protos.Watcher_SubscribeTickerServer][]*protos.TickerRequest{}
			}
		}
	}

}

// ExchangeRates is a struct composed of ExchangeRate
// It is used to unmarshal FOREX JSON data from the Alpha Vantage API
type ExchangeRates struct {
	R ExchangeRate `json:"Realtime Currency Exchange Rate"`
}
type ExchangeRate struct {
	Rate string `json:"5. Exchange Rate"`
}

func convert(original float64, dest string) (float64, error) {
	// If a destination currency is USD simply return the original stock price, which was already in USD
	if dest == "USD" {
		return original, nil
	}
	// Load the API key
	key, err := data.LoadKey("../key.txt")
	if err != nil {
		return -1, err
	}
	// Send a request to the Alpha Vantage API
	resp, err := http.Get("https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=USD&to_currency=" + dest + "&apikey=" + key)
	if err != nil {
		return -1, err
	}
	// Check for errors and expected status code
	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
	// Convert the JSON body to a ExchangeRate struct
	er := &ExchangeRates{}
	err = data.FromJSON(er, resp.Body)
	if err != nil {
		return -1, err
	}

	// Convert the rate to a float
	newRate, err := strconv.ParseFloat(er.R.Rate, 64)
	if err != nil {
		return -1, err
	}
	// Return the stock price in the destination currency
	newPrice := newRate * original

	// Round the price to 2 decimal places
	roundedPrice, err := strconv.ParseFloat(fmt.Sprintf("%.2f", newPrice), 64)
	if err != nil {
		return newPrice, err
	}
	return roundedPrice, nil
}
