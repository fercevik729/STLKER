// Package classification of STLKER API
//
// Documentation for STLKER API
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

import "github.com/fercevik729/STLKER/grpc/data"

// A portfoliosResponse is a list of portfolios returned to the client
// swagger:response portfoliosResponse
type portfoliosResponseWrapper struct {
	// All portfolios for a user
	// in: body
	Body []Portfolio
}

// A portfolioResponse is a single portfolio returned to the client
// swagger:response portfolioResponse
type portfolioResponseWrapper struct {
	// A single portfolio
	// in: body
	Body Portfolio
}

// A securityResponse is a single security's information returned to the client
// swagger:response securityResponse
type securityResponseWrapper struct {
	// A single security
	// in: body
	Body Security
}

// A stockResponse is information about a single stock returned to the client
type stockResponseWrapper struct {
	// A single stock
	// in: body
	Body data.Stock
}

// noContent is used to signify no content is returned to the sdk
// swagger:response noContent
type noContentWrapper struct{}

// swagger:parameters getPortfolio createSecurity updateSecurity deleteSecurity
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

type tickerCurrencyParamWrapper struct {
	// Ticker of the security
	// in: path
	Ticker string `json:"ticker"`
	// Destination currency for the security's unit prices
	Currency string `json:"currency"`
}

// An errorResponse is an empty data structure to represent an http error
// swagger:response errorResponse
type errorResponseWrapper struct{}
