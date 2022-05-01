package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// ReadEnvVar reads an environmental variable specified by key after loading vars.env
func ReadEnvVar(key string) (string, error) {
	err := godotenv.Load("vars.env")
	if err != nil {
		return "", err
	}
	return os.Getenv(key), nil
}

func (c *ControlHandler) setCache(r *http.Request, value interface{}) error {
	ctx := context.Background()
	key := retrieveUsername(r) + r.RequestURI

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
	c.l.Printf("[ERROR] %s\n", errorMsg)
	http.Error(w, fmt.Sprintf("Error: %s", errorMsg), errorCode)
}

// updateDB updates the database entry for a portfolio "port" by calling updatePrices
// and subsequently replacePortfolio
func (c *ControlHandler) updateDB(port *Portfolio) error {
	// Update prices using gRPC API
	c.updatePrices(port)
	// Delete previous portfolio and replace it with updated one
	return replacePortfolio(port.Name, port.Username, port)

}

// updateSecurities calls the grpc microservice and updates a given
// security with the new prices
func (c *ControlHandler) updateSecurities(s *Security) {
	// Get security information using Info method defined in driver.go
	st, err := Info(s.Ticker, s.Currency, c.client)
	if err != nil {
		c.l.Println("[ERROR] Couldn't get info for ticker:", s.Ticker)
		return
	}
	// Parse the stock price
	price, err := strconv.ParseFloat(st.Price, 64)
	if err != nil {
		c.l.Println("[ERROR] Couldn't parse stock price for ticker:", s.Ticker, "price:", st.Price)
		return
	}
	// Set stock price in target currency (USD by default)
	if s.Currency == "" {
		s.Currency = "USD"
	}
	c.l.Println("[DEBUG] Got price for ticker:", s.Ticker, "in", s.Currency)
	s.CurrPrice = price

	// Update the individual security's gains and percent changes
	gain, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (s.CurrPrice-s.BoughtPrice)*s.Shares), 64)
	s.setMoves(gain, fmt.Sprintf("%.2f%%", (s.CurrPrice-s.BoughtPrice)/s.BoughtPrice*100))

}

// updatePrices concurrently retrieves stock prices for all the securities
// in the portfolio by calling updateSecurities
func (c *ControlHandler) updatePrices(port *Portfolio) {
	// Concurrently retrieve stock prices
	wg := &sync.WaitGroup{}
	for _, sec := range port.Securities {
		wg.Add(1)
		go func(s *Security) {
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
