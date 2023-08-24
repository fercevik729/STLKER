package handlers_test

import (
	"bytes"
	"log"
	"log/slog"
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

// Create the security
func TestCreateSecurity(t *testing.T) {
	// Create the default portfolio
	TestCreatePortfolio(t)
	jsonStr := []byte(`{"Ticker": "PG", "Shares": 12}`)
	expectedStr := `{"Message":"Created PG security with 12.00 shares for portfolio CollegeFund"}`
	req, err := http.NewRequest("POST", "/portfolio", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error("couldn't create POST request for TestCreateSecurity")
	}
	req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"name": "CollegeFund",
	}
	req = mux.SetURLVars(req, vars)
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
	control := handlers.NewControlHandler(slog.Default(), protos.NewWatcherClient(conn), nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.CreateSecurity)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want 201",
			status)
	}
	if !strings.Contains(rr.Body.String(), expectedStr) {
		t.Errorf("expected %v got %v", expectedStr, rr.Body.String())
	}

}

func TestReadSecurity(t *testing.T) {
	mustContain := `"Ticker":"PG"`
	req, err := http.NewRequest("GET", "/portfolio", nil)
	if err != nil {
		t.Error("couldn't create GET request for TestReadSecurity")
	}
	req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"name":   "CollegeFund",
		"ticker": "PG",
	}
	req = mux.SetURLVars(req, vars)
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
	control := handlers.NewControlHandler(slog.Default(), protos.NewWatcherClient(conn), nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.ReadSecurity)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want 200",
			status)
	}
	if !strings.Contains(rr.Body.String(), mustContain) {
		t.Errorf("expected %v got %v", mustContain, rr.Body.String())
	}
}

func TestUpdateSecurity1(t *testing.T) {
	// Create the default portfolio
	jsonStr := []byte(`{"Ticker": "PG", "Shares": 19}`)
	expectedStr := `{"Message":"Updated security with ticker PG"}`
	req, err := http.NewRequest("PUT", "/portfolio", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error("couldn't create PUT request for TestUpdateSecurity1")
	}
	req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"name": "CollegeFund",
	}
	req = mux.SetURLVars(req, vars)
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
	control := handlers.NewControlHandler(slog.Default(), protos.NewWatcherClient(conn), nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.UpdateSecurity)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want 200",
			status)
	}
	if !strings.Contains(rr.Body.String(), expectedStr) {
		t.Errorf("expected %v got %v", expectedStr, rr.Body.String())
	}

}

func TestReadSecurity2(t *testing.T) {
	mustContain := `"Shares":19`
	req, err := http.NewRequest("GET", "/portfolio", nil)
	if err != nil {
		t.Error("couldn't create GET request for TestReadSecurity")
	}
	req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"name":   "CollegeFund",
		"ticker": "PG",
	}
	req = mux.SetURLVars(req, vars)
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
	control := handlers.NewControlHandler(slog.Default(), protos.NewWatcherClient(conn), nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.ReadSecurity)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want 200",
			status)
	}
	if !strings.Contains(rr.Body.String(), mustContain) {
		t.Errorf("expected %v got %v", mustContain, rr.Body.String())
	}
}

func TestUpdateSecurity2(t *testing.T) {
	// Create the default portfolio
	jsonStr := []byte(`{"Ticker": "V", "Shares": 124.2}`)
	expectedStr := `{"Message":"Created V security with 124.20 shares for portfolio CollegeFund"}`
	req, err := http.NewRequest("PUT", "/portfolio", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error("couldn't create PUT request for TestUpdateSecurity1")
	}
	req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"name": "CollegeFund",
	}
	req = mux.SetURLVars(req, vars)
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
	control := handlers.NewControlHandler(slog.Default(), protos.NewWatcherClient(conn), nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.UpdateSecurity)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want 201",
			status)
	}
	if !strings.Contains(rr.Body.String(), expectedStr) {
		t.Errorf("expected %v got %v", expectedStr, rr.Body.String())
	}
}

func TestReadSecurity3(t *testing.T) {
	mustContain := `"Shares":124.2`
	req, err := http.NewRequest("GET", "/portfolio", nil)
	if err != nil {
		t.Error("couldn't create GET request for TestReadSecurity")
	}
	req.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"name":   "CollegeFund",
		"ticker": "V",
	}
	req = mux.SetURLVars(req, vars)
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
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)
	// Create a handler to listen for incoming requests
	control := handlers.NewControlHandler(slog.Default(), protos.NewWatcherClient(conn), nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.ReadSecurity)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want 200",
			status)
	}
	if !strings.Contains(rr.Body.String(), mustContain) {
		t.Errorf("expected %v got %v", mustContain, rr.Body.String())
	}
}

func TestCleanup(t *testing.T) {
	TestDeletePortfolio(t)
}
