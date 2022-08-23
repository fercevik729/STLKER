package handlers_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fercevik729/STLKER/control/handlers"
	"github.com/fercevik729/STLKER/grpc/protos"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCreatePortfolio(t *testing.T) {
	jsonStr := []byte(`{"Name": "CollegeFund","Securities":[{"Ticker": "T","Bought Price":12.50,"Shares":50},{"Ticker":"TSLA","Bought Price":120.21,"Shares":25},{"Ticker": "AMC","Bought Price":5.07,"Shares":1000}]}}`)
	req, err := http.NewRequest("POST", "/portfolio", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error("couldn't create POST request for TestCreatePortfolio")
	}
	req.Header.Set("Content-Type", "application/json")
	// Login and set the token and username
	req, err = loginMockUser(req)
	if err != nil {
		t.Error(err)
	}
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
	control := handlers.NewControlHandler(log.Default(), protos.NewWatcherClient(conn), nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.CreatePortfolio)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want 201",
			status)
	}

	// Check response message
	expected := `{"Message":"Created portfolio named CollegeFund for wolfofwallstreet"}`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected %v got %v", expected, rr.Body.String())
	}

}

func TestGetPortfolio(t *testing.T) {
	expectedStr := `{"Portfolio Name":"CollegeFund","Original Value":8700.25,"Current Value":33110,"Net Gain":24409.75,"Net Change":"280.56%","Securities":[{"Ticker":"T","Bought Price":12.5,"Current Price":18.13,"Shares":50,"Gain":281.5,"Percent Change":"45.04%","Currency":"USD"},{"Ticker":"TSLA","Bought Price":120.21,"Current Price":869.74,"Shares":25,"Gain":18738.25,"Percent Change":"623.52%","Currency":"USD"},{"Ticker":"AMC","Bought Price":5.07,"Current Price":10.46,"Shares":1000,"Gain":5390,"Percent Change":"106.31%","Currency":"USD"}]}`
	req, err := http.NewRequest("GET", "/portfolio", nil)
	if err != nil {
		t.Error("couldn't create GET request for TestGetPortfolio")
	}
	// Set mux URL variables
	vars := map[string]string{
		"name": "CollegeFund",
	}
	req = mux.SetURLVars(req, vars)

	// login and set the token and username
	req, err = loginMockUser(req)
	if err != nil {
		t.Error(err)
	}
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
	control := handlers.NewControlHandler(log.Default(), protos.NewWatcherClient(conn), nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.GetPortfolio)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want 200",
			status)
	}

	// Check response message
	fmt.Println(rr.Body.String())
	if !strings.Contains(rr.Body.String(), string(expectedStr)) {
		t.Errorf("expected %v got %v", expectedStr, rr.Body.String())
	}
}

func TestUpdatePortfolio(t *testing.T) {
	// Create request and json body
	jsonStr := []byte(`{"Name": "CollegeFund","Securities":[{"Ticker": "V","Bought Price":12.50,"Shares":50},{"Ticker":"GME","Bought Price":120.21,"Shares":25},{"Ticker": "ZM","Bought Price":5.07,"Shares":1000}]}}`)
	req, err := http.NewRequest("PUT", "/portfolio", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error("couldn't create PUT request for TestUpdatePortfolio")
	}
	req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"name": "CollegeFund",
	}
	req = mux.SetURLVars(req, vars)
	// login and set the token and username
	req, err = loginMockUser(req)
	if err != nil {
		t.Error(err)
	}
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
	control := handlers.NewControlHandler(log.Default(), protos.NewWatcherClient(conn), nil, "../database/stlker.db")
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
	req, err := http.NewRequest("DELETE", "/portfolio", nil)
	if err != nil {
		t.Error("couldn't create DELETE request for TestDeletePortfolio")
	}
	// Set mux URL variables
	vars := map[string]string{
		"name": "CollegeFund",
	}
	req = mux.SetURLVars(req, vars)
	// Set token and username
	req, err = loginMockUser(req)
	if err != nil {
		t.Fail()
	}

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
	control := handlers.NewControlHandler(log.Default(), protos.NewWatcherClient(conn), nil, "../database/stlker.db")
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
