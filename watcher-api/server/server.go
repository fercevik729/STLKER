package server

import (
	"context"
	"io"
	"log"
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
		w.l.Println("[INFO] Handle client request, ticker:", tr.Ticker, "dest currency:", tr.Destinatation.String())

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

func (w *Watcher) GetInfo(ctx context.Context, tr *protos.TickerRequest) *protos.TickerResponse {
	s := w.stockPrices.GetInfo(tr.Ticker)
	return &protos.TickerResponse{
		Symbol:        s.Symbol,
		Open:          s.Open,
		High:          s.High,
		Low:           s.Low,
		Price:         s.Price,
		Volume:        s.Volume,
		PrevClose:     s.PrevClose,
		PercentChange: s.PercentChange,
	}
}

func (w *Watcher) handleUpdates() {
	su := w.stockPrices.MonitorStocks(60 * time.Second)

	for range su {
		w.l.Println("[INFO] Got updated stock prices")
		// Loop over subscribed clients
		for k, stocks := range w.subs {
			// Loop over subbed stocks
			for _, tr := range stocks {

				// Get the stock info
				stock := w.stockPrices.GetInfo(tr.Ticker)
				// Get the price in USD
				price, err := strconv.ParseFloat(stock.Price, 64)
				if err != nil {
					w.l.Println("[ERROR] couldn't parse the stock price")
					continue
				}
				// Get the destination currency
				destC := tr.Destinatation.String()

				// Convert the price
				convPrice := convert(price, destC)
				err = k.Send(&protos.PriceResponse{StockPrice: convPrice})
				if err != nil {
					w.l.Println("[ERROR] Couldn't send the ticker response to the client")
				}
			}
		}
	}

}

func convert(original float64, dest string) float64 {

	// TODO: query the FOREX exchange endpoint from the Alpha Vantage API
	return 0
}
