package models

import "fmt"

// Security defines the structure for a security
// swagger:model
type Security struct {
	// swagger: ignore
	STLKERModel
	// swagger: ignore
	SecurityID int `gorm:"primary_key" json:"-"`
	// ticker of the security
	Ticker      string  `json:"Ticker"`
	BoughtPrice float64 `json:"BoughtPrice"`
	CurrPrice   float64 `json:"CurrentPrice"`
	Shares      float64 `json:"Shares"`
	Gain        float64 `json:"Gain"`
	Change      string  `json:"PercentChange"`
	// Currency is the destination currency of the stock
	Currency string `json:"Currency" gorm:"default:USD"`
	// Foreign key
	// swagger: ignore
	PortfolioID uint `json:"-"`
}

// SetMoves sets the gain and change variables of s to the new parameters
func (s *Security) SetMoves(gain float64, change string) {
	s.Gain = gain
	s.Change = change
}

func (s *Security) String() string {
	return fmt.Sprintf("Security: ticker=%v, bought price=%v, curr price=%v", s.Ticker, s.BoughtPrice, s.CurrPrice)
}
