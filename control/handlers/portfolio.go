package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/fercevik729/STLKER/control/data"
	"github.com/gorilla/mux"
)

type Username struct{}

type IsAdmin struct{}

type NamePair struct {
	Name     string
	Username string
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

// swagger: parameters createPortfolio updatePortfolio
type ReqPortfolio struct {
	// A single portfolio
	// in: body
	Body Portfolio
}

// A Portfolio defines the structure for an API portfolio
// swagger:model
type Portfolio struct {
	// swagger: ignore
	STLKERModel
	// the name of the portfolio
	//
	// required: true
	Name string `json:"Name"`
	// username of the portfolio's owner
	//
	// required: true
	// example: MoneyLover123
	Username string `json:"Username"`
	// Stocks is a list of Security structures
	//
	// required: true
	Securities []*Security `json:"Securities" gorm:"foreignKey:PortfolioID"`
}

// A Profits struct defines the structure for the profits of an API portfoli
// swagger:model
type Profits struct {
	// portfolio name
	//
	// example: Retirement Account
	Name string `json:"Portfolio Name"`
	// original value of the portfolio
	//
	OriginalValue float64 `json:"Original Value"`
	// new value of the portfolio
	NewValue float64 `json:"Current Value"`
	NetGain  float64 `json:"Net Gain"`
	// change of the portfolio's value as a percentage
	NetChange string `json:"Net Change"`
	// list of all the securities
	Moves []*Security `json:"Securities"`
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
		Name:          p.Name,
		OriginalValue: original,
		NewValue:      new,
		Moves:         p.Securities,
		NetGain:       netGain,
		NetChange:     percChange,
	}, nil
}

// swagger:route POST /portfolios portfolios createPortfolio
// Creates a new portfolio for a user
// responses:
//  200: messageResponse
//  400: errorResponse
//  500: errorResponse
func (c *ControlHandler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	// Retrieve the portfolio from the request body
	reqPort := Portfolio{}
	data.FromJSON(&reqPort, r.Body)

	// Validate the portfolio
	ok, msg := validatePortfolio(&reqPort)
	if !ok {
		c.logHTTPError(w, msg, http.StatusBadRequest)
		return
	}
	// Set username of the requested portfolio
	username := retrieveUsername(r)
	reqPort.Username = username

	// Open sqlite db connection
	db, err := newGormDBConn(c.dbName)
	if err != nil {
		c.logHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}
	// Migrate schema
	db.AutoMigrate(&Portfolio{}, &Security{})

	sqlPort := Portfolio{}
	// Check if a portfolio with that name for that user already exists
	db.Debug().Where("name=? AND username=?", reqPort.Name, username).First(&sqlPort)

	// If a portfolio with that name does exist return an error
	if !reflect.DeepEqual(&sqlPort, &Portfolio{}) {
		c.logHTTPError(w, "A portfolio with that name already exists", http.StatusBadRequest)
		return
	}
	// Retrieve prices for all stocks in the portfolio
	c.l.Println("[INFO] Retrieving updated stock prices")
	c.updatePrices(&reqPort)

	// Create portfolio entry
	db.Debug().Create(&reqPort)
	msg = fmt.Sprintf("Created portfolio named %s for %s", reqPort.Name, reqPort.Username)
	c.l.Printf("[DEBUG] %s", msg)

	// Write to response body
	w.WriteHeader(http.StatusCreated)
	data.ToJSON(&ResponseMessage{
		Msg: msg,
	}, w)

}

