package handlers

import (
	"net/http"
	"reflect"

	"github.com/fercevik729/STLKER/octopus/data"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO: change sql database schema

// A Portfolio is a GORM model that is intended to mirror the structure
// of a simple portfolio
type Portfolio struct {
	gorm.Model
	ID uint `gorm:"primary_key"`
	// Name is the name of the portfolio
	Name string `json:"Name"`
	// Stocks is a slice of Security structs
	Securities []Security `json:"Securities" gorm:"foreignKey:ID"`
}

type Security struct {
	gorm.Model
	ID          uint
	Ticker      string  `json:"Ticker"`
	BoughtPrice float64 `json:"Bought Price"`
	CurrPrice   float64 `json:"Current Price"`
	Shares      float64 `json:"Shares"`
}

func (c *ControlHandler) SavePortfolio(w http.ResponseWriter, r *http.Request) {
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

	// Retrieve the portfolio from the request body
	reqPortfolio := Portfolio{}
	// body, _ := ioutil.ReadAll(r.Body)
	//fmt.Println(string(body))
	data.FromJSON(&reqPortfolio, r.Body)

	sqlPort := Portfolio{}
	// Check if a portfolio with that name already exists
	db.First(&sqlPort, "name = ?", reqPortfolio.Name)

	// If a portfolio with that name does exist return an error
	if !reflect.DeepEqual(&sqlPort, &Portfolio{}) {
		c.l.Println("[ERROR] A portfolio with that name already exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Retrieve prices for all stocks in the portfolio
	c.l.Println("[INFO] Retrieving updated stock prices")
	// c.updatePortfolio(&port)

	// Create portfolio entry
	db.Create(&reqPortfolio)
	c.l.Println("[DEBUG] Created portfolio named", reqPortfolio.Name)
	db.Save(&reqPortfolio)

	// Close database connection
	sqlDB.Close()

}

/*
func (c *ControlHandler) updatePortfolio(port *Portfolio) {
	// Concurrently retrieve stock prices
	stocks := port.St.List
	var wg *sync.WaitGroup
	for _, stock := range stocks {
		wg.Add(1)
		go func(s *Stock) {
			st, err := Info(s.Ticker, s.Destination, c.client)
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
			// Set stock price
			c.l.Println("[DEBUG] Got price for ticker:", s.Ticker)
			s.updatePrice(price)
			wg.Done()
		}(&stock)
	}
	wg.Wait()

}
*/

func (c *ControlHandler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	// Retrieve portfolio name parameter
	name := mux.Vars(r)["name"]

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
	db.First(&port, "name = ?", name)

	// Check if portfolio is empty
	if reflect.DeepEqual(port, Portfolio{}) {
		c.l.Println("[DEBUG] No results found")
		return
	}
	data.ToJSON(port, w)

}
