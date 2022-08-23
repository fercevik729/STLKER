package handlers_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/fercevik729/STLKER/control/handlers"
	"github.com/fercevik729/STLKER/grpc/data"
	"github.com/fercevik729/STLKER/grpc/protos"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Gets a mock stock for a test case
func getMockStock(ticker, currency string) (*data.Stock, int, error) {
	req, err := http.NewRequest("GET", "/stocks", nil)
	if err != nil {
		return nil, 500, fmt.Errorf("couldn't create post request to create a new portfolio: %s", err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"ticker":   ticker,
		"currency": currency,
	})
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()
	// Dial gRPC server
	conn, err := grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("[ERROR] dialing gRPC server")
		return nil, 500, fmt.Errorf("couldn't dial gRPC server: %s", err)
	}
	defer conn.Close()
	// Create a handler to listen for incoming requests
	control := handlers.NewControlHandler(log.Default(), protos.NewWatcherClient(conn), nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.GetInfo)
	handler.ServeHTTP(rr, req)

	// Get stock
	var stock data.Stock
	data.FromJSON(&stock, rr.Body)
	return &stock, rr.Result().StatusCode, nil

}

// Logins in a mock user and returns the request with the token as a http cookie
func loginMockUser(r *http.Request) (*http.Request, error) {

	// Get token from /login endroute
	// Create request
	mockBody := handlers.User{
		Username: mockUser,
		Password: mockPass,
	}
	var buf bytes.Buffer
	err := data.ToJSON(mockBody, &buf)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login", &buf)
	if err != nil {
		return nil, fmt.Errorf("couldn't create post request to create a new portfolio: %s", err)
	}
	req.Header.Set("Content-Type", contentType)

	rr := httptest.NewRecorder()
	// Create a handler to listen for incoming requests
	control := handlers.NewControlHandler(log.Default(), nil, nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.LogIn)
	handler.ServeHTTP(rr, req)

	// Set the cookie and username
	var res map[string]string
	err = data.FromJSON(&res, rr.Body)
	cookie := &http.Cookie{
		Name:  "token",
		Value: res["Access-Token"],
	}
	var u handlers.Username
	r = r.WithContext(context.WithValue(r.Context(), u, mockUser))
	r.AddCookie(cookie)
	return r, err

}
