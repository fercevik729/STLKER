package handlers

import (
	"fmt"
	m "github.com/fercevik729/STLKER/control/models"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/fercevik729/STLKER/control/data"
)

// TODO: improve error handling to be more granular

type Username struct{}

type IsAdmin struct{}

type NamePair struct {
	Name     string
	Username string
}

type ResponseMessage struct {
	Msg string `json:"Message"`
}

// swagger:parameters createPortfolio updatePortfolio
type ReqPortfolio struct {
	// A single portfolio
	// in: body
	Body m.Portfolio
}

// swagger:route POST /portfolios portfolios createPortfolio
// Creates a new portfolio for a user
// responses:
//
//	200: messageResponse
//	400: errorResponse
//	500: errorResponse
func (c *ControlHandler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	// Retrieve the portfolio from the request body
	reqPort := m.Portfolio{}
	data.FromJSON(&reqPort, r.Body)

	// Validate the portfolio
	ok, msg := validatePortfolio(&reqPort)
	if !ok {
		c.logHTTPError(w, msg, http.StatusBadRequest)
		return
	}
	// Set username of the requested portfolio
	username := retrieveUsername(r)
	reqPort.Username = username

	// Retrieve prices for all stocks in the portfolio from the gRPC service
	c.l.Info("Retrieving updated stock prices")
	c.updatePrices(&reqPort)

	// Use the repo to create the new portfolio
	err := c.portRepo.CreateNewPortfolio(reqPort.Name, username, reqPort)
	if err != nil {
		c.logHTTPError(w, err.Error(), http.StatusBadRequest)
	}

	// Write to response body
	w.WriteHeader(http.StatusCreated)
	data.ToJSON(&ResponseMessage{
		Msg: msg,
	}, w)
}

// swagger:route GET /portfolios portfolios getPortfolios
// Outputs all of the portfolios for a user
// responses:
//
//	200: profitsResponse
//	400: errorResponse
//	500: errorResponse
func (c *ControlHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	username := retrieveUsername(r)
	isAdmin := retrieveAdmin(r)
	// If the username is admin, retrieve all portfolio names and associated usernames
	// Then return in an ordered format
	if isAdmin {
		// Return the table in json format
		table := c.portRepo.GetAllPortfoliosAdmin()
		data.ToJSON(table, w)
		return
	}
	// Otherwise retrieve all the user's portfolios
	ports := c.portRepo.GetAllPortfolios(username)

	// Calculate their total profits for each portfolio
	profits := make([]*m.Profits, 0)
	for i := range ports {
		prof, err := ports[i].CalcProfits()
		if err != nil {
			c.logHTTPError(
				w,
				fmt.Sprintf("Couldn't calculate profits for %s", ports[i].Name),
				http.StatusInternalServerError,
			)
			return
		}
		profits = append(profits, prof)
	}

	data.ToJSON(profits, w)
}

// swagger:route GET /portfolios/{name} portfolios getPortfolio
// Outputs a particular portfolio for a user
// responses:
//
//	200: profitResponse
//	400: errorResponse
//	500: errorResponse
func (c *ControlHandler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	// Retrieve portfolio name parameter and username
	name := mux.Vars(r)["name"]
	username := retrieveUsername(r)

	// Get the portfolio from the repository
	port, err := c.portRepo.GetPortfolio(name, username)
	if err != nil {
		c.logHTTPError(
			w,
			fmt.Sprintf("no portfolios found with name %s, and user %s", name, username),
			http.StatusNotFound,
		)
	}
	// Update the database with the new prices from the gRPC service
	c.updatePrices(&port)
	err = c.portRepo.UpdatePortfolio(name, username, port)
	if err != nil {
		c.logHTTPError(
			w,
			fmt.Sprintf("couldn't update portfolio with name %s, and user %s", name, username),
			http.StatusNotFound,
		)
	}
	// Calculate the profits
	profits, err := port.CalcProfits()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data.ToJSON(profits, w)
}

// swagger:route PUT /portfolios portfolios updatePortfolio
// Updates a given portfolio
// responses:
//
//	200: messageResponse
//	400: errorResponse
//	500: errorResponse
func (c *ControlHandler) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	// Get variables
	username := retrieveUsername(r)

	// Retrieve the portfolio from the request body
	reqPort := m.Portfolio{}
	data.FromJSON(&reqPort, r.Body)
	reqPort.Username = username
	name := reqPort.Name
	err := c.portRepo.UpdatePortfolio(name, username, reqPort)
	if err != nil {
		c.logHTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send response message
	msg := fmt.Sprintf("Updated portfolio with name %s", name)
	data.ToJSON(&ResponseMessage{
		Msg: msg,
	}, w)
}

// swagger:route DELETE /portfolios/{name} portfolios deletePortfolio
// Deletes a given portfolio for a user
// responses:
//
//	200: messageResponse
//	400: errorResponse
//	500: errorResponse
func (c *ControlHandler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	username := retrieveUsername(r)
	// Delete portfolio and all child securities
	err := c.portRepo.DeletePortfolio(name, username)
	if err != nil {
		c.logHTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}
	data.ToJSON(&ResponseMessage{
		Msg: fmt.Sprintf("Deleted portfolio %s", name),
	}, w)
}
