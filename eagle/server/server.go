package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fercevik729/STLKER/eagle/data"
	pb "github.com/fercevik729/STLKER/eagle/protos"
)

type WatcherServer struct {
	pb.UnimplementedWatcherServer
	stockPrices *data.StockPrices
	l           *log.Logger
}

func NewWatcher(sp *data.StockPrices, l *log.Logger) *WatcherServer {
	w := &WatcherServer{
		stockPrices: sp,
		l:           l,
	}
	return w
}

// GetInfo returns a TickerResponse containing the price of the security in USD
func (w *WatcherServer) GetInfo(ctx context.Context, tr *pb.TickerRequest) (*pb.TickerResponse, error) {
	s := w.stockPrices.GetInfo(tr.Ticker)
	return &pb.TickerResponse{
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

// MoreInfo returns a CompanyResponse containing important financial ratios
func (w *WatcherServer) MoreInfo(ctx context.Context, tr *pb.TickerRequest) (*pb.CompanyResponse, error) {
	ms := w.stockPrices.MoreInfo(tr.Ticker)
	return &pb.CompanyResponse{
		Ticker:            tr.Ticker,
		Name:              ms.Name,
		Exchange:          ms.Exchange,
		Sector:            ms.Sector,
		MarketCap:         ms.MarketCap,
		PERatio:           ms.PERatio,
		PEGRatio:          ms.PEGRatio,
		DivPerShare:       ms.DivPerShare,
		EPS:               ms.EPS,
		RevPerShare:       ms.RevPerShare,
		ProfitMargin:      ms.ProfitMargin,
		YearHigh:          ms.YearHigh,
		YearLow:           ms.YearLow,
		SharesOutstanding: ms.SharesOutstanding,
		PriceToBookRatio:  ms.PriceToBookRatio,
		Beta:              ms.Beta,
	}, nil
}

// Echo returns the request that was passed to it
func (w *WatcherServer) Echo(ctx context.Context, tr *pb.TickerRequest) (*pb.TickerRequest, error) {
	return tr, nil
}

// ExchangeRates is a struct composed of ExchangeRate
// It is used to unmarshal FOREX JSON data from the Alpha Vantage API
type ExchangeRates struct {
	R ExchangeRate `json:"Realtime Currency Exchange Rate"`
}
type ExchangeRate struct {
	Rate string `json:"5. Exchange Rate"`
}

// convert calls the AV's FOREX endpoint and converts the original stock price as needed
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
