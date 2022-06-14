package handlers_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fercevik729/STLKER/control/handlers"
	"github.com/fercevik729/STLKER/grpc/data"
	"github.com/fercevik729/STLKER/grpc/protos"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Info Tests
// Should pass
func TestInfo1(t *testing.T) {
	stock, status, err := getMockStock("NFLX", "USD")
	if err != nil {
		t.Error(err)
	}
	if status != http.StatusOK {
		t.Errorf("Wanted status %d, got %d", http.StatusOK, status)
	}
	if !(stock.Symbol == "NFLX" && len(stock.Open) > 0 && len(stock.High) > 0) {
		t.Errorf("bad response stock: %#v", stock)
	}
}

func TestInfo2(t *testing.T) {
	_, status, err := getMockStock("PSDA", "USD")
	if err != nil {
		t.Error(err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Wanted status %d, got %d", 500, status)
	}
}

func TestInfo3(t *testing.T) {
	validTickers := []string{"SPY", "AAPL", "TSLA", "AMC", "GME"}
	invalidTickers := []string{"", "2231", "AAAAPL", "PLA!2D"}

	// Test several valid tickers
	for _, tick := range validTickers {
		stock, status, err := getMockStock(tick, currency)
		if err != nil {
			t.Error(err)
		}
		if status != http.StatusOK {
			t.Errorf("wanted status %d, got %d for ticker %s", http.StatusOK, status, tick)
		}
		if !(stock.Symbol == tick && len(stock.Open) > 0 && len(stock.High) > 0) {
			t.Errorf("bad response stock: %#v", stock)
		}
		time.Sleep(time.Second)

	}
	// Test invalid tickers
	for _, ticker := range invalidTickers {
		_, status, err := getMockStock(ticker, currency)
		if err != nil {
			t.Error(err)
		}
		if status != http.StatusBadRequest {
			t.Errorf("Wanted status %d, got %d", 500, status)
		}
	}
}

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
