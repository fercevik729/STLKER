package handlers_test

import (
	"fmt"
	"testing"

	"github.com/fercevik729/STLKER/octopus/handlers"
)

func TestCreatePortfolio(t *testing.T) {
	// Open sqlite db connection
	db, err := handlers.NewGormDBConn("../portfolios.db")
	if err != nil {
		t.Error("couldn't connect to database")
	}
	db.AutoMigrate(&handlers.Portfolio{}, &handlers.Security{})

	Port1 := handlers.Portfolio{
		Name: "CollegeFund",
		Securities: []*handlers.Security{
			{
				Ticker:      "SPY",
				BoughtPrice: 121.98,
				Shares:      10,
			},
			{
				Ticker:      "TSLA",
				BoughtPrice: 130.12,
				Shares:      50,
			},
		},
	}
	db.Create(&Port1)
	portfolio := &handlers.Portfolio{}

	db.Where("name=?", "CollegeFund").Preload("Securities").Find(&portfolio)

	if len(portfolio.Securities) < 2 {
		t.Errorf("Expected 2 securities got %d\n", len(portfolio.Securities))
	}

}

func TestGetPortfolio(t *testing.T) {
	// Open sqlite db connection
	db, err := handlers.NewGormDBConn("../portfolios.db")
	if err != nil {
		t.Error("couldn't connect to database")
	}
	expPort := handlers.Portfolio{
		Name: "CollegeFund",
		Securities: []*handlers.Security{
			{
				Ticker:      "SPY",
				BoughtPrice: 121.98,
				Shares:      10,
			},
			{
				Ticker:      "TSLA",
				BoughtPrice: 130.12,
				Shares:      50,
			},
		},
	}
	var dbPort handlers.Portfolio
	db.Where("name=?", "CollegeFund").Preload("Securities").Find(&dbPort)

	// Make sure the number of securities matches
	if len(dbPort.Securities) != len(expPort.Securities) {
		fmt.Printf("%#v\n", dbPort.Securities[0])
		t.Error("Did not receive expected portfolio from database")
	}
}

func TestDeletePortfolio(t *testing.T) {
	// Open db conn
	db, err := handlers.NewGormDBConn("../portfolios.db")
	if err != nil {
		t.Error("couldn't connect to database")
	}
	var (
		sec  handlers.Security
		port handlers.Portfolio
	)

	name := "CollegeFund"

	// Delete portfolio and underlying securities
	db.Model(port).Where("name=?", name).Find(&port)
	db.Model(sec).Where("portfolio_id=?", port.ID).Delete(&sec)
	db.Model(port).Delete(&port)

	// Check if deletion worked
	db.Where("name=?", name).Preload("Securities").Find(&port)
	if len(port.Securities) != 0 {
		t.Errorf("Expected no securities got %d\n", len(port.Securities))
	}
}

// TODO: add more tests for update portfolio and crud operations for securities
