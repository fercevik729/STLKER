package repository

import (
	m "github.com/fercevik729/STLKER/control/models"
	"gorm.io/gorm"
	"reflect"
)

type ISecurityRepository interface {
	CreateSecurity(security m.Security)
	Exists(portfolio uint, ticker string) bool
	GetSecurity(portId uint, ticker string) m.Security
	UpdateSecurity(security m.Security)
	DeleteSecurity(ticker string, portId uint)
}

// securityRepository is a struct used to abstract data access operations that implements the ISecurityRepository
// interface
type securityRepository struct {
	db *gorm.DB
}

// NewSecurityRepository constructs a new ISecurityRepository struct
func NewSecurityRepository(db *gorm.DB) ISecurityRepository {
	return &securityRepository{db: db}
}

// CreateSecurity creates a new security
func (s securityRepository) CreateSecurity(security m.Security) {
	s.db.Create(&security)
}

func (s securityRepository) Exists(portfolioId uint, ticker string) bool {
	res := s.GetSecurity(portfolioId, ticker)
	return reflect.DeepEqual(&res, &m.Security{})
}

// GetSecurity retrieves a security of a particular ticker within a user's portfolio
func (s securityRepository) GetSecurity(portId uint, ticker string) m.Security {
	var security m.Security
	s.db.Model(&m.Security{}).Select([]string{"ticker", "bought_price", "curr_price", "shares", "gain", "change"}).Where("ticker=?", ticker).Where("portfolio_id=?", portId).First(&security)
	return security
}

// UpdateSecurity updates a security in a user's portfolio if it exists, otherwise it creates a new one
func (s securityRepository) UpdateSecurity(security m.Security) {
	var res m.Security
	s.db.Model(&res).Where("portfolio_id=?", security.PortfolioID).Where("ticker=?", security.Ticker).Update("shares", security.Shares)
}

// DeleteSecurity deletes a security from a user's portfolio
func (s securityRepository) DeleteSecurity(ticker string, portId uint) {
	s.db.Model(&m.Security{}).Where("ticker=?", ticker).Where("portfolio_id=?", portId).Delete(&s)
}
