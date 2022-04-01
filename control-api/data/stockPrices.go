package data

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/fercevik729/STLKER/watcher-api/data"
	"github.com/fercevik729/STLKER/watcher-api/protos"
)

type StockPricesDB struct {
	watcher protos.WatcherClient
	log     *log.Logger
	prices  map[string]float64
	client  protos.Watcher_SubscribeTickerClient
}

func NewStockPricesDB(w protos.WatcherClient, l *log.Logger) *StockPricesDB {
	spdb := &StockPricesDB{
		watcher: w,
		log:     l,
		prices:  make(map[string]float64),
		client:  nil,
	}

	go spdb.handleUpdates()
	return spdb
}

func (s *StockPricesDB) handleUpdates() {
	sub, err := s.watcher.SubscribeTicker(context.Background())
	if err != nil {
		s.log.Println("[ERROR] unable to subscribe for stock prices, error:", err)
	}
	// Assign sub as the client
	s.client = sub
	for {
		// Receive messages from the client
		pr, err := sub.Recv()
		s.log.Println("[INFO] received updatd prices from server for ticker:", pr.Ticker, "dest currency:", pr.Currency)
		// Store the ticker's price
		s.prices[pr.Ticker] = pr.StockPrice

		if err != nil {
			s.log.Println("[ERROR] receiving request from user error:", err)
		}

	}

}

func (s *StockPricesDB) getInfo(ticker, destination string) (*data.Stock, error) {
	if len(ticker) > 5 {
		return nil, fmt.Errorf("ticker symbol is too long")
	}
	tr := &protos.TickerRequest{
		Ticker:      ticker,
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}
	resp, err := s.watcher.GetInfo(context.Background(), tr)
	if err != nil {
		return nil, err
	}
	// update list of prices
	price, err := strconv.ParseFloat(resp.Price, 64)
	if err != nil {
		return nil, err
	}
	s.prices[ticker] = price

	// Return a pointer to a Stock struct
	return &data.Stock{
		Symbol:        ticker,
		Open:          resp.Open,
		High:          resp.High,
		Low:           resp.Low,
		Price:         resp.Price,
		Volume:        resp.Volume,
		LTD:           resp.LTD,
		PrevClose:     resp.PrevClose,
		PercentChange: resp.PercentChange,
	}, nil
}
