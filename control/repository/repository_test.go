package repository

import (
	"errors"
	mocks "github.com/fercevik729/STLKER/control/mocks/repository"
	"github.com/fercevik729/STLKER/control/models"
	"testing"
)

func TestPortfolioRepository_CreateNewPortfolio(t *testing.T) {
	type args struct {
		portfolio models.Portfolio
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "missing_securities",
			args: args{
				portfolio: models.Portfolio{
					Name:       "Malicious portfolio",
					Username:   "maliciousIndividual",
					Securities: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "missing_portfolio",
			args: args{
				portfolio: models.Portfolio{},
			},
			wantErr: true,
		},
		{
			name: "successful_create_portfolio",
			args: args{
				portfolio: models.Portfolio{
					Name:     "My first portfolio",
					Username: "fercevik",
					Securities: []*models.Security{
						{
							Ticker:      "F",
							BoughtPrice: 5.0,
							CurrPrice:   12.03,
							Shares:      100,
							Currency:    "USD",
							PortfolioID: 0,
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.IPortfolioRepository{}
			if !tt.wantErr {
				repo.On("CreateNewPortfolio", tt.args.portfolio).Return(nil)
			} else {
				repo.On("CreateNewPortfolio", tt.args.portfolio).Return(errors.New("failed to create new portfolio"))
			}
			err := repo.CreateNewPortfolio(tt.args.portfolio)
			if (err != nil) != tt.wantErr {
				t.Errorf("PortfolioRepository.CreatePortfolio() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}
