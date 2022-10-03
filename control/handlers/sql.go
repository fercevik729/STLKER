package handlers

import (
	"fmt"
	"reflect"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// replacePortfolio replaces a portfolio of name "portName" for a user "username" with
// a new portfolio struct "newPort"
func (c *ControlHandler) replacePortfolio(portName string, username string, newPort *Portfolio) error {
	// Declare vars
	var (
		port Portfolio
		sec  Security
	)
	// Create a new gorm db connection
	db, err := newGormDBConn(c.dsn)
	if err != nil {
		return err
	}
	// Check if any results were found
	db.Where("name=?", portName).Where("username=?", username).Preload("Securities").Find(&port)
	if reflect.DeepEqual(port, &Portfolio{}) {
		return fmt.Errorf("no results could be found for portfolio %s and username %s", portName, username)
	}
	// Delete the securities and then the portfolio
	db.Model(&sec).Where("portfolio_id=?", port.ID).Delete(sec)
	db.Model(&port).Delete(&port)

	// If a new portfolio is specified create it in place of the old one
	if newPort != nil {
		db.Create(newPort)
	}

	return nil

}

// getPortfolioId returns a portfolio's id provided its name and the username associated with it
func getPortfolioId(db *gorm.DB, portName string, username string) int {
	var port Portfolio
	db.Model(&Portfolio{}).Select("id").Where("name=?", portName).Where("username=?", username).First(&port)
	return int(port.ID)

}

// newGormDBConn opens a new gorm database connection
func newGormDBConn(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil

}
