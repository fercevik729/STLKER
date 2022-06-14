package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/fercevik729/STLKER/grpc/data"
	pb "github.com/fercevik729/STLKER/grpc/protos"
)

// Info returns a pointer to a Stock struct and an error if one arises
func Info(ticker, destination string, client pb.WatcherClient) (*data.Stock, error) {
	if len(ticker) > 5 {
		return nil, fmt.Errorf("ticker symbol is too long")
	}
	if len(ticker) == 0 {
		return nil, fmt.Errorf("ticker symbol is too short")
	}
	tr := &pb.TickerRequest{
		Ticker:      ticker,
		Destination: pb.Currencies(pb.Currencies_value[destination]),
	}
	// Have the request timeout after 15 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	stockInfo, err := client.GetInfo(ctx, tr)
	if err != nil {
		return nil, err
	}
	if stockInfo.Open == "" {
		return nil, fmt.Errorf("couldn't find a security with symbol %s", ticker)
	}
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
func CompanyOverview(ticker string, client pb.WatcherClient) (*data.MoreStock, error) {
	if len(ticker) > 5 {
		return nil, fmt.Errorf("ticker symbol is too long")
	}
	tr := &pb.TickerRequest{
		Ticker: ticker,
	}
	// Time out the request after 15 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	moreStockInfo, err := client.MoreInfo(ctx, tr)
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
