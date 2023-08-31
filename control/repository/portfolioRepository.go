package repository

import (
	"fmt"
	m "github.com/fercevik729/STLKER/control/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"reflect"
)

type PortfolioRepository interface {
	GetPortfolio(portName, username string) (m.Portfolio, error)
	GetAllPortfolios(username string) []m.Portfolio
	GetAllPortfoliosAdmin() map[string][]string
	CreateNewPortfolio(portName, username string, portfolio m.Portfolio) error
	DeletePortfolio(portName, username string) error
	UpdatePortfolio(portName, username string, portfolio m.Portfolio) error
	GetPortfolioId(portName, username string) uint
}

// portfolioRepositoryImpl is a struct used to abstract data access operations that implements the PortfolioRepository
// interface
type portfolioRepositoryImpl struct {
	db *gorm.DB
}

// NewPortfolioRepository constructs a new PortfolioRepository struct and returns a pointer to it
func NewPortfolioRepository(db *gorm.DB) PortfolioRepository {
	return &portfolioRepositoryImpl{db: db}
}

// GetPortfolio retrieves a portfolio. Returns an error if a portfolio couldn't be found
func (r *portfolioRepositoryImpl) GetPortfolio(portName, username string) (m.Portfolio, error) {
	// Run query
	var res m.Portfolio
	r.db.Where("name=?", portName).Where("username=?", username).First(&res)

	// Check if a portfolio couldn't be found
	if !reflect.DeepEqual(&res, &m.Portfolio{}) {
		return m.Portfolio{}, errors.Errorf("no portfolio of name %s, belonging to user %s",
			portName, username)
	}
	return res, nil
}

// GetAllPortfolios retrieves all portfolios for a user
func (r *portfolioRepositoryImpl) GetAllPortfolios(username string) []m.Portfolio {
	var ports []m.Portfolio
	r.db.Where("username=?", username).Preload("Securities").Find(&ports)
	return ports
}

// GetAllPortfoliosAdmin retrieves all portfolio names for all users. Intended to be used by admin users only
func (r *portfolioRepositoryImpl) GetAllPortfoliosAdmin() map[string][]string {
	var (
		usernames      []string
		portfolioNames []string
	)

	r.db.Model(&m.User{}).Not("username=?", "admin").Select("username").Find(&usernames)

	// Create a map of usernames to slices of portfolio names
	table := make(map[string][]string)

	// Create the map of usernames -> list of portfolios
	for _, user := range usernames {
		r.db.Table("portfolios").Where("username=?", user).Select("name").Find(&portfolioNames)
		table[user] = portfolioNames
	}

	return table
}

// CreateNewPortfolio creates a new portfolio if a portfolio with the same name doesn't already exist for a user
func (r *portfolioRepositoryImpl) CreateNewPortfolio(portName, username string, portfolio m.Portfolio) error {
	// Check if portfolio is empty
	if reflect.DeepEqual(portfolio, m.Portfolio{}) {
		return errors.Errorf("portfolio cannot be empty")
	}

	// Check if a portfolio already exists
	var res m.Portfolio
	r.db.Where("name=?", portName).Where("username=?", username).First(&res)
	if reflect.DeepEqual(&res, &portfolio) {
		return errors.Errorf("a portfolio of name %s, belonging to user %s already exists", portName, username)
	}
	// Create the portfolio
	r.db.Create(&portfolio)

	return nil
}

// DeletePortfolio deletes a portfolio and all its associated securities
func (r *portfolioRepositoryImpl) DeletePortfolio(portName, username string) error {
	var (
		port m.Portfolio
		sec  m.Security
	)
	// Check if any matching portfolios were found
	r.db.Where("name=?", portName).Where("username=?", username).Preload("Securities").Find(&port)
	if reflect.DeepEqual(port, &m.Portfolio{}) {
		return fmt.Errorf("no results could be found for portfolio %s and username %s", portName, username)
	}
	// Delete the securities and then the portfolio
	r.db.Model(&sec).Where("portfolio_id=?", port.ID).Delete(sec)
	r.db.Model(&port).Delete(&port)

	return nil
}

// UpdatePortfolio updates a portfolio and all its associated securities by deleting the previous version of the
// portfolio and creating a new version
func (r *portfolioRepositoryImpl) UpdatePortfolio(portName, username string, portfolio m.Portfolio) error {
	// Delete the previous version of the portfolio
	if err := r.DeletePortfolio(portName, username); err != nil {
		return err
	}
	// Create the new version of the portfolio
	return r.CreateNewPortfolio(portName, username, portfolio)
}

// GetPortfolioId retrieves the id of a portfolio
func (r *portfolioRepositoryImpl) GetPortfolioId(portName, username string) uint {
	var res m.Portfolio
	r.db.Where("name=?", portName).Where("username=?", username).First(&res)
	return res.ID
}
