package models

import (
	"database/sql"
	"time"
)

type STLKERModel struct {
	ID        uint         `gorm:"primaryKey" json:"-"`
	CreatedAt time.Time    `                  json:"-"`
	UpdatedAt time.Time    `                  json:"-"`
	DeletedAt sql.NullTime `gorm:"index"      json:"-"`
}
