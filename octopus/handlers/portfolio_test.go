package handlers_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fercevik729/STLKER/eagle/protos"
	"github.com/fercevik729/STLKER/octopus/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: Update tests to utilize cookies
func TestCreatePortfolio(t *testing.T) {
	jsonStr := []byte(`{"Name": "CollegeFund","Securities":[{"Ticker": "T","Bought Price":12.50,"Shares":50},{"Ticker":"TSLA","Bought Price":120.21,"Shares":25},{"Ticker": "AMC","Bought Price":5.07,"Shares":1000}]}}`)
	req, err := http.NewRequest("POST", "/portfolio", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error("couldn't create post request to create a new portfolio:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create http recorder
	rr := httptest.NewRecorder()
	// Dial gRPC server
	conn, err := grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("[ERROR] dialing gRPC server")
		panic(err)
	}
	defer conn.Close()
	// Create a handler to listen for incoming requests
	control := handlers.NewControlHandler(log.Default(), protos.NewWatcherClient(conn))
	handler := http.HandlerFunc(control.CreatePortfolio)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want 201",
			status)
	}

	// Check response message
	expected := `{"Message":"Created portfolio named CollegeFund"}`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected %v got %v", expected, rr.Body.String())
	}

}

func TestGetPortfolio(t *testing.T) {
	expectedStr := `{"Original Value":`
	req, err := http.NewRequest("GET", "/portfolio/CollegeFund", nil)
	if err != nil {
		t.Error("couldn't create get request for CollegeFund")
	}
	// Set mux URL variables
	vars := map[string]string{
		"name": "CollegeFund",
	}
	req = mux.SetURLVars(req, vars)

	// Create http recorder
	rr := httptest.NewRecorder()
	// Dial gRPC server
	conn, err := grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("[ERROR] dialing gRPC server")
		panic(err)
	}
	defer conn.Close()
	// Create a handler to listen for incoming requests
	control := handlers.NewControlHandler(log.Default(), protos.NewWatcherClient(conn))
	handler := http.HandlerFunc(control.GetPortfolio)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want 200",
			status)
	}

	// Check response message
	if !strings.Contains(rr.Body.String(), string(expectedStr)) {
		t.Errorf("expected %v got %v", expectedStr, rr.Body.String())
	}
}

func TestUpdatePortfolio(t *testing.T) {
	// Create request and json body
	jsonStr := []byte(`{"Name": "CollegeFund","Securities":[{"Ticker": "V","Bought Price":12.50,"Shares":50},{"Ticker":"GME","Bought Price":120.21,"Shares":25},{"Ticker": "ZM","Bought Price":5.07,"Shares":1000}]}}`)
	req, err := http.NewRequest("PUT", "/portfolio/CollegeFund", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error("couldn't create put request for portfolio College Fund")
	}
	req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"name": "CollegeFund",
	}
	req = mux.SetURLVars(req, vars)
	// Create http recorder
	rr := httptest.NewRecorder()
	// Dial gRPC server
	conn, err := grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("[ERROR] dialing gRPC server")
		panic(err)
	}
	defer conn.Close()
	// Create a handler to listen for incoming requests
	control := handlers.NewControlHandler(log.Default(), protos.NewWatcherClient(conn))
	handler := http.HandlerFunc(control.UpdatePortfolio)
	handler.ServeHTTP(rr, req)
	// Check status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want 200",
			status)
	}

	expStr := `{"Message":"Updated portfolio with name CollegeFund"}`
	// Check response body
	if !strings.Contains(rr.Body.String(), expStr) {
		t.Errorf("expected %s got %s", expStr, rr.Body.String())
	}

}

func TestDeletePortfolio(t *testing.T) {
	expectedStr := `{"Message":"Deleted portfolio CollegeFund"}`
	req, err := http.NewRequest("DELETE", "/portfolio/CollegeFund", nil)
	if err != nil {
		t.Error("couldn't create delete request:", err)
	}
	// Set mux URL variables
	vars := map[string]string{
		"name": "CollegeFund",
	}
	req = mux.SetURLVars(req, vars)

	// Create http recorder
	rr := httptest.NewRecorder()
	// Dial gRPC server
	conn, err := grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("[ERROR] dialing gRPC server")
		panic(err)
	}
	defer conn.Close()
	// Create a handler to listen for incoming requests
	control := handlers.NewControlHandler(log.Default(), protos.NewWatcherClient(conn))
	handler := http.HandlerFunc(control.DeletePortfolio)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want 200",
			status)
	}

	// Check response message
	if !strings.Contains(rr.Body.String(), string(expectedStr)) {
		t.Errorf("expected %v got %v", expectedStr, rr.Body.String())
	}
}

// TODO: add more tests for update portfolio and crud operations for securities
