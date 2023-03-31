// Filename: BIOAFF/backend/internal/data/models.go
package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// A wrapper for all of our models
type Models struct {
	PublicUsers PublicUserModel
	Forms       FormModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		PublicUsers: PublicUserModel{DB: db},
		Forms:       FormModel{DB: db},
	}
}
