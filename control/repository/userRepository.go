package repository

import (
	m "github.com/fercevik729/STLKER/control/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user m.User)
	GetUser(username string) m.User
	DeleteUser(username string)
}

// userRepository is a struct used to abstract data access operations that implements the IUserRepository
// interface
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) CreateUser(user m.User) {
	u.db.Model(&m.User{}).Create(&user)
}

func (u *userRepository) GetUser(username string) m.User {
	var res m.User
	u.db.Where("username=?", username).First(&res)
	return res
}

func (u *userRepository) DeleteUser(username string) {
	var res m.User
	u.db.Model(&m.User{}).Where("username=?", username).Delete(&res)
}