// swagger:route GET /portfolios portfolios getPortfolios
// Outputs all of the portfolios for a user
// responses:
//  200: profitsResponse
//  400: errorResponse
//  500: errorResponse
func (c *ControlHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	username := retrieveUsername(r)
	isAdmin := retrieveAdmin(r)
	var ports []Portfolio

	// Open database
	db, err := newGormDBConn(c.dbName)
	if err != nil {
		c.logHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}
	// If the username is admin, retrieve all portfolio names and associated usernames
	// Then return in an ordered format
	if isAdmin {
		var (
			usernames []string
			portnames []string
		)
		if err != nil {
			c.logHTTPError(w, "Couldn't open user's database", http.StatusInternalServerError)
			return
		}
		// Get all usernames except for admin
		db.Model(&User{}).Not("username=?", "admin").Select("username").Find(&usernames)

		// Create a map of usernames to slices of portfolio names
		table := make(map[string][]string)

		// Iterate ove all users
		for _, user := range usernames {
			db.Table("portfolios").Where("username=?", user).Select("name").Find(&portnames)
			table[user] = portnames
		}
		// Return the table in json format
		data.ToJSON(table, w)
		if err != nil {
			c.logHTTPError(w, "Couldn't set value into cache", http.StatusInternalServerError)
			return
		}
	}
	// Otherwise retrieve all portfolio data for a user
	db.Where("username=?", username).Preload("Securities").Find(&ports)

	profits := make([]*Profits, 0)
	for i := range ports {
		prof, err := ports[i].calcProfits()
		if err != nil {
			c.logHTTPError(w, fmt.Sprintf("Couldn't calculate profits for %s", ports[i].Name), http.StatusInternalServerError)
			return
		}
		profits = append(profits, prof)
	}

	data.ToJSON(profits, w)

}

// swagger:route GET /portfolios/{name} portfolios getPortfolio
// Outputs a particulare portfolio for a user
// responses:
//  200: profitResponse
//  400: errorResponse
//  500: errorResponse
func (c *ControlHandler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	// Retrieve portfolio name parameter and username
	name := mux.Vars(r)["name"]
	username := retrieveUsername(r)

	// Open sqlite db connection
	db, err := newGormDBConn(c.dbName)
	if err != nil {
		c.logHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}

	var port Portfolio
	// Check if a portfolio with that name for that user can be found
	db.Where("name=?", name).Where("username=?", username).Preload("Securities").Find(&port)
	// Check if any results were found
	if port.ID == 0 {
		c.logHTTPError(w, fmt.Sprintf("no results found with name %s, and user %s", name, username), http.StatusBadRequest)
		return
	}
	// Update the database entry with the new prices
	err = c.updateDB(&port)
	if err != nil {
		c.logHTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Calculate the profits
	profits, err := port.calcProfits()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data.ToJSON(profits, w)

}

// swagger:route PUT /portfolios portfolios updatePortfolio
// Updates a given portfolio
// responses:
//  200: messageResponse
//  400: errorResponse
//  500: errorResponse
func (c *ControlHandler) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	// Get variables
	username := retrieveUsername(r)

	// Retrieve the portfolio from the request body
	reqPort := Portfolio{}
	data.FromJSON(&reqPort, r.Body)
	reqPort.Username = username
	name := reqPort.Name

	// Check if request payload is empty
	if reflect.DeepEqual(reqPort, Portfolio{}) {
		c.logHTTPError(w, "Bad request payload", http.StatusBadRequest)
		return
	}
	// Call helper method
	err := c.replacePortfolio(name, username, &reqPort)
	if err != nil {
		c.logHTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send response message
	msg := fmt.Sprintf("Updated portfolio with name %s", name)
	data.ToJSON(&ResponseMessage{
		Msg: msg,
	}, w)

}

// swagger:route DELETE /portfolios/{name} securities deletePortfolio
// Deletes a given portfolio for a user
// responses:
//  200: messageResponse
//  400: errorResponse
//  500: errorResponse
func (c *ControlHandler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	username := retrieveUsername(r)
	// Delete portfolio and all child securities
	err := c.replacePortfolio(name, username, nil)
	if err != nil {
		c.logHTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}
	data.ToJSON(&ResponseMessage{
		Msg: fmt.Sprintf("Deleted portfolio %s", name),
	}, w)
}
