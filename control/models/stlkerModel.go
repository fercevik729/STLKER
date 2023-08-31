package models

import (
	"database/sql"
	"time"
)

// STLKERModel is a Gorm Model to be composed in other models
type STLKERModel struct {
	ID        uint         `gorm:"primaryKey" json:"-"`
	CreatedAt time.Time    `                  json:"-"`
	UpdatedAt time.Time    `                  json:"-"`
	DeletedAt sql.NullTime `gorm:"index"      json:"-"`
}
