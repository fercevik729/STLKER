package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fercevik729/STLKER/control/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

const dummyPortName = "My first portfolio"
const dummyUsername = "fercevik"

var dummyPortfolio = models.Portfolio{
	Name:     dummyPortName,
	Username: dummyUsername,
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
}

func TestPortfolioRepository_CreateNewPortfolio1(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	pg := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(pg, &gorm.Config{})
	portRepo := NewPortfolioRepository(db)
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT")

	err := portRepo.CreateNewPortfolio(dummyPortfolio)
	if err != nil {
		t.Errorf("got unexpected error %v", err)
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("all expectations were not met: %v", err)
		return
	}
}

/*
func TestPortfolioRepository_GetPortfolio(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	pg := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	// Add the data first
	db, _ := gorm.Open(pg, &gorm.Config{})
	portRepo := NewPortfolioRepository(db)
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT")
	mock.ExpectCommit()

	if err := portRepo.CreateNewPortfolio(dummyPortfolio); err != nil {
		t.Errorf("go unexpected error %v", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("all expectations were not met: %v", err)
		return

	}
	// Get the data
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("0"))

	res, err := portRepo.GetPortfolio(dummyPortName, dummyUsername)
	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}
	t.Log("Log:", res.Name, res.Username)
	if dummyPortName != res.Name || dummyUsername != res.Username {
		t.Errorf("Got %v, wanted %v", res, dummyPortfolio)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("all expectations were not met: %v", err)
		return
	}
}
*/
