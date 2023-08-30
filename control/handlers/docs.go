// Package classification of STLKER API
//
// # Documentation for STLKER API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import m "github.com/fercevik729/STLKER/control/models"

// A profits is a list of profits for each of the portfolios returned to the client
// swagger:response profitsResponse
type profitsResponseWrapper struct {
	// Profits for all portfolios for a user
	// in: body
	Body []m.Profits
}

// A profitResponse is a single portfolio returned to the client
// swagger:response profitResponse
type profitResponseWrapper struct {
	// A single portfolio's profits
	// in: body
	Body m.Profits
}

// A securityResponse is a single security's information returned to the client
// swagger:response securityResponse
type securityResponseWrapper struct {
	// A single security
	// in: body
	Body m.Security
}

// A stockResponse is information about a single stock returned to the client
// swagger:response stockResponse
type stockResponseWrapper struct {
	// A single stock
	// in: body
	Body Stock
}

// A moreStockResponse is more information about a single stock
// swagger:response moreStockResponse
type moreStockResponseWrapper struct {
	// A single stock
	// in: body
	Body MoreStock
}

// A messageResponse returns diagnostic information to the client
// swagger:response messageResponse
type messageResponseWrapper struct {
	// A message
	// in: body
	Body ResponseMessage
}

// swagger:parameters getPortfolio createSecurity updateSecurity deletePortfolio
type portfolioNameParamWrapper struct {
	// Name of the portfolio
	// in: path
	// required: true
	Name string `json:"name"`
}

// swagger:parameters deleteSecurity readSecurity
type portfolioNameTickerParamWrapper struct {
	// Name of the portfolio
	// in: path
	Name string `json:"name"`
	// Ticker of the security
	// in: path
	Ticker string `json:"ticker"`
}

// swagger:parameters moreInfo
type tickerParamWrapper struct {
	// Ticker of the security
	// in: path
	Ticker string `json:"ticker"`
}

// swagger:parameters getInfo
type tickerCurrencyParamWrapper struct {
	// Ticker of the security
	// in: path
	Ticker string `json:"ticker"`
	// Destination currency for the security's unit prices
	// in: path
	Currency string `json:"currency"`
}

// An errorResponse is an empty data structure to represent an http error
// swagger:response errorResponse
type errorResponseWrapper struct{}
