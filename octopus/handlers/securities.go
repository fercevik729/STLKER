package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Security struct {
	STLKERModel
	SecurityID  int     `gorm:"primary_key"`
	Ticker      string  `json:"Ticker"`
	BoughtPrice float64 `json:"Bought Price"`
	CurrPrice   float64 `json:"Current Price"`
	Shares      float64 `json:"Shares"`
	Gain        float64 `json:"Gain"`
	Change      string  `json:"Percent Change"`
	// Currency is the destination currency of the stock
	Currency string `json:"Currency" gorm:"default:USD"`
	// Foreign key
	PortfolioID uint
}

func (s *Security) setMoves(gain float64, change string) {
	s.Gain = gain
	s.Change = change
}

func NewDBConn(databaseName string) (*sql.DB, error) {
	sqliteDatabase, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return nil, err
	}

	return sqliteDatabase, nil

}

func (c *ControlHandler) DeleteSecurity(w http.ResponseWriter, r *http.Request) {
	portName, ticker := c.getSecurityVars("Delete Security", r)

	// Connect to database
	db, err := NewDBConn("portfolios.db")
	if err != nil {
		c.l.Println("[ERROR] Couldn't connect to database:", err)
		return
	}
	defer db.Close()
	// Get portfolio with the name specified by the mux variable
	portId, err := getPortfolioId(db, portName)
	if err != nil {
		c.LogHTTPError(w, fmt.Sprintf("Couldn't get portfolio id for name: %s", portName), http.StatusBadRequest)
	}
	// Delete the security if it could be found and update database entry
	stmt, err := db.Prepare(`DELETE FROM securities WHERE ticker=? AND portfolio_id=?`)
	if err != nil {
		c.LogHTTPError(w, "Couldn't prepare delete query string", http.StatusInternalServerError)
	}
	res, err := stmt.Exec(ticker, portId)
	if err != nil {
		c.LogHTTPError(w, "Couldn't execute delete query", http.StatusInternalServerError)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		c.LogHTTPError(w, "Couldn't retrieve affected rows", http.StatusInternalServerError)
	}
	c.l.Printf("[DEBUG] Deleted %d rows\n", affected)
}

func (c *ControlHandler) AddSecurity(w http.ResponseWriter, r *http.Request) {

}

func (c *ControlHandler) EditSecurity(w http.ResponseWriter, r *http.Request) {
	//portName, ticker := c.getSecurityVars("Edit Security", r)

}

// getPortfolioId returns a portfolio's id provided its name
func getPortfolioId(db *sql.DB, portName string) (int, error) {
	// Execute query
	rows, err := db.Query("SELECT id FROM portfolios WHERE name=?", portName)
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

// getVars returns a portfolioName and ticker string vars. It also logs
// the method being handled
func (c *ControlHandler) getSecurityVars(method string, r *http.Request) (portName string, ticker string) {
	// Get URI vars
	vars := mux.Vars(r)
	portName = vars["name"]
	ticker = vars["ticker"]
	// Log endpoint
	c.l.Printf("[INFO] Handle %s for portfolio: %s, for ticker: %s", method, portName, ticker)
	return portName, ticker
}
