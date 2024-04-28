package migrations_schema

import (
	"database/sql"
	"errors"

	"github.com/gouthamkrishnakv/zerocounter/database/schema/create_table"
)

// -- variables --

// DatabaseNotFoundErr returns an error when database provided is null
var ErrDatabaseNotFound = errors.New("database not found")

// -- functions --

// Run method runs all the database migrations on the database
func Run(db *sql.DB) (string, error) {
	if db == nil {
		return "", ErrDatabaseNotFound
	}

	// Create table methods
	for table, createTableFunc := range create_table.CreateTables {
		if createTableErr := createTableFunc(db); createTableErr != nil {
			return table, createTableErr
		}
	}

	return "", nil
}
