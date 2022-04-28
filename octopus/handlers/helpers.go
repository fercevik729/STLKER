package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// logHTTPError logs the error message for a handler with the specified message and status code
func (c *ControlHandler) logHTTPError(w http.ResponseWriter, errorMsg string, errorCode int) {
	c.l.Printf("[ERROR] %s\n", errorMsg)
	http.Error(w, fmt.Sprintf("Error: %s", errorMsg), errorCode)
}

// retrieveUsername retrieves the username of the user who made the request from the request context
func (c *ControlHandler) retrieveUsername(r *http.Request) string {
	// Get username from request context
	username := r.Context().Value(Username{})
	// c.l.Println("[INFO] Got username:", username)

	v, ok := username.(string)
	if ok {
		return v
	}
	return ""
}

// retrieveAdmin retrieves a boolean value from a request's context
// to specify if the user was the admin
func (c *ControlHandler) retrieveAdmin(r *http.Request) bool {
	// Get email from request context
	isAdmin := r.Context().Value(IsAdmin{})
	// c.l.Println("[INFO] User is admin:", isAdmin)

	v, ok := isAdmin.(bool)
	if ok {
		return v
	}
	return false
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

// getSecurityVars returns a portfolioName, ticker string vars, and username from the request.
// It also logs the method being handled
func (c *ControlHandler) getSecurityVars(method string, r *http.Request) (portName string, ticker string, username string) {
	// Get URI vars
	vars := mux.Vars(r)
	portName = vars["name"]
	ticker = vars["ticker"]
	// Log endpoint
	c.l.Printf("[INFO] Handle %s for portfolio: %s, for ticker: %s", method, portName, ticker)
	return portName, ticker, c.retrieveUsername(r)
}

// replacePortfolio replaces a portfolio of name "portName" for a user "username" with
// a new portfolio struct "newPort"
func replacePortfolio(portName string, username string, newPort *Portfolio) error {
	// Declare vars
	var (
		port Portfolio
		sec  Security
	)
	// Create a new gorm db connection
	db, err := newGormDBConn(databasePath)
	if err != nil {
		return err
	}
	// Check if any results were found
	db.Where("name=?", portName).Where("username=?", username).Preload("Securities").Find(&port)
	if reflect.DeepEqual(port, &Portfolio{}) {
		return fmt.Errorf("no results could be found for portfolio %s and username %s", portName, username)
	}
	// Delete the securities and then the portfolio
	db.Model(&sec).Where("portfolio_id=?", port.ID).Delete(sec)
	db.Model(&port).Delete(&port)

	// If a new portfolio is specified create it in place of the old one
	if newPort != nil {
		db.Create(newPort)
	}

	return nil

}

// getPortfolioId returns a portfolio's id provided its name and the username associated with it
func getPortfolioId(db *sql.DB, portName string, username string) (int, error) {
	// Execute query
	rows, err := db.Query("SELECT id FROM portfolios WHERE name=? AND username=?", portName, username)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return -1, err
	}
	// Grab the id of the portfolio
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return -1, err
		}
	}
	return id, nil

}

// newGormDBConn opens a new gorm database connection
func newGormDBConn(databaseName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil

}

// newSqlDBConn opens a new sqlite3 database connection
func newSqlDBConn(databaseName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// validatePortfolio validates a portfolio's name
func validatePortfolio(port *Portfolio) bool {
	// Check the length of the name and if it contains spaces
	if len(port.Name) < 3 || len(port.Name) > 30 || strings.Contains(port.Name, " ") {
		return false
	}
	// Check if the name is alphanumeric
	re := regexp.MustCompile(`[a-zA-Z0-9]+`)
	matches := re.FindAllString(port.Name, -1)

	return len(matches) == 1

}

// validateUser validates a user
func validateUser(usr User) bool {
	// Check the lengths
	if len(usr.Username) < 6 || len(usr.Username) > 30 || len(usr.Password) < 10 || len(usr.Password) > 100 {
		return false
	}
	// Check if the username or pwd contain invalid chars or the password contains the username
	if strings.ContainsAny(usr.Username, "(){}[]|!%^@:;&_'-+<>") || strings.ContainsAny(usr.Password, "(){}[]|!%^@:;&_'-+<>") || strings.Contains(usr.Password, usr.Username) {
		return false
	}
	return true
}
