package repository

import (
	m "github.com/fercevik729/STLKER/control/models"
	"gorm.io/gorm"
	"reflect"
)

type SecurityRepository interface {
	CreateSecurity(security m.Security)
	Exists(portfolio uint, ticker string) bool
	GetSecurity(portId uint, ticker string) m.Security
	UpdateSecurity(security m.Security)
	DeleteSecurity(ticker string, portId uint)
}

// securityRepositoryImpl is a struct used to abstract data access operations that implements the SecurityRepository
// interface
type securityRepositoryImpl struct {
	db *gorm.DB
}

// NewSecurityRepository constructs a new SecurityRepository struct
func NewSecurityRepository(db *gorm.DB) SecurityRepository {
	return &securityRepositoryImpl{db: db}
}

// CreateSecurity creates a new security
func (s *securityRepositoryImpl) CreateSecurity(security m.Security) {
	s.db.Create(&security)
}

func (s *securityRepositoryImpl) Exists(portfolioId uint, ticker string) bool {
	res := s.GetSecurity(portfolioId, ticker)
	return reflect.DeepEqual(&res, &m.Security{})
}

// GetSecurity retrieves a security of a particular ticker within a user's portfolio
func (s *securityRepositoryImpl) GetSecurity(portId uint, ticker string) m.Security {
	var security m.Security
	s.db.Model(&m.Security{}).Select([]string{"ticker", "bought_price", "curr_price", "shares", "gain", "change"}).Where("ticker=?", ticker).Where("portfolio_id=?", portId).First(&security)
	return security
}

// UpdateSecurity updates a security in a user's portfolio if it exists, otherwise it creates a new one
func (s *securityRepositoryImpl) UpdateSecurity(security m.Security) {
	var res m.Security
	s.db.Model(&res).Where("portfolio_id=?", security.PortfolioID).Where("ticker=?", security.Ticker).Update("shares", security.Shares)
}

// DeleteSecurity deletes a security from a user's portfolio
func (s *securityRepositoryImpl) DeleteSecurity(ticker string, portId uint) {
	s.db.Model(&m.Security{}).Where("ticker=?", ticker).Where("portfolio_id=?", portId).Delete(&s)
}
