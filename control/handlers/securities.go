package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fercevik729/STLKER/control/data"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Security struct {
	STLKERModel
	SecurityID  int     `gorm:"primary_key" json:"-"`
	Ticker      string  `json:"Ticker"`
	BoughtPrice float64 `json:"Bought Price"`
	CurrPrice   float64 `json:"Current Price"`
	Shares      float64 `json:"Shares"`
	Gain        float64 `json:"Gain"`
	Change      string  `json:"Percent Change"`
	// Currency is the destination currency of the stock
	Currency string `json:"Currency" gorm:"default:USD"`
	// Foreign key
	PortfolioID uint `json:"-"`
}

// Used for destructuring POST and PUT data
type securityData struct {
	Ticker string  `json:"Ticker"`
	Shares float64 `json:"Shares"`
}

// setMoves sets the gain and change variables of s to the new parameters
func (s *Security) setMoves(gain float64, change string) {
	s.Gain = gain
	s.Change = change
}

func (c *ControlHandler) CreateSecurity(w http.ResponseWriter, r *http.Request) {
	// Get URI vars
	portName := mux.Vars(r)["name"]
	username := retrieveUsername(r)

	// Get ticker and shares info from JSON body
	var params securityData
	data.FromJSON(&params, r.Body)

	ticker := params.Ticker
	shares := params.Shares

	// Create sql db instance
	db, err := newGormDBConn(c.dbName)
	if err != nil {
		c.logHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}
	// Get portfolio id
	portId := getPortfolioId(db, portName, username)
	if portId == 0 {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get updated price for %s", portName), http.StatusBadRequest)
		return
	}
	// Get stock info
	stock, err := Info(ticker, "USD", c.client)
	if err != nil {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get updated price for %s", ticker), http.StatusBadRequest)
		return
	}
	price, _ := strconv.ParseFloat(stock.Price, 64)
	// Create the security struct
	newSecurity := Security{
		Ticker:      ticker,
		BoughtPrice: price,
		CurrPrice:   price,
		Shares:      shares,
		Currency:    "USD",
		PortfolioID: uint(portId),
	}
	db.Debug().Create(&newSecurity)
}

func (c *ControlHandler) ReadSecurity(w http.ResponseWriter, r *http.Request) {
	// Get URI vars
	portName, ticker, username := retrieveSecurityVars(r)
	db, err := newGormDBConn(c.dbName)
	if err != nil {
		c.logHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}
	// Get portfolio_id
	portId := getPortfolioId(db, portName, username)
	if portId == 0 {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get updated price for %s", portName), http.StatusBadRequest)
		return
	}
	var security Security
	db.Model(&Security{}).Select([]string{"ticker", "bought_price", "curr_price", "shares", "gain", "change"}).Where("ticker=?", ticker).Where("portfolio_id=?", portId).First(&security)
	// Update the security
	c.updateSecurities(&security)

	// Set the cache
	if c.cache != nil {
		err = c.setCache(r, &security)
	}
	if err != nil {
		c.logHTTPError(w, "Couldn't set value into cache", http.StatusInternalServerError)
	}
	// Write to responsewriter
	data.ToJSON(&security, w)
}

func (c *ControlHandler) UpdateSecurity(w http.ResponseWriter, r *http.Request) {
	// Get request vars
	portName := mux.Vars(r)["name"]
	username := retrieveUsername(r)
	var sd securityData
	data.FromJSON(&sd, r.Body)

	// Check if the payload is empty
	if sd == (securityData{}) {
		c.logHTTPError(w, "Bad request payload", http.StatusBadRequest)
		return
	}

	// Create sql db instance
	db, err := newGormDBConn(c.dbName)
	if err != nil {
		c.logHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}
	// Get portfolio id
	portId := getPortfolioId(db, portName, username)
	if portId == 0 {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get portfolio id for name: %s", portName), http.StatusBadRequest)
		return
	}
	// Update the portfolio
	db.Model(&Security{}).Where("portfolio_id=?", portId).Where("ticker=?", sd.Ticker).Update("shares", sd.Shares)
	data.ToJSON(ResponseMessage{Msg: fmt.Sprintf("Updated security of ticker %s", sd.Ticker)},
		w,
	)

}

func (c *ControlHandler) DeleteSecurity(w http.ResponseWriter, r *http.Request) {
	portName, ticker, username := retrieveSecurityVars(r)
	// Connect to database
	db, err := newGormDBConn(c.dbName)
	if err != nil {
		c.logHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}
	// Get portfolio with the name specified by the mux variable
	portId := getPortfolioId(db, portName, username)
	if portId == 0 {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get portfolio id for name: %s", portName), http.StatusBadRequest)
		return
	}
	// Delete the security if it could be found and update database entry
	var s Security
	db.Model(&Security{}).Where("ticker=?", ticker).Where("portfolio_id=?", portId).Delete(&s)

	// Write to the response writer
	data.ToJSON(ResponseMessage{Msg: fmt.Sprintf("Deleted security of ticker %s", ticker)}, w)
}
