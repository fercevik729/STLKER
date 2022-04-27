package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fercevik729/STLKER/octopus/data"
	"github.com/gorilla/mux"
)

type Username struct{}

type IsAdmin struct{}

type NamePair struct {
	Name     string
	Username string
}

const databasePath string = "./database/stlker.db"

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
	Name     string `json:"Name"`
	Username string `json:"Username"`
	// Stocks is a slice of Security structs
	Securities []*Security `json:"Securities" gorm:"foreignKey:PortfolioID"`
}

type Profits struct {
	Name          string      `json:"Portfolio Name"`
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
		Name:          p.Name,
		OriginalValue: original,
		NewValue:      new,
		Moves:         p.Securities,
		NetGain:       netGain,
		NetChange:     percChange,
	}, nil
}

func (c *ControlHandler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	// Retrieve the portfolio from the request body
	reqPort := Portfolio{}
	data.FromJSON(&reqPort, r.Body)

	// Check if name is empty which generally signifies that the json body was misconstrued
	// or if the name contains spaces, since it isn't compatible with the URI
	if reqPort.Name == "" || strings.Contains(reqPort.Name, " ") {
		c.logHTTPError(w, "Bad portfolio request. Name shouldn't be empty or contain spaces", http.StatusBadRequest)
		return
	}

	// Set username of the requested portfolio
	username := c.retrieveUsername(r)
	reqPort.Username = username

	c.l.Println("[INFO] Creating portfolio:", reqPort.Name, "for user:", username)

	// Open sqlite db connection
	db, err := newGormDBConn(databasePath)
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
	msg := fmt.Sprintf("Created portfolio named %s for %s", reqPort.Name, reqPort.Username)
	c.l.Printf("[DEBUG] %s", msg)

	// Write to response body
	w.WriteHeader(http.StatusCreated)
	data.ToJSON(&ResponseMessage{
		Msg: msg,
	}, w)

}

func (c *ControlHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	username := c.retrieveUsername(r)
	isAdmin := c.retrieveAdmin(r)
	c.l.Printf("[INFO] Getting all portfolios for user: %s\n", username)

	var ports []Portfolio

	// Open database
	db, err := newGormDBConn(databasePath)
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
		return

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

func (c *ControlHandler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	// Retrieve portfolio name parameter and username
	name := mux.Vars(r)["name"]
	username := c.retrieveUsername(r)
	c.l.Println("[INFO] Getting portfolio:", name, "and user:", username)

	// Open sqlite db connection
	db, err := newGormDBConn(databasePath)
	if err != nil {
		c.logHTTPError(w, "Couldn't connect to database", http.StatusInternalServerError)
		return
	}

	var port Portfolio
	// Check if a portfolio with that name for that user can be found
	db.Where("name=?", name).Where("username=?", username).Preload("Securities").Find(&port)
	// Check if any results were found
	if port.ID == 0 {
		c.l.Println("[DEBUG] No results found")
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

func (c *ControlHandler) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	// Get variables
	name := mux.Vars(r)["name"]
	username := c.retrieveUsername(r)

	// Retrieve the portfolio from the request body
	reqPort := Portfolio{}
	data.FromJSON(&reqPort, r.Body)
	reqPort.Username = username

	c.l.Println("[INFO] Updating portfolio:", name, "for user:", username)
	// Call helper method
	err := replacePortfolio(name, username, &reqPort)
	if err != nil {
		c.logHTTPError(w, err.Error(), http.StatusBadRequest)
	}

	// Send response message
	msg := fmt.Sprintf("Updated portfolio with name %s", name)
	data.ToJSON(&ResponseMessage{
		Msg: msg,
	}, w)

}

func (c *ControlHandler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	username := c.retrieveUsername(r)
	c.l.Println("[INFO] Deleting portfolio for:", name, "and user:", username)
	// Delete portfolio and all child securities
	err := replacePortfolio(name, username, nil)
	if err != nil {
		c.logHTTPError(w, err.Error(), http.StatusBadRequest)
	}
	data.ToJSON(&ResponseMessage{
		Msg: fmt.Sprintf("Deleted portfolio %s", name),
	}, w)
}
