package repository

import (
	m "github.com/fercevik729/STLKER/control/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) CreateUser(user m.User) {
	u.db.Model(&m.User{}).Create(&user)
}

func (u *UserRepository) GetUser(username string) m.User {
	var res m.User
	u.db.Where("username=?", username).First(&res)
	return res
}

func (u *UserRepository) DeleteUser(username string) {
	var res m.User
	u.db.Model(&m.User{}).Where("username=?", username).Delete(&res)
}
