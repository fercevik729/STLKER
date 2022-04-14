package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
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

// TODO: Use sqlite driver instead of ORM for these CRUD ops because GORM doesn't work well with nested associations
func (c *ControlHandler) DeleteSecurity(w http.ResponseWriter, r *http.Request) {
	// Get URI vars
	vars := mux.Vars(r)
	portName := vars["name"]
	ticker := vars["ticker"]
	c.l.Println("[INFO] Handle Delete Security for portfolio:", portName, "and security:", ticker)

	// Connect to database
	// Delete the security if it could be found and update database entry
	// TODO: replace this comment block with a sqlite query executer
	/*
		DELETE FROM securities
		WHERE Ticker = (?) AND PortfolioId = (?);
	*/

}

func (c *ControlHandler) AddSecurity(w http.ResponseWriter, r *http.Request) {

}

func (c *ControlHandler) EditSecurity(w http.ResponseWriter, r *http.Request) {

}
