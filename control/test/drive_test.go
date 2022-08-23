package handlers_test

import (
	"net/http"
	"testing"
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
	s, status, err := getMockStock("", "USD")
	if err != nil {
		t.Error(err)
	}
	if status != http.StatusBadRequest {
		t.Errorf("Wanted status %d, got %d", 500, status)
	}
	if len(s.Price) != 0 {
		t.Errorf("Wanted empty price, got %s", s.Price)
	}
}
