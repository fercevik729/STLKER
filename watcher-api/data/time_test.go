package data_test

import (
	"testing"
	"time"

	"github.com/fercevik729/STLKER/watcher-api/data"
)

func TestMarketsClosed(t *testing.T) {
	format := "2006-01-02 15:04:05 MST"

	times := map[string]bool{
		"2022-04-01 08:04:55 EDT": true,
		"2022-04-01 13:35:20 EDT": false,
		"2022-04-01 15:38:21 MST": false,
	}

	for st, expected := range times {
		parsedTime, err := time.Parse(format, st)
		if err != nil {
			t.Error(err)
		}
		actual := data.MarketsClosed(parsedTime)
		if actual != expected {
			t.Errorf("expected %v, got %v for time %s", expected, actual, st)
		}
	}

	// TODO: Update this as needed
	currTime := time.Now()
	if !data.MarketsClosed(currTime) {
		t.Errorf("expected %v, got %v", true, false)
	}
}
