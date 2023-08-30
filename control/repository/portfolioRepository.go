package repository

import (
	"fmt"
	"github.com/fercevik729/STLKER/control/handlers"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
)

// PortfolioRepository is a struct used to abstract data access operations
type PortfolioRepository struct {
	db *gorm.DB
}

// NewPortfolioRepository constructs a new PortfolioRepository struct and returns a pointer to it
func NewPortfolioRepository(dsn string) (*PortfolioRepository, error) {
	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Init the schemas
	err = db.AutoMigrate(&handlers.Portfolio{}, &handlers.Security{}, &handlers.User{})
	return &PortfolioRepository{db: db}, nil
}

// GetPortfolio retrieves a portfolio. Returns an error if a portfolio couldn't be found
func (r *PortfolioRepository) GetPortfolio(portName, username string) (handlers.Portfolio, error) {
	// Run query
	var res handlers.Portfolio
	r.db.Debug().Where("name=?", portName).Where("username=?", username).First(&res)

	// Check if a portfolio couldn't be found
	if !reflect.DeepEqual(&res, &handlers.Portfolio{}) {
		return handlers.Portfolio{}, errors.Errorf("no portfolio of name %s, belonging to user %s",
			portName, username)
	}
	return res, nil
}

// CreateNewPortfolio creates a new portfolio if a portfolio with the same name doesn't already exist for a user
func (r *PortfolioRepository) CreateNewPortfolio(portName, username string, portfolio handlers.Portfolio) error {
	// Check if portfolio is empty
	if reflect.DeepEqual(portfolio, handlers.Portfolio{}) {
		return errors.Errorf("portfolio cannot be empty")
	}

	// Run query
	var res handlers.Portfolio
	r.db.Debug().Where("name=?", portName).Where("username=?", username).First(&res)
	// Check if a portfolio already exists
	if reflect.DeepEqual(&res, &portfolio) {
		return errors.Errorf("a portfolio of name %s, belonging to user %s already exists", portName, username)
	}
	r.db.Create(&portfolio)

	return nil
}

// DeletePortfolio deletes a portfolio and all its associated securities
func (r *PortfolioRepository) DeletePortfolio(portName, username string) error {
	// Run queries
	var (
		port handlers.Portfolio
		sec  handlers.Security
	)
	// Check if any results were found
	r.db.Where("name=?", portName).Where("username=?", username).Preload("Securities").Find(&port)
	if reflect.DeepEqual(port, &handlers.Portfolio{}) {
		return fmt.Errorf("no results could be found for portfolio %s and username %s", portName, username)
	}
	// Delete the securities and then the portfolio
	r.db.Model(&sec).Where("portfolio_id=?", port.ID).Delete(sec)
	r.db.Model(&port).Delete(&port)

	return nil
}

// UpdatePortfolio updates a portfolio and all its associated securities by deleting the previous version of the
// portfolio and creating a new version
func (r *PortfolioRepository) UpdatePortfolio(portName, username string, portfolio handlers.Portfolio) error {
	// Delete the previous version of the portfolio
	if err := r.DeletePortfolio(portName, username); err != nil {
		return err
	}
	// Create the new version of the portfolio
	return r.CreateNewPortfolio(portName, username, portfolio)
}

// GetPortfolioId retrieves the id of a portfolio
func (r *PortfolioRepository) GetPortfolioId(portName, username string) uint {
	var res handlers.Portfolio
	r.db.Where("name=?", portName).Where("username=?", username).First(&res)
	return res.ID
}
