package handlers_test

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fercevik729/STLKER/control/data"
	"github.com/fercevik729/STLKER/control/handlers"
)

func TestSignUp(t *testing.T) {
	jsonStr := []byte(fmt.Sprintf(`{"Username": "%s", "Password": "%s"}`, mockUser, mockPass))
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", contentType)
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
	tokenFields, err := login()
	if err != nil {
		t.Error(err)
	}
	if len(tokenFields["Access-Token"]) == 0 {
		t.Error("no access-token field in response body")
	}
	if len(tokenFields["Refresh-Token"]) == 0 {
		t.Error("no refresh-token field in response body")
	}
}

func TestDeleteUser(t *testing.T) {
	// First log in as the user to be deleted
	tokens, err := login()
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("DELETE", "/deleteuser", nil)
	if err != nil {
		t.Error("Couldn't create DELETE request to /deleteuser:", err)
	}
	req.Header.Add("Authorization", tokens["Access-Token"])
	// Create http recorder
	rr := httptest.NewRecorder()
	control := handlers.NewControlHandler(log.Default(), nil, nil, "../database/stlker.db")
	delete := http.HandlerFunc(control.DeleteUser)
	handler := handlers.Authenticate(delete)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Error(fmt.Sprintf("handler returned wrong status code: got %v want 200",
			status))
	}

}

// login is a helper function to login a mock user
// it returns a map[string]string containing access and refresh tokens as well as the user name
// it also returns an error if at some point it was unsuccessful
func login() (map[string]string, error) {
	jsonStr := []byte(fmt.Sprintf(`{"Username": "%s", "Password": "%s"}`, mockUser, mockPass))
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", contentType)
	if err != nil {
		return nil, errors.New(fmt.Sprint("couldn't create post request to /login", err.Error()))
	}
	// Create http recorder
	rr := httptest.NewRecorder()
	control := handlers.NewControlHandler(log.Default(), nil, nil, "../database/stlker.db")
	handler := http.HandlerFunc(control.LogIn)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		return nil, fmt.Errorf("handler returned wrong status code: got %v want 200",
			status)
	}
	tokenFields := make(map[string]string)
	data.FromJSON(&tokenFields, rr.Body)
	return tokenFields, nil
}
