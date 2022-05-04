package handlers_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fercevik729/STLKER/octopus/handlers"
)

func TestSignUp(t *testing.T) {
	jsonStr := []byte(`{"Username": "jordanbelfort", "Password": "ilovemoney"}`)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Error("couldn't create post request to /signup:", err)
	}
	// Create http recorder
	rr := httptest.NewRecorder()
	control := handlers.NewControlHandler(log.Default(), nil, nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.SignUp)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want 201",
			status)
	}
}
func TestLogIn(t *testing.T) {
	jsonStr := []byte(`{"Username": "wolfofwallstreet", "Password": "ilovemoney"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Error("couldn't create post request to /login:", err)
	}
	// Create http recorder
	rr := httptest.NewRecorder()
	control := handlers.NewControlHandler(log.Default(), nil, nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.LogIn)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want 200",
			status)
	}

	_, err = rr.Result().Request.Cookie("token")
	if err != nil {
		t.Error("couldn't retrieve 'token' cookie:", err)
	}
	_, err = rr.Result().Request.Cookie("refreshToken")
	if err != nil {
		t.Error("couldn't retrieve 'refreshToken' cookie:", err)
	}
}
