package handlers

import (
	"fmt"
	m "github.com/fercevik729/STLKER/control/models"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	"github.com/fercevik729/STLKER/control/data"
	"github.com/gorilla/mux"
)

// swagger:parameters createSecurity updateSecurity
type ReqSecurityWrapper struct {
	// A single portfolio
	// in: body
	Body ReqSecurity
}

// swagger:model
type ReqSecurity struct {
	// ticker of the security
	//
	// required: true
	// example: AAPL
	Ticker string `json:"Ticker"`
	// number of shares of the security
	//
	// required: true
	// example: 5.32
	Shares float64 `json:"Shares"`
}

// convertSecurity converts a ReqSecurity struct to a Security struct
func (c *ControlHandler) convertSecurity(params ReqSecurity, portName string, username string) (m.Security, error) {
	ticker := params.Ticker
	shares := params.Shares

	// Get portfolio id
	portId := c.portRepo.GetPortfolioId(portName, username)
	if portId == 0 {
		return m.Security{}, errors.Errorf("Couldn't get portfolio %s", portName)
	}

	// Get stock info from gRPC service
	stock, err := Info(ticker, "USD", c.client)
	if err != nil {
		return m.Security{}, errors.Errorf("Couldn't get updated price for %s", ticker)
	}
	// Get the price
	price, _ := strconv.ParseFloat(stock.Price, 64)

	// Create the security struct
	newSecurity := m.Security{
		Ticker:      ticker,
		BoughtPrice: price,
		CurrPrice:   price,
		Shares:      shares,
		Currency:    "USD",
		PortfolioID: portId,
	}

	return newSecurity, nil
}

// swagger:route POST /portfolios/{name} securities createSecurity
// Creates a new security
// responses:
//
//	200: messageResponse
//	400: errorResponse
//	500: errorResponse
func (c *ControlHandler) CreateSecurity(w http.ResponseWriter, r *http.Request) {
	// Get URI vars
	portName := mux.Vars(r)["name"]
	username := retrieveUsername(r)

	// Get ticker and shares info from JSON body
	var params ReqSecurity
	data.FromJSON(&params, r.Body)
	// Check if the payload is empty
	if params == (ReqSecurity{}) {
		c.logHTTPError(w, "Bad request payload", http.StatusBadRequest)
		return
	}

	// Convert from DTO -> DAO
	security, err := c.convertSecurity(params, portName, username)
	if err != nil {
		c.logHTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	portId, ticker := security.PortfolioID, security.Ticker
	// Security already exists
	if exists := c.secRepo.Exists(portId, ticker); exists {
		w.WriteHeader(http.StatusBadRequest)
	}

	c.secRepo.CreateSecurity(security)
	w.WriteHeader(http.StatusCreated)
}

// swagger:route GET /portfolios/{name}/{ticker} securities readSecurity
// Outputs a security's details to the client
// responses:
//
//	200: securityResponse
//	400: errorResponse
//	500: errorResponse
func (c *ControlHandler) ReadSecurity(w http.ResponseWriter, r *http.Request) {
	// Get URI vars
	portName, ticker, username := retrieveSecurityVars(r)
	// Get portfolio_id
	portId := c.portRepo.GetPortfolioId(portName, username)
	if portId == 0 {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get updated price for %s", portName), http.StatusBadRequest)
		return
	}
	// Get the security from the database
	security := c.secRepo.GetSecurity(portId, ticker)

	// Update the security
	c.updateSecurities(&security)

	// Write it to the response
	data.ToJSON(&security, w)
}

// swagger:route PUT /portfolios/{name} securities updateSecurity
// Updates a security's information for a given portfolio
// responses:
//
//	200: messageResponse
//	400: errorResponse
//	500: errorResponse
func (c *ControlHandler) UpdateSecurity(w http.ResponseWriter, r *http.Request) {
	// Get request vars
	portName := mux.Vars(r)["name"]
	username := retrieveUsername(r)
	var params ReqSecurity
	data.FromJSON(&params, r.Body)

	// Convert from DTO to DAO
	security, err := c.convertSecurity(params, portName, username)
	if err != nil {
		c.logHTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Security doesn't exist so create it
	portId, ticker := security.PortfolioID, security.Ticker
	if exists := c.secRepo.Exists(portId, ticker); !exists {
		c.secRepo.CreateSecurity(security)
		w.WriteHeader(http.StatusCreated)
		return
	}

	// Update the portfolio
	c.secRepo.UpdateSecurity(security)
}

// swagger:route DELETE /portfolios/{name}/{ticker} securities deleteSecurity
// Deletes a security from a given portfolio
// responses:
//
//	200: messageResponse
//	400: errorResponse
//	500: errorResponse
func (c *ControlHandler) DeleteSecurity(w http.ResponseWriter, r *http.Request) {
	portName, ticker, username := retrieveSecurityVars(r)
	// Get portfolio with the name specified by the mux variable
	portId := c.portRepo.GetPortfolioId(portName, username)
	if portId == 0 {
		c.logHTTPError(w, fmt.Sprintf("Couldn't get portfolio id for name: %s", portName), http.StatusBadRequest)
		return
	}
	// Delete the security if it could be found and update database entry
	c.secRepo.DeleteSecurity(ticker, portId)

	// Write to the response writer
	data.ToJSON(ResponseMessage{Msg: fmt.Sprintf("Deleted security of ticker %s", ticker)}, w)
}
