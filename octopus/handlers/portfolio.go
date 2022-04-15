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

func (c *ControlHandler) LogHTTPError(w http.ResponseWriter, errorMsg string, errorCode int) {
	c.l.Printf("[ERROR] %s\n", errorMsg)
	http.Error(w, errorMsg, errorCode)
}

type ResponseMessage struct {
	Msg string `json:"Message"`
}

type STLKERModel struct {
	ID        uint         `gorm:"primaryKey" json:"-"`
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

type Profits struct {
	OriginalValue float64     `json:"Original Value"`
	NewValue      float64     `json:"Current Value"`
	NetGain       float64     `json:"Net Gain"`
	NetChange     string      `json:"Net Change"`
	Moves         []*Security `json:"Securities"`
}

func (p *Portfolio) calcProfits() (*Profits, error) {
	original := 0.
	new := 0.
	// Iterate over securities and calculate change and percent change
	for _, sec := range p.Securities {
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

	// Return profits
	return &Profits{
		OriginalValue: original,
		NewValue:      new,
		Moves:         p.Securities,
		NetGain:       netGain,
		NetChange:     percChange,
	}, nil
}

func (c *ControlHandler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	c.l.Println("[INFO] Handle Create Portfolio")
	// Retrieve the portfolio from the request body
	reqPort := Portfolio{}
	data.FromJSON(&reqPort, r.Body)

	// Check if name is empty which generally signifies that the json body was misconstrued
	// or if the name contains spaces, since it isn't compatible with the URI
	if reqPort.Name == "" || strings.Contains(reqPort.Name, " ") {
		c.LogHTTPError(w, "Bad portfolio request. Name shouldn't be empty or contain spaces", http.StatusBadRequest)
	}

	// Open sqlite db connection
	db, err := gorm.Open(sqlite.Open("portfolios.db"), &gorm.Config{})
	if err != nil {
		c.LogHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}

	// Get generic sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		c.LogHTTPError(w, "Couldn't create sql instance", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()
	// Migrate schema
	db.AutoMigrate(&Portfolio{}, &Security{})

	sqlPort := Portfolio{}
	// Check if a portfolio with that name already exists
	db.First(&sqlPort, "name = ?", reqPort.Name)

	// If a portfolio with that name does exist return an error
	if !reflect.DeepEqual(&sqlPort, &Portfolio{}) {
		c.LogHTTPError(w, "A portfolio with that name already exists", http.StatusBadRequest)
		return
	}
	// Retrieve prices for all stocks in the portfolio
	c.l.Println("[INFO] Retrieving updated stock prices")
	c.updatePrices(&reqPort)

	// Create portfolio entry
	db.Create(&reqPort)
	msg := fmt.Sprintf("Created portfolio named %s", reqPort.Name)
	c.l.Printf("[DEBUG] %s", msg)

	// Write to response body
	w.WriteHeader(http.StatusCreated)
	data.ToJSON(&ResponseMessage{
		Msg: msg,
	}, w)

}

func (c *ControlHandler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	// Retrieve portfolio name parameter
	name := mux.Vars(r)["name"]
	c.l.Println("[INFO] Handle Get Portfolio for:", name)

	// Open sqlite db connection
	db, err := gorm.Open(sqlite.Open("portfolios.db"), &gorm.Config{})
	if err != nil {
		c.LogHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}

	// Get generic sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		c.LogHTTPError(w, "Couldn't create sql instance", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	var port Portfolio
	// Check if a portfolio with that name already exists
	db.Where("name = ?", name).Preload("Securities").Find(&port)
	// Check if any results were found
	if port.ID == 0 {
		c.l.Println("[DEBUG] No results found")
		data.ToJSON(&ResponseMessage{
			Msg: fmt.Sprintf("A portfolio with name %s could not be found", name),
		}, w)
		return
	}
	// Update the database entry with the new prices
	c.updateDB(w, &port)
	// Calculate the profits
	profits, err := port.calcProfits()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data.ToJSON(profits, w)

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
		c.LogHTTPError(w, "Couldn't create sql db instance", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	// Delete portfolio and all child securities
	var (
		sec  Security
		port Portfolio
	)
	// Update the portfolio in the db
	db.Model(port).Where("name=?", name).Find(&port)
	// Check if any results were found
	if port.ID == 0 {
		c.l.Println("[DEBUG] No results found")
		data.ToJSON(&ResponseMessage{
			Msg: fmt.Sprintf("A portfolio with name %s could not be found", name),
		}, w)
		return
	}
	db.Model(sec).Where("portfolio_id=?", port.ID).Delete(&sec)
	db.Model(port).Delete(&port)

	data.ToJSON(&ResponseMessage{
		Msg: fmt.Sprintf("Deleted portfolio %s", name),
	}, w)
}

func (c *ControlHandler) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	// Retrieve the portfolio from the request body
	reqPort := Portfolio{}
	data.FromJSON(&reqPort, r.Body)

	c.l.Println("[INFO] Handle Update Portfolio for:", name)

	// Open sqlite db connection
	db, err := gorm.Open(sqlite.Open("portfolios.db"), &gorm.Config{})
	if err != nil {
		c.LogHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}

	// Get generic sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		c.LogHTTPError(w, "Couldn't create sql instance", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	var (
		sqlPort Portfolio
		sec     Security
	)
	// Find the portfolio in the db
	db.Model(sqlPort).Where("name=?", name).Find(&sqlPort)
	// Check if any results were found
	if sqlPort.ID == 0 {
		c.l.Println("[DEBUG] No results found")
		data.ToJSON(&ResponseMessage{
			Msg: fmt.Sprintf("A portfolio with name %s could not be found", name),
		}, w)
		return
	}
	//Delete previous portfolio
	db.Model(sec).Where("portfolio_id=?", sqlPort.ID).Delete(&sec)
	db.Model(sqlPort).Delete(&sqlPort)
	// Create new one
	db.Create(&reqPort)

	msg := fmt.Sprintf("Updated portfolio with name %s", name)
	data.ToJSON(&ResponseMessage{
		Msg: msg,
	}, w)

}
func (c *ControlHandler) updateDB(w http.ResponseWriter, port *Portfolio) {
	// Update prices using gRPC API
	c.updatePrices(port)

	// Update database entry using GORM
	// Open sqlite db connection
	db, err := gorm.Open(sqlite.Open("portfolios.db"), &gorm.Config{})
	if err != nil {
		c.LogHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}

	// Get generic sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		c.LogHTTPError(w, "Couldn't create sql instance", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	// Update associations
	var dbPort Portfolio
	db.Model(&dbPort).Association("Securities").Replace(&port.Securities)

}

func (c *ControlHandler) updatePrices(port *Portfolio) {
	// Concurrently retrieve stock prices
	wg := &sync.WaitGroup{}
	for _, sec := range port.Securities {
		wg.Add(1)
		go func(s *Security) {
			defer wg.Done()
			// Get security information using Info method defined in driver.go
			st, err := Info(s.Ticker, s.Currency, c.client)
			if err != nil {
				c.l.Println("[ERROR] Couldn't get info for ticker:", s.Ticker)
				return
			}
			// Parse the stock price
			price, err := strconv.ParseFloat(st.Price, 64)
			if err != nil {
				c.l.Println("[ERROR] Couldn't parse stock price for ticker:", s.Ticker, "price:", st.Price)
				return
			}
			// Set stock price in target currency (USD by default)
			if s.Currency == "" {
				s.Currency = "USD"
			}
			c.l.Println("[DEBUG] Got price for ticker:", s.Ticker, "in", s.Currency)
			s.CurrPrice = price

			// Update the individual security's gains and percent changes
			gain, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (s.CurrPrice-s.BoughtPrice)*s.Shares), 64)
			s.setMoves(gain, fmt.Sprintf("%.2f%%", (s.CurrPrice-s.BoughtPrice)/s.BoughtPrice*100))

		}(sec)
	}
	wg.Wait()

}
