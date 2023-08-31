package repository

import (
	m "github.com/fercevik729/STLKER/control/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user m.User)
	GetUser(username string) m.User
	DeleteUser(username string)
}

// userRepositoryImpl is a struct used to abstract data access operations that implements the UserRepository
// interface
type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) CreateUser(user m.User) {
	u.db.Model(&m.User{}).Create(&user)
}

func (u *userRepositoryImpl) GetUser(username string) m.User {
	var res m.User
	u.db.Where("username=?", username).First(&res)
	return res
}

func (u *userRepositoryImpl) DeleteUser(username string) {
	var res m.User
	u.db.Model(&m.User{}).Where("username=?", username).Delete(&res)
}
