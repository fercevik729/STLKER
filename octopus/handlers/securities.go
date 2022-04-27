package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fercevik729/STLKER/octopus/data"
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

func (s *Security) setMoves(gain float64, change string) {
	s.Gain = gain
	s.Change = change
}

func (c *ControlHandler) AddSecurity(w http.ResponseWriter, r *http.Request) {
	// Get URI vars
	portName := mux.Vars(r)["name"]
	username := c.retrieveUsername(r)

	// Get ticker and shares info from JSON body
	var params securityData
	data.FromJSON(&params, r.Body)

	ticker := params.Ticker
	shares := params.Shares

	// Create sql db instance
	db, err := newSqlDBConn(databasePath)
	if err != nil {
		c.logHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
	}
	defer db.Close()

	// Get portfolio id
	portId, err := getPortfolioId(db, portName, username)
	if err != nil {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get portfolio id for name: %s, and user: %s", portName, username), http.StatusBadRequest)
	}

	// Create insert sql query
	stmt, err := db.Prepare(`INSERT INTO securities(created_at, security_id, ticker, bought_price, curr_price, shares, currency, portfolio_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		c.logHTTPError(w, "Could't prepare insert query string", http.StatusInternalServerError)
		return
	}
	// Get current price of security
	stock, err := Info(ticker, "USD", c.client)
	if err != nil {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get updated price for %s", ticker), http.StatusBadRequest)
		return
	}
	// Execute query
	res, err := stmt.Exec(time.Now(), 0, ticker, stock.Price, stock.Price, shares, "USD", portId)
	if err != nil {
		c.logHTTPError(w, "Couldn't execute insert query", http.StatusInternalServerError)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		c.logHTTPError(w, "Couldn't retrieve affected rows", http.StatusInternalServerError)
		return
	}
	c.l.Printf("[DEBUG] Added %d row\n", rows)

}

func (c *ControlHandler) ReadSecurity(w http.ResponseWriter, r *http.Request) {
	// Get URI vars
	portName, ticker, username := c.getSecurityVars("Read Security", r)
	db, err := newSqlDBConn(databasePath)
	if err != nil {
		c.logHTTPError(w, "Couldn't connct to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get portfolio_id
	portId, err := getPortfolioId(db, portName, username)
	if err != nil {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get updated price for %s", portName), http.StatusBadRequest)
		return
	}
	rows, err := db.Query("SELECT ticker, bought_price, curr_price, shares, gain, change FROM securities WHERE ticker=? AND portfolio_id=?", ticker, portId)
	if err != nil {
		c.logHTTPError(w, fmt.Sprintf("Couldn't query database for %s security in %s portfolio", ticker, portName), http.StatusBadRequest)
		return
	}
	// Iterate over all the securities with that ticker
	for rows.Next() {
		var (
			ticker      string
			boughtPrice float64
			currPrice   float64
			shares      float64
			gain        float64
			change      string
		)

		rows.Scan(&ticker, &boughtPrice, &currPrice, &shares, &gain, &change)
		security := Security{
			Ticker:      ticker,
			BoughtPrice: boughtPrice,
			CurrPrice:   currPrice,
			Shares:      shares,
			Gain:        gain,
			Change:      change,
		}
		// Update the security
		c.updateSecurities(&security)

		// Write to response writer
		data.ToJSON(&security, w)
	}
}

func (c *ControlHandler) EditSecurity(w http.ResponseWriter, r *http.Request) {
	// Get URI vars
	portName := mux.Vars(r)["name"]
	username := c.retrieveUsername(r)

	// Get ticker and shares info from JSON body
	var params securityData
	data.FromJSON(&params, r.Body)

	ticker := params.Ticker
	shares := params.Shares

	// Create sql db instance
	db, err := newSqlDBConn(databasePath)
	if err != nil {
		c.logHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get portfolio id
	portId, err := getPortfolioId(db, portName, username)
	if err != nil {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get portfolio id for name: %s", portName), http.StatusBadRequest)
		return
	}
	stmt, err := db.Prepare(`UPDATE securities SET shares=?, updated_at=? WHERE ticker=? AND portfolio_id=?`)
	if err != nil {
		c.logHTTPError(w, "Couldn't prepare update query string", http.StatusInternalServerError)
		return
	}
	res, err := stmt.Exec(shares, time.Now(), ticker, portId)
	if err != nil {
		c.logHTTPError(w, "Couldn't execute update query", http.StatusInternalServerError)
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		c.logHTTPError(w, "Couldn't retrieve affected rows", http.StatusInternalServerError)
		return
	}
	// Write a response to the client to tell them if the security was there in the first place
	if affected > 0 {
		msg := fmt.Sprintf("Updated %d security", affected)
		c.l.Printf("[DEBUG] %s\n", msg)

		data.ToJSON(&ResponseMessage{
			Msg: msg,
		}, w)
	} else {
		msg := fmt.Sprintf("%s did not have security: %s", portName, ticker)
		c.l.Printf("[DEBUG] %s\n", msg)
		w.Header().Set("Content-Type", "application/json")
		data.ToJSON(&ResponseMessage{
			Msg: msg,
		}, w)
	}

}

func (c *ControlHandler) DeleteSecurity(w http.ResponseWriter, r *http.Request) {
	portName, ticker, username := c.getSecurityVars("Delete Security", r)
	// Connect to database
	db, err := newSqlDBConn(databasePath)
	if err != nil {
		c.logHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	// Get portfolio with the name specified by the mux variable
	portId, err := getPortfolioId(db, portName, username)
	if err != nil {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get portfolio id for name: %s", portName), http.StatusBadRequest)
		return
	}
	// Delete the security if it could be found and update database entry
	stmt, err := db.Prepare(`DELETE FROM securities WHERE ticker=? AND portfolio_id=?`)
	if err != nil {
		c.logHTTPError(w, "Couldn't prepare delete query string", http.StatusInternalServerError)
		return
	}
	// Execute the query
	res, err := stmt.Exec(ticker, portId)
	if err != nil {
		c.logHTTPError(w, "Couldn't execute delete query", http.StatusInternalServerError)
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		c.logHTTPError(w, "Couldn't retrieve affected rows", http.StatusInternalServerError)
		return
	}
	// Write a response to the client to tell them if the security was there in the first place
	if affected > 0 {
		msg := fmt.Sprintf("Deleted %d security", affected)
		c.l.Printf("[DEBUG] %s\n", msg)

		data.ToJSON(&ResponseMessage{
			Msg: msg,
		}, w)
	} else {
		msg := fmt.Sprintf("%s did not have security: %s", portName, ticker)
		c.l.Printf("[DEBUG] %s\n", msg)
		w.Header().Set("Content-Type", "application/json")
		data.ToJSON(&ResponseMessage{
			Msg: msg,
		}, w)
	}
}
