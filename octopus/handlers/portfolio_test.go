package handlers_test

import (
	"testing"

	"github.com/fercevik729/STLKER/octopus/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSavePortfolio(t *testing.T) {
	// Open sqlite db connection
	db, err := gorm.Open(sqlite.Open("portfolios.db"), &gorm.Config{})
	if err != nil {
		t.Error("[ERROR] Couldn't connect to database")
	}

	db.AutoMigrate(&handlers.Portfolio{}, &handlers.Security{})

	// Get generic sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		t.Error("[ERROR] Couldn't create sqlDB instance:", err)
	}
	defer sqlDB.Close()

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

	db.Where("name = ?", "CollegeFund").Preload("Securities").Find(&portfolio)

	if len(portfolio.Securities) < 2 {
		t.Errorf("Expected 2 securities got %d\n", len(portfolio.Securities))
	}

}
