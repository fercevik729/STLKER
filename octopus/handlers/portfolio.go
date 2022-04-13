package handlers

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/fercevik729/STLKER/octopus/data"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// A Portfolio is a GORM model that is intended to mirror the structure
// of a simple portfolio
type Portfolio struct {
	gorm.Model
	ID uint `gorm:"primary_key"`
	// Name is the name of the portfolio
	Name string `json:"Name"`
	// Stocks is a slice of Security structs
	Securities []*Security `json:"Securities" gorm:"foreignKey:PortfolioID"`
}

type Security struct {
	gorm.Model
	SecurityID  int     `gorm:"primary_key"`
	Ticker      string  `json:"Ticker"`
	BoughtPrice float64 `json:"Bought Price"`
	CurrPrice   float64 `json:"Current Price"`
	Shares      float64 `json:"Shares"`
	// Currency is the destination currency of the stock
	Currency string `json:"Currency" gorm:"default:USD"`
	// Foreign key
	PortfolioID uint
}

func (c *ControlHandler) SavePortfolio(w http.ResponseWriter, r *http.Request) {
	c.l.Println("[INFO] Handle Save Portfolio")
	// Retrieve the portfolio from the request body
	reqPort := Portfolio{}
	data.FromJSON(&reqPort, r.Body)

	// Check if name is empty which generally signifies that the json body was misconstrued
	// or if the name contains spaces, since it isn't compatible with the URI
	if reqPort.Name == "" || strings.Contains(reqPort.Name, " ") {
		c.l.Println("[ERROR] Bad portfolio request")
		w.WriteHeader(http.StatusBadRequest)
	}

	// Open sqlite db connection
	db, err := gorm.Open(sqlite.Open("portfolios.db"), &gorm.Config{})
	if err != nil {
		c.l.Println("[ERROR] Couldn't connect to database")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Get generic sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		c.l.Println("[ERROR] Couldn't create sqlDB instance:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Migrate schema
	db.AutoMigrate(&Portfolio{}, &Security{})

	sqlPort := Portfolio{}
	// Check if a portfolio with that name already exists
	db.First(&sqlPort, "name = ?", reqPort.Name)

	// If a portfolio with that name does exist return an error
	if !reflect.DeepEqual(&sqlPort, &Portfolio{}) {
		c.l.Println("[ERROR] A portfolio with that name already exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Retrieve prices for all stocks in the portfolio
	c.l.Println("[INFO] Retrieving updated stock prices")
	c.updatePortfolio(&reqPort)

	// Create portfolio entry
	db.Create(&reqPort)
	c.l.Println("[DEBUG] Created portfolio named", reqPort.Name)

	// Close database connection
	sqlDB.Close()

}

func (c *ControlHandler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	// Retrieve portfolio name parameter
	name := mux.Vars(r)["name"]
	c.l.Println("[INFO] Handle Get Portfolio for:", name)

	// Open sqlite db connection
	db, err := gorm.Open(sqlite.Open("portfolios.db"), &gorm.Config{})
	if err != nil {
		c.l.Println("[ERROR] Couldn't connect to database")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Get generic sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		c.l.Println("[ERROR] Couldn't create sqlDB instance:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer sqlDB.Close()

	var port Portfolio
	// Check if a portfolio with that name already exists
	db.Where("name = ?", name).Preload("Securities").Find(&port)
	// Check if portfolio is empty
	if reflect.DeepEqual(port, Portfolio{}) {
		c.l.Println("[DEBUG] No results found")
		return
	}
	data.ToJSON(port, w)

}

func (c *ControlHandler) updatePortfolio(port *Portfolio) {
	// Concurrently retrieve stock prices
	wg := &sync.WaitGroup{}
	for _, sec := range port.Securities {
		wg.Add(1)
		go func(s *Security) {
			// Get security information using Info method defined in driver.go
			st, err := Info(s.Ticker, s.Currency, c.client)
			if err != nil {
				c.l.Println("[ERROR] Couldn't get info for ticker:", s.Ticker)
				return
			}
			// Parse the stock price
			price, err := strconv.ParseFloat(st.Price, 64)
			if err != nil {
				c.l.Println("[ERROR] Couldn't parse stock price for ticker:", s.Ticker)
				return
			}
			// Set stock price in target currency (USD by default)
			if s.Currency == "" {
				s.Currency = "USD"
			}
			c.l.Println("[DEBUG] Got price for ticker:", s.Ticker, "in", s.Currency)
			s.CurrPrice = price
			wg.Done()
		}(sec)
	}
	wg.Wait()

}
