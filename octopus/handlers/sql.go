package handlers

import (
	"database/sql"
	"fmt"
	"reflect"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// replacePortfolio replaces a portfolio of name "portName" for a user "username" with
// a new portfolio struct "newPort"
func replacePortfolio(portName string, username string, newPort *Portfolio) error {
	// Declare vars
	var (
		port Portfolio
		sec  Security
	)
	// Create a new gorm db connection
	db, err := newGormDBConn(databasePath)
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
func getPortfolioId(db *sql.DB, portName string, username string) (int, error) {
	// Execute query
	rows, err := db.Query("SELECT id FROM portfolios WHERE name=? AND username=?", portName, username)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return -1, err
	}
	// Grab the id of the portfolio
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return -1, err
		}
	}
	return id, nil

}

// newGormDBConn opens a new gorm database connection
func newGormDBConn(databaseName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil

}

// newSqlDBConn opens a new sqlite3 database connection
func newSqlDBConn(databaseName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
