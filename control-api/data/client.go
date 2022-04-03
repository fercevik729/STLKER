package data

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/fercevik729/STLKER/watcher-api/data"
	"github.com/fercevik729/STLKER/watcher-api/protos"
	pb "github.com/fercevik729/STLKER/watcher-api/protos"
)

type StockClientDB struct {
	client protos.WatcherClient
	log    *log.Logger
	prices map[string]float64
	sub    pb.Watcher_SubscribeTickerClient
}

type StockPrice struct {
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
	Dest   string  `json:"dest"`
}

func NewStockPricesDB(client protos.WatcherClient, l *log.Logger) *StockClientDB {
	spdb := &StockClientDB{
		client: client,
		log:    l,
		prices: make(map[string]float64),
		sub:    nil,
	}

	return spdb
}

// SubscribeTicker returns a read only channel of StockPrice pointers which are then to be converted
// to JSON
func (s *StockClientDB) SubscribeTicker(ticker, destination string, stocks chan<- *StockPrice) {
	in := &pb.TickerRequest{
		Ticker:      ticker,
		Destination: pb.Currencies(pb.Currencies_value[destination]),
	}
	stream, err := s.client.SubscribeTicker(context.Background(), in)
	if err != nil {
		s.log.Println("[ERROR] unable to subscribe for stock prices, error:", err)
	}
	waitC := make(chan struct{})

	// Assign sub as the client
	s.sub = stream
	go func() {
		for {
			// Receive messages from the client until EOF is reached
			// TODO: change SubscribeTicker to be only server-side streaming
			s.log.Println("[INFO] received updated prices from server for ticker:", in.Ticker, "dest currency:", in.Destination)
			// Add the stock price to the channel
			s.client.
				stocks <- &StockPrice{
				Ticker: ticker,
				Price:  pr.StockPrice,
				Dest:   pr.Currency,
			}

			if err != nil {
				s.log.Println("[ERROR] receiving request from user error:", err)
			}

		}
	}()

	<-waitC
	stream.CloseSend()

}

// GetInfo returns a pointer to a Stock struct and an error if one arises
func (s *StockClientDB) GetInfo(ticker, destination string) (*data.Stock, error) {
	if len(ticker) > 5 {
		return nil, fmt.Errorf("ticker symbol is too long")
	}
	tr := &protos.TickerRequest{
		Ticker:      ticker,
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}
	// Have the request timeout after 15 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	stockInfo, err := s.client.GetInfo(ctx, tr)
	if err != nil {
		return nil, err
	}
	// update map of prices (cache)
	price, err := strconv.ParseFloat(stockInfo.Price, 64)
	if err != nil {
		return nil, err
	}
	s.prices[ticker] = price

	// Return a pointer to a Stock struct
	return &data.Stock{
		Symbol:        ticker,
		Open:          stockInfo.Open,
		High:          stockInfo.High,
		Low:           stockInfo.Low,
		Price:         stockInfo.Price,
		Volume:        stockInfo.Volume,
		LTD:           stockInfo.LTD,
		PrevClose:     stockInfo.PrevClose,
		PercentChange: stockInfo.PercentChange,
	}, nil
}

// MoreInfo returns a pointer to a MoreStock struct and an error if one arises
func (s *StockClientDB) MoreInfo(ticker string) (*data.MoreStock, error) {
	if len(ticker) > 5 {
		return nil, fmt.Errorf("ticker symbol is too long")
	}
	tr := &protos.TickerRequest{
		Ticker: ticker,
	}
	// Time out the request after 15 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	moreStockInfo, err := s.client.MoreInfo(ctx, tr)
	defer cancel()

	if err != nil {
		return nil, err
	}
	// Return a pointer to a Stock struct
	return &data.MoreStock{
		Ticker:            moreStockInfo.Ticker,
		Name:              moreStockInfo.Name,
		Exchange:          moreStockInfo.Exchange,
		Sector:            moreStockInfo.Sector,
		MarketCap:         moreStockInfo.MarketCap,
		PERatio:           moreStockInfo.PERatio,
		PEGRatio:          moreStockInfo.PEGRatio,
		DivPerShare:       moreStockInfo.DivPerShare,
		EPS:               moreStockInfo.EPS,
		RevPerShare:       moreStockInfo.RevPerShare,
		ProfitMargin:      moreStockInfo.ProfitMargin,
		YearHigh:          moreStockInfo.YearHigh,
		YearLow:           moreStockInfo.YearLow,
		SharesOutstanding: moreStockInfo.SharesOutstanding,
		PriceToBookRatio:  moreStockInfo.PriceToBookRatio,
		Beta:              moreStockInfo.Beta,
	}, nil
}
