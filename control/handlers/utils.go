package handlers

import (
	"context"
	"fmt"
	m "github.com/fercevik729/STLKER/control/models"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/gorilla/mux"
)

func (c *ControlHandler) setStockCache(r *http.Request, value interface{}) error {
	ctx := context.Background()
	key := r.RequestURI

	if err := c.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   15 * time.Minute,
	}); err != nil {
		return err
	}
	return nil
}

// logHTTPError logs the error message for a handler with the specified message and status code
func (c *ControlHandler) logHTTPError(w http.ResponseWriter, errorMsg string, errorCode int) {
	c.l.Error(errorMsg)
	http.Error(w, fmt.Sprintf("Error: %s", errorMsg), errorCode)
}

// updateSecurities calls the grpc microservice and updates a given
// security with the new prices
func (c *ControlHandler) updateSecurities(s *m.Security) {
	// Get security information using Info method defined in driver.go
	st, err := Info(s.Ticker, s.Currency, c.client)
	if err != nil {
		c.l.Error("Couldn't get info for ticker:", s.Ticker)
		return
	}
	// Parse the stock price
	price, err := strconv.ParseFloat(st.Price, 64)
	if err != nil {
		c.l.Error("Couldn't parse stock price for:", "ticker", s.Ticker, "price", st.Price)
		return
	}
	// Set stock price in target currency (USD by default)
	if s.Currency == "" {
		s.Currency = "USD"
	}
	c.l.Debug("Got price for:", "ticker", s.Ticker, "currency", s.Currency)
	s.CurrPrice = price

	// Update the individual security's gains and percent changes
	gain, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (s.CurrPrice-s.BoughtPrice)*s.Shares), 64)
	s.SetMoves(gain, fmt.Sprintf("%.2f%%", (s.CurrPrice-s.BoughtPrice)/s.BoughtPrice*100))

}

// updatePrices concurrently retrieves stock prices for all the securities
// in the portfolio by calling updateSecurities
func (c *ControlHandler) updatePrices(port *m.Portfolio) {
	// Concurrently retrieve stock prices
	wg := &sync.WaitGroup{}
	for _, sec := range port.Securities {
		wg.Add(1)
		go func(s *m.Security) {
			c.updateSecurities(s)
			wg.Done()
		}(sec)
	}
	wg.Wait()

}

// retrieveUsername retrieves the username of the user who made the request from the request context
func retrieveUsername(r *http.Request) string {
	// Get username from request context
	username := r.Context().Value(Username{})

	v, ok := username.(string)
	if ok {
		return v
	}
	return ""
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

// retrieveAdmin retrieves a boolean value from a request's context
// to specify if the user was the admin
func retrieveAdmin(r *http.Request) bool {
	// Get email from request context
	isAdmin := r.Context().Value(IsAdmin{})
	// c.l.Println("[INFO] User is admin:", isAdmin)

	v, ok := isAdmin.(bool)
	if ok {
		return v
	}
	return false
}

// retrieveSecurityVars returns a portfolioName, ticker string vars, and username from the request.
// It also logs the method being handled
func retrieveSecurityVars(r *http.Request) (portName string, ticker string, username string) {
	// Get URI vars
	vars := mux.Vars(r)
	portName = vars["name"]
	ticker = vars["ticker"]
	// Log endpoint
	return portName, ticker, retrieveUsername(r)
}
