package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fercevik729/STLKER/grpc/data"
	pb "github.com/fercevik729/STLKER/grpc/protos"
)

type WatcherServer struct {
	pb.UnimplementedWatcherServer
	stockPrices       *data.StockPrices
	l                 *log.Logger
	subscribedTickers map[pb.Watcher_SubscribeTickerServer][]*pb.TickerRequest
}

func NewWatcher(sp *data.StockPrices, l *log.Logger) *WatcherServer {
	w := &WatcherServer{
		stockPrices:       sp,
		l:                 l,
		subscribedTickers: make(map[pb.Watcher_SubscribeTickerServer][]*pb.TickerRequest),
	}
	go w.handleUpdates()
	return w
}

// GetInfo returns a TickerResponse containing the price of the security in USD
func (w *WatcherServer) GetInfo(ctx context.Context, tr *pb.TickerRequest) (*pb.TickerResponse, error) {
	s := w.stockPrices.GetInfo(tr.Ticker)
	oldPrice, err := strconv.ParseFloat(s.Price, 64)
	if err != nil {
		return nil, err
	}
	newPrice, err := convert(oldPrice, tr.Destination.String())
	if err != nil {
		return nil, err
	}
	s.Price = fmt.Sprintf("%.2f", newPrice)
	return &pb.TickerResponse{
		Symbol:        s.Symbol,
		Open:          s.Open,
		High:          s.High,
		Low:           s.Low,
		Price:         s.Price,
		Destination:   tr.Destination.String(),
		Volume:        s.Volume,
		LTD:           s.LTD,
		PrevClose:     s.PrevClose,
		PercentChange: s.PercentChange,
	}, nil
}

// Helper function that handles updates to the requested tickers
func (w *WatcherServer) handleUpdates() {
	ru := make(chan struct{})

	// waits 15 seconds
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		for range ticker.C {
			ru <- struct{}{}
		}
	}()

	for range ru {
		w.l.Println("[INFO] Got updated stock information")
		// Loop over subscribed clients
		for k, v := range w.subscribedTickers {
			// loop over subscribed rates
			for _, tr := range v {
				resp, err := w.GetInfo(context.Background(), tr)
				if err != nil {
					w.l.Println("[ERROR] Unable to get updated stock information for", "ticker:", tr.Ticker)
				}
				err = k.Send(resp)
				if err != nil {
					w.l.Println("[ERROR] Couldn't send the updated rate")
				}
			}
		}
	}
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

// SubscribeTicker returns a stream of TickerResponses
func (w *WatcherServer) SubscribeTicker(src pb.Watcher_SubscribeTickerServer) error {
	// Receives ticker requests from the client and appends them to their list in the map
	for {
		tr, err := src.Recv()
		if err == io.EOF {
			w.l.Println("[INFO] Client has closed connection")
			return err
		}
		if err != nil {
			w.l.Println("[ERROR] Unable to read from client", "err", err)
			return err
		}
		w.l.Println("[INFO] Handle client request", "ticker:", tr.Ticker, "dest:", tr.Destination)
		trs, ok := w.subscribedTickers[src]
		if !ok {
			trs = []*pb.TickerRequest{}
		}
		trs = append(trs, tr)
		w.subscribedTickers[src] = trs
	}
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
