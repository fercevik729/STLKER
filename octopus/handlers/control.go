package handlers

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"

	"github.com/fercevik729/STLKER/octopus/data"
	pb "github.com/fercevik729/STLKER/watcher-api/protos"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ControlHandler is a http.Handler
type ControlHandler struct {
	l      *log.Logger
	client pb.WatcherClient
}

type StockRequest struct {
	Ticker      string `json:"ticker"`
	Destination string `json:"dest"`
}

type Stock struct {
	StockRequest
	Price float64
}

// A Portfolio is a GORM model that is a slice of Stock structs
type Portfolio struct {
	gorm.Model
	Name   string
	Stocks []*Stock

	/*
		{
		Portfolio 1:
			[
				{
					Ticker: TSLA,
					Price: 12223,
					Dest: USD

				},
				{

				}
			]
		}

	*/
}

// NewControlHandler is a constructor
func NewControlHandler(log *log.Logger, wc pb.WatcherClient) *ControlHandler {
	return &ControlHandler{
		l:      log,
		client: wc,
	}
}

func (c *ControlHandler) SavePortfolio(w http.ResponseWriter, r *http.Request) {
	// Open sqlite db connection
	db, err := gorm.Open(sqlite.Open("portfolios.db"), &gorm.Config{})
	if err != nil {
		c.l.Println("[ERROR] Couldn't connect to database")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Get generic sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		c.l.Println("[ERROR] Couldn't create sqlDB instance:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Migrate schema
	db.AutoMigrate(&Portfolio{})

	// Retrieve the portfolio from the request body
	portfolio := getPortfolioParams(r)

	port := Portfolio{}
	// Check if a portfolio with that name already exists
	db.First(&port, 1)

	// If a portfolio with that name does exist return an error
	if !reflect.DeepEqual(port, Portfolio{}) {
		c.l.Println("[ERROR] A portfolio with that name already exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Retrieve prices for all stocks in the portfolio
	c.l.Println("[INFO] Retrieving updated stock prices")

	// Concurrently retrieve stock prices
	var wg *sync.WaitGroup
	for _, stock := range port.Stocks {
		wg.Add(1)
		go func(s *Stock) {
			st, err := Info(s.Ticker, s.Destination, c.client)
			if err != nil {
				c.l.Println("[ERROR] Couldn't get info for ticker:", s.Ticker)
				return
			}
			// Parse the stock price
			price, err := strconv.ParseFloat(st.Price, 64)
			if err != nil {
				c.l.Println("[ERROR] Couldn't parse stock price for ticker:", s.Ticker)
				return
			}
			// Set stock price
			c.l.Println("[DEBUG] Got price for ticker:", s.Ticker)
			s.Price = price
			wg.Done()
		}(stock)
	}
	wg.Wait()

	// Create portfolio entry
	db.Create(portfolio)
	c.l.Println("[DEBUG] Created portfolio named", portfolio.Name)

	// Close database connection
	sqlDB.Close()

}

func (c *ControlHandler) GetPortfolio(w http.ResponseWriter, r *http.Request) {

}

func (c *ControlHandler) GetInfo(w http.ResponseWriter, r *http.Request) {

	// Retrieve and set parameters
	sr, err := getParams(r)
	if err != nil {
		c.l.Println("[ERROR]", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ticker, destCurr := sr.Ticker, sr.Destination
	c.l.Println("[DEBUG] Handle GetInfo for", ticker, "in", destCurr)

	// Get the stock information
	stock, err := Info(ticker, destCurr, c.client)
	if err != nil {
		c.l.Println("[ERROR] Couldn't get ticker information, ensure ticker and destination currency are valid")
		w.WriteHeader(http.StatusBadRequest)
	}
	// Write the data to the client
	w.Header().Set("Content-Type", "application/json")
	data.ToJSON(stock, w)

}
func (c *ControlHandler) MoreInfo(w http.ResponseWriter, r *http.Request) {
	// Get and set parameters
	sr, err := getParams(r)
	if err != nil {
		c.l.Println("[ERROR]", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	ticker := sr.Ticker
	c.l.Println("[DEBUG] Handle MoreInfo for", ticker)

	// Get the company overview
	co, err := CompanyOverview(ticker, c.client)
	if err != nil {
		c.l.Println("[ERROR] Couldn't get company overview information for ticker:", ticker, "err:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data.ToJSON(co, w)
}

// getParams is a helper function to retrieve the parameters for the API's endpoints
func getParams(r *http.Request) (*StockRequest, error) {
	// Get the parameters from the request body
	params := &StockRequest{}
	data.FromJSON(params, r.Body)
	// Check if parameters were empty
	if reflect.DeepEqual(*params, StockRequest{}) {
		return nil, fmt.Errorf("must provide a ticker")
	}
	return params, nil

}

// getPortfolioParams is a helper function to retrieve the individual stocks in a call to SavePortfolio
func getPortfolioParams(r *http.Request) *Portfolio {
	// Initialize portfolio
	port := &Portfolio{}
	data.FromJSON(port, r.Body)
	return port
}
