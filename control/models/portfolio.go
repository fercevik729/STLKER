package models

import (
	"fmt"
	"strconv"
)

// swagger:model
type Portfolio struct {
	STLKERModel
	// the name of the portfolio
	//
	// in: string
	// required: true
	Name string `json:"Name"`
	// username of the portfolio's owner
	//
	// in: string
	// required: true
	// example: MoneyLover123
	Username string `json:"Username"`
	// Stocks is a list of Security structures
	//
	// in: Security
	// required: true
	Securities []*Security `json:"Securities" gorm:"foreignKey:PortfolioID"`
}

// A Profits struct defines the structure for the profits of an API portfolio
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

// CalcProfits is used to calculate the total curren profits of a portfolio
func (p *Portfolio) CalcProfits() (*Profits, error) {
	original := 0.
	newProfits := 0.
	// Iterate over securities and calculate change and percent change
	for _, sec := range p.Securities {
		original += sec.BoughtPrice * sec.Shares
		newProfits += sec.CurrPrice * sec.Shares
	}

	// Round original and new
	var err error
	original, err = strconv.ParseFloat(fmt.Sprintf("%.2f", original), 64)
	if err != nil {
		return nil, err
	}

	newProfits, err = strconv.ParseFloat(fmt.Sprintf("%.2f", newProfits), 64)
	if err != nil {
		return nil, err
	}

	// Compute and round change and profits
	percChange := fmt.Sprintf("%.2f%%", (newProfits-original)/original*100)

	netGain, err := strconv.ParseFloat(fmt.Sprintf("%.2f", newProfits-original), 64)
	if err != nil {
		return nil, err
	}

	// Return profits
	return &Profits{
		Name:          p.Name,
		OriginalValue: original,
		NewValue:      newProfits,
		Moves:         p.Securities,
		NetGain:       netGain,
		NetChange:     percChange,
	}, nil
}
