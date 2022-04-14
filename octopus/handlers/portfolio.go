package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fercevik729/STLKER/octopus/data"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type STLKERModel struct {
	ID        uint         `json:"-" gorm:"primaryKey"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `json:"-" gorm:"index"`
}

// A Portfolio is a GORM model that is intended to mirror the structure
// of a simple portfolio
type Portfolio struct {
	STLKERModel
	// Name is the name of the portfolio
	Name string `json:"Name"`
	// Stocks is a slice of Security structs
	Securities []*Security `json:"Securities" gorm:"foreignKey:PortfolioID"`
}

func (p *Portfolio) calcProfits() (*Profits, error) {
	original := 0.
	new := 0.

	for _, sec := range p.Securities {

		// Update the individual security's gains and percent changes
		gain, err := strconv.ParseFloat(fmt.Sprintf("%.2f", sec.CurrPrice-sec.BoughtPrice), 64)
		if err != nil {
			return nil, err
		}
		sec.setMoves(gain, fmt.Sprintf("%.2f%%", (gain)/sec.BoughtPrice*100))

		original += sec.BoughtPrice * sec.Shares
		new += sec.CurrPrice * sec.Shares
	}

	// Round original and new
	var err error
	original, err = strconv.ParseFloat(fmt.Sprintf("%.2f", original), 64)
	if err != nil {
		return nil, err
	}

	new, err = strconv.ParseFloat(fmt.Sprintf("%.2f", new), 64)
	if err != nil {
		return nil, err
	}

	// Compute and round change and profits
	percChange := fmt.Sprintf("%.2f%%", (new-original)/original*100)

	netGain, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (new-original)), 64)
	if err != nil {
		return nil, err
	}
	return &Profits{
		OriginalValue: original,
		NewValue:      new,
		Moves:         p.Securities,
		NetGain:       netGain,
		NetChange:     percChange,
	}, nil
}

type Security struct {
	STLKERModel
	SecurityID  int     `gorm:"primary_key" json:"-"`
	Ticker      string  `json:"Ticker"`
	BoughtPrice float64 `json:"Bought Price"`
	CurrPrice   float64 `json:"Current Price"`
	Shares      float64 `json:"Shares"`
	Gain        float64 `json:"Gain"`
	Change      string  `json:"Percent Change"`
	// Currency is the destination currency of the stock
	Currency string `json:"Currency" gorm:"default:USD"`
	// Foreign key
	PortfolioID uint `json:"-"`
}

func (s *Security) setMoves(gain float64, change string) {
	s.Gain = gain
	s.Change = change
}

func (c *ControlHandler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	c.l.Println("[INFO] Handle Create Portfolio")
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

	var port Portfolio
	err := c.updateDBPort(&port)
	if err != nil {
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
	defer sqlDB.Close()

	// Check if a portfolio with that name already exists
	db.Where("name = ?", name).Preload("Securities").Find(&port)
	// Check if portfolio is empty
	if reflect.DeepEqual(port, Portfolio{}) {
		c.l.Println("[DEBUG] No results found")
		return
	}
	data.ToJSON(port, w)

}

func (c *ControlHandler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	c.l.Println("[INFO] Handle Delete Portfolio for:", name)

	db, err := gorm.Open(sqlite.Open("portfolios.db"), &gorm.Config{})
	if err != nil {
		c.l.Println("[ERROR] Couldn't connect to database")
		w.WriteHeader(http.StatusInternalServerError)
	}

	sqlDB, err := db.DB()
	if err != nil {
		c.l.Println("[ERROR] Couldn't create sqlDB instance:", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer sqlDB.Close()
	var port Portfolio

	// Delete portfolio
	db.Model(port).Where("name = ?", name).Delete(&port)
}

func (c *ControlHandler) GetProfits(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	c.l.Println("[INFO] Handle Get Profits for:", name)

	db, err := gorm.Open(sqlite.Open("portfolios.db"), &gorm.Config{})
	if err != nil {
		c.l.Println("[ERROR] Couldn't connect to database")
		w.WriteHeader(http.StatusInternalServerError)
	}

	sqlDB, err := db.DB()
	if err != nil {
		c.l.Println("[ERROR] Couldn't create sqlDB instance:", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer sqlDB.Close()
	var port *Portfolio

	// Calculate profits
	db.Where("name = ?", name).Preload("Securities").Find(&port)
	profits, err := port.calcProfits()
	if err != nil {
		c.l.Println("[ERROR] Rounding profit:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	data.ToJSON(profits, w)

}

func (c *ControlHandler) updateDBPort(port *Portfolio) error {
	// Update prices using gRPC API
	c.updatePortfolio(port)

	// Update database entry using GORM
	// Open sqlite db connection
	db, err := gorm.Open(sqlite.Open("portfolios.db"), &gorm.Config{})
	if err != nil {
		c.l.Println("[ERROR] Couldn't connect to database")
		return err
	}

	// Get generic sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		c.l.Println("[ERROR] Couldn't create sqlDB instance:", err)
		return err
	}
	defer sqlDB.Close()
	db.Model(&Portfolio{}).Where("name = ?", port.Name).Update("Securities", port.Securities)

	return nil
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

type Profits struct {
	OriginalValue float64     `json:"Original Value"`
	NewValue      float64     `json:"New Value"`
	NetGain       float64     `json:"Net Gain"`
	Moves         []*Security `json:"Securities"`
	NetChange     string      `json:"Net Change"`
}
