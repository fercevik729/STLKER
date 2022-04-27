package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDBConn(databaseName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil

}

func NewSqlDBConn(databaseName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *ControlHandler) LogHTTPError(w http.ResponseWriter, errorMsg string, errorCode int) {
	c.l.Printf("[ERROR] %s\n", errorMsg)
	http.Error(w, fmt.Sprintf("Error: %s", errorMsg), errorCode)
}

func (c *ControlHandler) RetrieveUsername(r *http.Request) string {
	// Get email from request context
	username := r.Context().Value(Username{})
	c.l.Println("[INFO] Got username:", username)

	v, ok := username.(string)
	if ok {
		return v
	}
	return ""
}

func (c *ControlHandler) RetrieveAdmin(r *http.Request) bool {
	// Get email from request context
	isAdmin := r.Context().Value(IsAdmin{})
	c.l.Println("[INFO] User is admin:", isAdmin)

	v, ok := isAdmin.(bool)
	if ok {
		return v
	}
	return false
}
