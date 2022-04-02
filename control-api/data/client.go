package data

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/fercevik729/STLKER/watcher-api/data"
	"github.com/fercevik729/STLKER/watcher-api/protos"
)

type StockClientDB struct {
	client protos.WatcherClient
	log    *log.Logger
	prices map[string]float64
	sub    protos.Watcher_SubscribeTickerClient
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

func (s *StockClientDB) SubscribeTicker(ticker, destination string) {
	// TODO: rewrite this to make it compatible with RESTful frontend
	stream, err := s.client.SubscribeTicker(context.Background())
	if err != nil {
		s.log.Println("[ERROR] unable to subscribe for stock prices, error:", err)
	}
	waitC := make(chan struct{})
	// Assign sub as the client
	s.sub = stream
	go func() {
		for {
			// Receive messages from the client
			pr, err := stream.Recv()
			if err == io.EOF {
				close(waitC)

			}
			s.log.Println("[INFO] received updatd prices from server for ticker:", pr.Ticker, "dest currency:", pr.Currency)
			// Store the ticker's price
			s.prices[pr.Ticker] = pr.StockPrice

			if err != nil {
				s.log.Println("[ERROR] receiving request from user error:", err)
			}

		}
	}()

	<-waitC
	stream.CloseSend()

}

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
